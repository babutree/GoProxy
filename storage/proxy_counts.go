package storage

import "database/sql"

// GetAverageLatency 获取指定协议的平均延迟
func (s *Storage) GetAverageLatency(protocol string) (int, error) {
	var avg sql.NullFloat64
	err := s.db.QueryRow(
		`SELECT AVG(latency) FROM proxies
		 WHERE protocol = ? AND status = 'active' AND user_paused = 0 AND latency > 0
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')`,
		protocol,
	).Scan(&avg)
	if err != nil || !avg.Valid {
		return 0, err
	}
	return int(avg.Float64), nil
}

// GetQualityDistribution 获取质量分布统计
func (s *Storage) GetQualityDistribution() (map[string]int, error) {
	rows, err := s.db.Query(
		`SELECT quality_grade, COUNT(*) as count 
		 FROM proxies 
		 WHERE status = 'active' AND user_paused = 0 AND fail_count < 3
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')
		 GROUP BY quality_grade`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dist := make(map[string]int)
	for rows.Next() {
		var grade string
		var count int
		if err := rows.Scan(&grade, &count); err != nil {
			return nil, err
		}
		dist[grade] = count
	}
	return dist, nil
}

// CountAll 返回所有可用节点数量
func (s *Storage) CountAll() (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM proxies
		 WHERE status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')`,
	).Scan(&count)
	return count, err
}

// CountAvailableByProtocol 按协议统计所有可用节点数量
func (s *Storage) CountAvailableByProtocol(protocol string) (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM proxies
		 WHERE status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3 AND protocol = ?
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')`,
		protocol,
	).Scan(&count)
	return count, err
}

// CountBySource 按来源统计可用代理数量
func (s *Storage) CountBySource(source string) (int, error) {
	var count int
	err := s.db.QueryRow(
		`SELECT COUNT(*) FROM proxies
		 WHERE source = ? AND status IN ('active', 'degraded') AND user_paused = 0 AND fail_count < 3
		   AND NOT EXISTS (SELECT 1 FROM subscriptions WHERE subscriptions.id = proxies.subscription_id AND subscriptions.status = 'paused')`,
		source,
	).Scan(&count)
	return count, err
}
