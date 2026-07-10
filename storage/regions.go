package storage

import (
	"database/sql"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

var alpha2RegionPattern = regexp.MustCompile(`^[A-Za-z]{2}$`)

type RegionCount struct {
	Region string `json:"region"`
	Count  int    `json:"count"`
}

func (s *Storage) GetByRegion(region string, excludes []int64) ([]Proxy, error) {
	query := `SELECT ` + proxyColumns + `
		 FROM proxies
		 WHERE status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')`
	args := []interface{}{}
	if normalized := normalizeRegion(region); normalized != "" {
		query += ` AND region = ?`
		args = append(args, normalized)
	}
	excludeMap := makeExcludeMap(excludes)
	query += ` ORDER BY CASE WHEN latency <= 0 THEN 1 ELSE 0 END, latency ASC, RANDOM()`
	return s.queryProxies(query, args, excludeMap)
}

func (s *Storage) GetRandomByRegion(region string, excludes []int64) (*Proxy, error) {
	proxies, err := s.GetByRegion(region, excludes)
	if err != nil {
		return nil, err
	}
	if len(proxies) == 0 {
		return nil, fmt.Errorf("no available proxy for region: %s", normalizeRegion(region))
	}
	proxy := proxies[rand.Intn(len(proxies))]
	return &proxy, nil
}

func (s *Storage) CountByRegion() (map[string]int, error) {
	rows, err := s.db.Query(`
		SELECT region, COUNT(*)
		FROM proxies
		WHERE status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3 AND region != ''
		  AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')
		GROUP BY region`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var region string
		var count int
		if err := rows.Scan(&region, &count); err != nil {
			return nil, err
		}
		counts[region] = count
	}
	return counts, rows.Err()
}

func (s *Storage) GetRegionsWithCount() ([]RegionCount, error) {
	rows, err := s.db.Query(`
		SELECT region, COUNT(*)
		FROM proxies
		WHERE status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3 AND region != ''
		  AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')
		GROUP BY region
		ORDER BY region ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	regions := []RegionCount{}
	for rows.Next() {
		var item RegionCount
		if err := rows.Scan(&item.Region, &item.Count); err != nil {
			return nil, err
		}
		regions = append(regions, item)
	}
	return regions, rows.Err()
}

func (s *Storage) queryProxies(query string, args []interface{}, excludes map[int64]bool) ([]Proxy, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	proxies := []Proxy{}
	for rows.Next() {
		p, err := scanProxy(rows)
		if err != nil {
			return nil, err
		}
		if !excludes[p.ID] {
			proxies = append(proxies, *p)
		}
	}
	return proxies, rows.Err()
}

func makeExcludeMap(excludes []int64) map[int64]bool {
	excludeMap := make(map[int64]bool, len(excludes))
	for _, id := range excludes {
		excludeMap[id] = true
	}
	return excludeMap
}

func (s *Storage) GetProxyByAddress(address string) (*Proxy, error) {
	row := s.db.QueryRow(`SELECT `+proxyColumns+` FROM proxies WHERE address = ? ORDER BY source = ? DESC, subscription_id ASC LIMIT 1`, address, SourceManual)
	proxy, err := scanProxy(row)
	if err == nil {
		return proxy, nil
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("proxy %s not found", address)
	}
	return nil, err
}

func (s *Storage) GetProxyByIdentity(address, source string, subscriptionID int64) (*Proxy, error) {
	if source == SourceManual {
		subscriptionID = 0
	}
	row := s.db.QueryRow(
		`SELECT `+proxyColumns+` FROM proxies WHERE address = ? AND source = ? AND subscription_id = ?`,
		address, source, subscriptionID,
	)
	proxy, err := scanProxy(row)
	if err == nil {
		return proxy, nil
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("proxy %s/%s/%d not found", address, source, subscriptionID)
	}
	return nil, err
}

func (s *Storage) GetProxyByID(id int64) (*Proxy, error) {
	row := s.db.QueryRow(`SELECT `+proxyColumns+` FROM proxies WHERE id = ?`, id)
	proxy, err := scanProxy(row)
	if err == nil {
		return proxy, nil
	}
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("proxy id %d not found", id)
	}
	return nil, err
}

func normalizeProtocol(protocol string) string {
	return strings.ToLower(strings.TrimSpace(protocol))
}

func regionFromExitLocation(exitLocation string) string {
	fields := strings.Fields(exitLocation)
	if len(fields) == 0 || !alpha2RegionPattern.MatchString(fields[0]) {
		return ""
	}
	return strings.ToLower(fields[0])
}
