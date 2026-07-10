package storage

func (s *Storage) AddManualProxy(address, protocol, region, note string) error {
	region = normalizeRegion(region)
	regionSource := "auto"
	if region != "" {
		regionSource = "manual"
	}
	_, err := s.db.Exec(
		`INSERT INTO proxies (address, protocol, source, subscription_id, region, region_source, note)
		 VALUES (?, ?, 'manual', 0, ?, ?, ?)
		 ON CONFLICT(address, source, subscription_id) DO UPDATE SET
			protocol = excluded.protocol,
			region = excluded.region,
			region_source = excluded.region_source,
			note = excluded.note`,
		address, normalizeProtocol(protocol), region, regionSource, note,
	)
	return err
}

func (s *Storage) UpdateProxyRegion(address, region string, manual bool) error {
	regionSource := "auto"
	if manual {
		regionSource = "manual"
	}
	res, err := s.db.Exec(
		`UPDATE proxies SET region = ?, region_source = ? WHERE address = ?`,
		normalizeRegion(region), regionSource, address,
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
		normalizeRegion(region), regionSource, id,
	)
	if err != nil {
		return err
	}
	return requireRowsAffected(res.RowsAffected())
}

func (s *Storage) UpdateProxyNote(address, note string) error {
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
