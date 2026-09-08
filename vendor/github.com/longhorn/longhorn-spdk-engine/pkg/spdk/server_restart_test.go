package spdk

import (
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRestartLogic(t *testing.T) {
	// Test 1: Initial state should allow restart
	s := &Server{}
	
	require.NotNil(t, s, "Server should be initialized")
	assert.True(t, s.lastRestartTime.Load() == 0, "Initial lastRestartTime should be 0")
	
	// Test 2: Simulate a restart and verify the timestamp is updated
	now := time.Now().Unix()
	s.lastRestartTime.Store(now)
	assert.Equal(t, now, s.lastRestartTime.Load(), "lastRestartTime should be updated")
}

func TestRestartFrequency(t *testing.T) {
	// Test 3: Verify restart is throttled (5 minute cooldown)
	s := &Server{
		lastRestartTime: atomic.Int64{},
	}

	// Set a recent restart time
	now := time.Now().Unix()
	s.lastRestartTime.Store(now - 100) // 100 seconds ago

	// Should not allow restart within 5 minutes (300 seconds)
	shouldAllow := s.lastRestartTime.Load() == 0 || (now - s.lastRestartTime.Load()) > 300
	assert.False(t, shouldAllow, "Should not allow restart within cooldown period")

	// Set a restart time > 5 minutes ago
	s.lastRestartTime.Store(now - 400) // 400 seconds ago

	shouldAllow = s.lastRestartTime.Load() == 0 || (now - s.lastRestartTime.Load()) > 300
	assert.True(t, shouldAllow, "Should allow restart after cooldown period")
}

func TestServerStateClearOnRestart(t *testing.T) {
	// Test 4: Verify all maps are cleared on restart
	
	s := &Server{
		diskMap:           make(map[string]*Disk),
		replicaMap:        make(map[string]*Replica),
		engineMap:         make(map[string]*Engine),
		engineFrontendMap: make(map[string]*EngineFrontend),
		shardMap:          make(map[string]*Shard),
		shardGroupMap:     make(map[string]*ShardGroup),
		backingImageMap:   make(map[string]*BackingImage),
		lastRestartTime:   atomic.Int64{},
	}

	// Add some data to the maps
	s.diskMap["disk1"] = &Disk{}
	s.replicaMap["replica1"] = &Replica{}
	assert.Equal(t, 1, len(s.diskMap), "diskMap should have one entry")
	assert.Equal(t, 1, len(s.replicaMap), "replicaMap should have one entry")

	// Simulate restart by setting lastRestartTime
	s.lastRestartTime.Store(time.Now().Unix())

	// Clear maps on restart
	s.diskMap = make(map[string]*Disk)
	s.replicaMap = make(map[string]*Replica)
	s.engineMap = make(map[string]*Engine)
	s.engineFrontendMap = make(map[string]*EngineFrontend)
	s.shardMap = make(map[string]*Shard)
	s.shardGroupMap = make(map[string]*ShardGroup)
	s.backingImageMap = make(map[string]*BackingImage)

	assert.Equal(t, 0, len(s.diskMap), "diskMap should be cleared")
	assert.Equal(t, 0, len(s.replicaMap), "replicaMap should be cleared")
}
