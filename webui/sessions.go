package webui

import (
	"log"
	"net/http"
)

type sessionRow struct {
	SessionID           string `json:"session_id"`
	Node                string `json:"node"`
	Region              string `json:"region"`
	RemainingTTLSeconds int64  `json:"remaining_ttl_seconds"`
}

// proxyOccupancyRow is the per-node occupancy snapshot for lease observability (#16).
type proxyOccupancyRow struct {
	ProxyID                  int64  `json:"proxy_id"`
	Address                  string `json:"address"`
	ActiveSessions           int    `json:"active_sessions"`
	MaxSessions              int    `json:"max_sessions"`
	CooldownRemainingSeconds int64  `json:"cooldown_remaining_seconds"`
}

func (s *Server) apiSessions(w http.ResponseWriter, _ *http.Request) {
	bindings := s.affinity.List()
	rows := make([]sessionRow, 0, len(bindings))
	for _, binding := range bindings {
		rows = append(rows, sessionRow{
			SessionID:           binding.SessionID,
			Node:                binding.NodeAddress,
			Region:              binding.Region,
			RemainingTTLSeconds: int64(s.affinity.RemainingTTL(binding).Seconds()),
		})
	}
	jsonOK(w, rows)
}

// apiProxyOccupancy returns per-proxy active session counts for authenticated admins.
// Only proxies with at least one non-expired binding are included. No credential fields.
func (s *Server) apiProxyOccupancy(w http.ResponseWriter, _ *http.Request) {
	if s.affinity == nil {
		jsonOK(w, []proxyOccupancyRow{})
		return
	}
	bindings := s.affinity.List()
	counts := make(map[int64]int)
	addressByID := make(map[int64]string)
	for _, binding := range bindings {
		if binding.ProxyID <= 0 {
			continue
		}
		counts[binding.ProxyID]++
		if _, ok := addressByID[binding.ProxyID]; !ok {
			addressByID[binding.ProxyID] = binding.NodeAddress
		}
	}
	maxSessions := 0
	if s.cfg != nil {
		maxSessions = s.cfg.MaxSessionsPerProxy
	}
	rows := make([]proxyOccupancyRow, 0, len(counts))
	for proxyID, active := range counts {
		addr := addressByID[proxyID]
		if p, err := s.storage.GetProxyByID(proxyID); err == nil && p != nil {
			addr = p.Address
		} else if err != nil {
			log.Printf("[webui] proxy occupancy lookup id=%d: %v", proxyID, err)
		}
		rows = append(rows, proxyOccupancyRow{
			ProxyID:                  proxyID,
			Address:                  addr,
			ActiveSessions:           active,
			MaxSessions:              maxSessions,
			CooldownRemainingSeconds: 0, // #15 not merged
		})
	}
	jsonOK(w, rows)
}
