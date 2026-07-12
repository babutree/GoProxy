package affinity

import (
	"sync"
	"time"
)

type Binding struct {
	ProxyID     int64
	NodeAddress string
	Region      string
	LastActive  time.Time
}

type Store struct {
	mu       sync.RWMutex
	bindings map[string]Binding
	// reverse maps proxy_id -> set of session_ids currently bound (non-expired).
	// Kept in sync with bindings by SetProxy/Remove/Get-expiry/GC.
	reverse map[int64]map[string]struct{}
	ttl     time.Duration
	now     func() time.Time

	// GC lifecycle fields, guarded by mu.
	gcStarted bool
	stopCh    chan struct{}
	doneCh    chan struct{}
}

// SessionBinding is a read-only snapshot of a single active session binding,
// suitable for exposing to a WebUI session-monitor panel.
type SessionBinding struct {
	SessionID   string
	ProxyID     int64
	NodeAddress string
	Region      string
	LastActive  time.Time
}

func New(ttl time.Duration) *Store {
	return NewWithClock(ttl, time.Now)
}

func NewWithClock(ttl time.Duration, now func() time.Time) *Store {
	return &Store{
		bindings: map[string]Binding{},
		reverse:  map[int64]map[string]struct{}{},
		ttl:      ttl,
		now:      now,
	}
}

func (s *Store) Get(sessionID string) (Binding, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	binding, ok := s.bindings[sessionID]
	if !ok {
		return Binding{}, false
	}
	if s.expired(binding) {
		s.removeBindingLocked(sessionID, binding)
		return Binding{}, false
	}
	binding.LastActive = s.now()
	s.bindings[sessionID] = binding
	return binding, true
}

func (s *Store) Set(sessionID string, nodeAddress string, region string) {
	s.SetProxy(sessionID, 0, nodeAddress, region)
}

func (s *Store) SetProxy(sessionID string, proxyID int64, nodeAddress string, region string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if old, ok := s.bindings[sessionID]; ok {
		s.detachReverseLocked(sessionID, old.ProxyID)
	}
	s.bindings[sessionID] = Binding{ProxyID: proxyID, NodeAddress: nodeAddress, Region: region, LastActive: s.now()}
	s.attachReverseLocked(sessionID, proxyID)
}

func (s *Store) Remove(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if binding, ok := s.bindings[sessionID]; ok {
		s.removeBindingLocked(sessionID, binding)
	}
}

// CountByProxy returns the number of non-expired sessions currently bound to proxyID.
// It purges any reverse-index entries whose forward binding is missing or expired.
func (s *Store) CountByProxy(proxyID int64) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.countByProxyLocked(proxyID)
}

func (s *Store) countByProxyLocked(proxyID int64) int {
	sessions, ok := s.reverse[proxyID]
	if !ok || len(sessions) == 0 {
		return 0
	}
	for sessionID := range sessions {
		binding, ok := s.bindings[sessionID]
		if !ok || binding.ProxyID != proxyID {
			delete(sessions, sessionID)
			continue
		}
		if s.expired(binding) {
			s.removeBindingLocked(sessionID, binding)
			continue
		}
	}
	if len(sessions) == 0 {
		delete(s.reverse, proxyID)
		return 0
	}
	return len(sessions)
}

func (s *Store) attachReverseLocked(sessionID string, proxyID int64) {
	set, ok := s.reverse[proxyID]
	if !ok {
		set = map[string]struct{}{}
		s.reverse[proxyID] = set
	}
	set[sessionID] = struct{}{}
}

func (s *Store) detachReverseLocked(sessionID string, proxyID int64) {
	set, ok := s.reverse[proxyID]
	if !ok {
		return
	}
	delete(set, sessionID)
	if len(set) == 0 {
		delete(s.reverse, proxyID)
	}
}

func (s *Store) removeBindingLocked(sessionID string, binding Binding) {
	delete(s.bindings, sessionID)
	s.detachReverseLocked(sessionID, binding.ProxyID)
}

func (s *Store) expired(binding Binding) bool {
	return s.ttl > 0 && s.now().Sub(binding.LastActive) >= s.ttl
}

// StartGC starts a background goroutine that scans bindings every interval and
// deletes expired ones. Subsequent calls are ignored while a GC goroutine is
// already running (call Stop first to restart). It is a no-op for interval <= 0.
func (s *Store) StartGC(interval time.Duration) {
	if interval <= 0 {
		return
	}

	s.mu.Lock()
	if s.gcStarted {
		s.mu.Unlock()
		return
	}
	s.gcStarted = true
	stopCh := make(chan struct{})
	doneCh := make(chan struct{})
	s.stopCh = stopCh
	s.doneCh = doneCh
	s.mu.Unlock()

	go s.gcLoop(interval, stopCh, doneCh)
}

func (s *Store) gcLoop(interval time.Duration, stopCh, doneCh chan struct{}) {
	defer close(doneCh)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			s.collectExpired()
		}
	}
}

// collectExpired removes every expired binding in a single locked pass.
func (s *Store) collectExpired() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for sessionID, binding := range s.bindings {
		if s.expired(binding) {
			s.removeBindingLocked(sessionID, binding)
		}
	}
}

// Stop gracefully stops the GC goroutine. It is safe to call once even if
// StartGC was never called, and safe to call after a prior Stop; it never
// panics and does not leak the goroutine.
func (s *Store) Stop() {
	s.mu.Lock()
	if !s.gcStarted {
		s.mu.Unlock()
		return
	}
	s.gcStarted = false
	stopCh := s.stopCh
	doneCh := s.doneCh
	s.stopCh = nil
	s.doneCh = nil
	s.mu.Unlock()

	close(stopCh)
	<-doneCh
}

// List returns a snapshot of all active (non-expired) bindings. It is
// read-only: it does not refresh LastActive and does not delete expired
// entries. Expired bindings are skipped.
func (s *Store) List() []SessionBinding {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]SessionBinding, 0, len(s.bindings))
	for sessionID, binding := range s.bindings {
		if s.expired(binding) {
			continue
		}
		result = append(result, SessionBinding{
			SessionID:   sessionID,
			ProxyID:     binding.ProxyID,
			NodeAddress: binding.NodeAddress,
			Region:      binding.Region,
			LastActive:  binding.LastActive,
		})
	}
	return result
}

// Count returns the number of active (non-expired) bindings. It is read-only:
// it does not refresh LastActive and does not delete expired entries.
func (s *Store) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	count := 0
	for _, binding := range s.bindings {
		if !s.expired(binding) {
			count++
		}
	}
	return count
}

// TTL returns the configured session time-to-live. The UI can combine this
// with SessionBinding.LastActive to compute a countdown.
func (s *Store) TTL() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ttl
}

func (s *Store) SetTTL(ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ttl = ttl
}

// RemainingTTL returns how long until the given binding expires, based on the
// store's clock. It returns 0 once the binding is at or past expiry, and 0 when
// no TTL is configured (ttl <= 0). This is read-only.
func (s *Store) RemainingTTL(binding SessionBinding) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.ttl <= 0 {
		return 0
	}
	remaining := s.ttl - s.now().Sub(binding.LastActive)
	if remaining < 0 {
		return 0
	}
	return remaining
}
