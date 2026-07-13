package custom

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"
)

func pickCrossShardPair(prefix string, n int) (ParsedNode, ParsedNode, bool) {
	for i := 0; i < 300; i++ {
		a := tunnelNode(fmt.Sprintf("%s-a-%d", prefix, i), fmt.Sprintf("%s-a-%d.example.com", prefix, i), "p")
		b := tunnelNode(fmt.Sprintf("%s-b-%d", prefix, i), fmt.Sprintf("%s-b-%d.example.com", prefix, i), "p")
		ia := shardIndexForKey(a.NodeKey(), n)
		ib := shardIndexForKey(b.NodeKey(), n)
		if ia != ib {
			if ia == 0 {
				return a, b, true
			}
			return b, a, true
		}
	}
	return ParsedNode{}, ParsedNode{}, false
}

// oncePartialShard: first Reload after arm() becomes Partial success; later reloads full success.
type oncePartialShard struct {
	*spyShard
	mu   sync.Mutex
	armN int
}

func (s *oncePartialShard) Reload(nodes []ParsedNode) error {
	s.mu.Lock()
	if s.armN > 0 {
		s.armN--
		s.spyShard.forcePartial = true
	} else {
		s.spyShard.forcePartial = false
	}
	s.mu.Unlock()
	return s.spyShard.Reload(nodes)
}

// failAfterShard fails reloads after the first successful one (used for rollback failure).
type failAfterShard struct {
	*spyShard
	mu        sync.Mutex
	successes int
	failAfter int
	failMsg   string
}

func (s *failAfterShard) Reload(nodes []ParsedNode) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.successes >= s.failAfter {
		s.reloadCalls++
		return errors.New(s.failMsg)
	}
	err := s.spyShard.Reload(nodes)
	if err == nil {
		s.successes++
	}
	return err
}

func TestShardedReloadRollsBackEarlierSuccessfulShardWhenLaterShardFails(t *testing.T) {
	var spies []*spyShard
	factory := func(shardIndex, shardBasePort int) singBoxShard {
		s := newSpyShard()
		spies = append(spies, s)
		return s
	}
	sb := newShardedSingBoxWithFactory(20000, 2, factory)

	old0, old1, ok := pickCrossShardPair("old", 2)
	if !ok {
		t.Fatal("no cross-shard seed pair")
	}
	if err := sb.Reload([]ParsedNode{old0, old1}); err != nil {
		t.Fatalf("seed Reload: %v", err)
	}
	seed0 := spies[0].calls()

	new0, new1, ok := pickCrossShardPair("new", 2)
	if !ok {
		t.Fatal("no cross-shard new pair")
	}
	spies[1].reloadErr = errors.New("injected shard1 reload failure")

	err := sb.Reload([]ParsedNode{new0, new1})
	if err == nil {
		t.Fatal("Reload expected error from shard1, got nil")
	}
	if !strings.Contains(err.Error(), "injected shard1 reload failure") {
		t.Fatalf("error = %v", err)
	}

	keys0 := map[string]bool{}
	for _, n := range spies[0].GetNodes() {
		keys0[n.NodeKey()] = true
	}
	if keys0[new0.NodeKey()] {
		t.Fatalf("shard0 still has new node after rollback; keys=%v", keys0)
	}
	if !keys0[old0.NodeKey()] {
		t.Fatalf("shard0 missing old node after rollback; keys=%v", keys0)
	}
	if spies[0].calls() < seed0+2 {
		t.Fatalf("shard0 calls=%d seed=%d; want forward+rollback", spies[0].calls(), seed0)
	}
}

func TestShardedReloadRollsBackPartialFailingShard(t *testing.T) {
	var partial *oncePartialShard
	factory := func(shardIndex, shardBasePort int) singBoxShard {
		partial = &oncePartialShard{spyShard: newSpyShard()}
		return partial
	}
	sb := newShardedSingBoxWithFactory(21000, 1, factory)

	old := tunnelNode("partial-old", "partial-old.example.com", "old")
	if err := sb.Reload([]ParsedNode{old}); err != nil {
		t.Fatalf("seed: %v", err)
	}
	seedCalls := partial.calls()
	partial.armN = 1 // next Reload is Partial; rollback Reload is clean

	newNode := tunnelNode("partial-new", "partial-new.example.com", "new")
	err := sb.Reload([]ParsedNode{newNode})
	if err == nil {
		t.Fatal("Reload expected partial commit error, got nil")
	}

	keys := map[string]bool{}
	for _, n := range partial.GetNodes() {
		keys[n.NodeKey()] = true
	}
	if keys[newNode.NodeKey()] {
		t.Fatalf("new node remained after partial rollback; keys=%v", keys)
	}
	if !keys[old.NodeKey()] {
		t.Fatalf("old node missing after partial rollback; keys=%v", keys)
	}
	if partial.calls() < seedCalls+2 {
		t.Fatalf("calls=%d seed=%d; want forward+rollback", partial.calls(), seedCalls)
	}
}

func TestShardedReloadReportsPrimaryAndRollbackErrors(t *testing.T) {
	var s0 *failAfterShard
	var s1 *spyShard
	factory := func(shardIndex, shardBasePort int) singBoxShard {
		if shardIndex == 0 {
			s0 = &failAfterShard{spyShard: newSpyShard(), failAfter: 2, failMsg: "rollback shard0 failure"}
			// seed success + forward success = 2, then rollback fails
			return s0
		}
		s1 = newSpyShard()
		return s1
	}
	sb := newShardedSingBoxWithFactory(22000, 2, factory)

	old0, old1, ok := pickCrossShardPair("rb-old", 2)
	if !ok {
		t.Fatal("no seed pair")
	}
	if err := sb.Reload([]ParsedNode{old0, old1}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	new0, new1, ok := pickCrossShardPair("rb-new", 2)
	if !ok {
		t.Fatal("no new pair")
	}
	s1.reloadErr = errors.New("primary shard1 failure")

	err := sb.Reload([]ParsedNode{new0, new1})
	if err == nil {
		t.Fatal("expected joined primary+rollback error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "primary shard1 failure") {
		t.Fatalf("missing primary: %v", msg)
	}
	if !strings.Contains(msg, "rollback shard0 failure") {
		t.Fatalf("missing rollback: %v", msg)
	}
}
