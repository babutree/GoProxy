package storage

import "database/sql"

func (s *Storage) AddManualProxy(address, protocol, region, note string) error {
	return s.addManualProxyExec(s.db, address, protocol, region, note, "", "", "")
}

// AddManualProxyWithCredentials 与 AddManualProxy 相同，但持久化上游认证凭据。
// 凭据用于拨号/验证时注入出站握手；绝不写入日志或错误串。空串表示无需认证。
func (s *Storage) AddManualProxyWithCredentials(address, protocol, region, note, username, password string) error {
	return s.addManualProxyExec(s.db, address, protocol, region, note, username, password, "")
}

// AddManualProxyWithNodeKey 手工入库并可写入稳定 node_key（隧道手工节点用 ParsedNode.NodeKey）。
func (s *Storage) AddManualProxyWithNodeKey(address, protocol, region, note, nodeKey string) error {
	return s.addManualProxyExec(s.db, address, protocol, region, note, "", "", nodeKey)
}

type proxyExec interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (s *Storage) addManualProxyExec(exec proxyExec, address, protocol, region, note, username, password, nodeKey string) error {
	region = normalizeManualRegion(region)
	regionSource := "auto"
	if region != "" {
		regionSource = "manual"
	}
	// 手工节点默认 disabled：须经连通/出口/纯净度/AI/CF 验证通过后才 active 入选路。
	// proxy_username/proxy_password 持久化上游认证凭据（拨号时注入，绝不入日志）。
	// node_key：调用方传入隧道稳定键；空则用 protocol|address 派生（直连场景）。
	if nodeKey == "" {
		// 仅用 DSL 安全字符（无 |），避免复制 -node-key- 时被字符集拒绝。
		nodeKey = "manual:" + normalizeProtocol(protocol) + ":" + address
	}
	_, err := exec.Exec(
		`INSERT INTO proxies (address, protocol, source, subscription_id, region, region_source, note, status, proxy_username, proxy_password, node_key)
		 VALUES (?, ?, 'manual', 0, ?, ?, ?, 'disabled', ?, ?, ?)
		 ON CONFLICT(address, source, subscription_id) DO UPDATE SET
			protocol = excluded.protocol,
			region = excluded.region,
			region_source = excluded.region_source,
			note = excluded.note,
			status = 'disabled',
			proxy_username = excluded.proxy_username,
			proxy_password = excluded.proxy_password,
			node_key = excluded.node_key`,
		address, normalizeProtocol(protocol), region, regionSource, note, username, password, nodeKey,
	)
	return err
}

func (s *Storage) AddManualProxies(proxies []Proxy, region, note string) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, proxy := range proxies {
		if err := s.addManualProxyExec(tx, proxy.Address, proxy.Protocol, region, note, proxy.Username, proxy.Password, proxy.NodeKey); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Storage) UpdateProxyRegion(address, region string, manual bool) error {
	if err := s.requireUnambiguousAddress(address); err != nil {
		return err
	}
	regionSource := "auto"
	if manual {
		regionSource = "manual"
	}
	res, err := s.db.Exec(
		`UPDATE proxies SET region = ?, region_source = ? WHERE address = ?`,
		normalizeManualRegion(region), regionSource, address,
	)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}

func (s *Storage) UpdateProxyRegionByID(id int64, region string, manual bool) error {
	regionSource := "auto"
	if manual {
		regionSource = "manual"
	}
	res, err := s.db.Exec(
		`UPDATE proxies SET region = ?, region_source = ? WHERE id = ?`,
		normalizeManualRegion(region), regionSource, id,
	)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}

func (s *Storage) UpdateProxyNote(address, note string) error {
	if err := s.requireUnambiguousAddress(address); err != nil {
		return err
	}
	res, err := s.db.Exec(`UPDATE proxies SET note = ? WHERE address = ?`, note, address)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}

func (s *Storage) UpdateProxyNoteByID(id int64, note string) error {
	res, err := s.db.Exec(`UPDATE proxies SET note = ? WHERE id = ?`, note, id)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}

func (s *Storage) DeleteManualProxy(address string) error {
	res, err := s.db.Exec(`DELETE FROM proxies WHERE address = ? AND source = 'manual'`, address)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}
