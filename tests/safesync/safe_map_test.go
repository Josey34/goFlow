package safesync_test

import (
	"fmt"
	"sync"
	"testing"

	"goflow/safesync"
)

func TestSafeMapBasicOperations(t *testing.T) {
	sm := safesync.New[string, string]()

	// Set
	sm.Set("key1", "value1")

	// Get
	val, ok := sm.Get("key1")
	if !ok {
		t.Error("expected key1 to exist")
	}
	if val != "value1" {
		t.Errorf("expected 'value1', got '%s'", val)
	}

	// Non-existent key
	_, ok = sm.Get("nonexistent")
	if ok {
		t.Error("expected nonexistent key to not exist")
	}
}

func TestSafeMapDelete(t *testing.T) {
	sm := safesync.New[string, int]()

	sm.Set("key1", 42)
	val, ok := sm.Get("key1")
	if !ok || val != 42 {
		t.Error("expected key1 to exist")
	}

	sm.Delete("key1")
	_, ok = sm.Get("key1")
	if ok {
		t.Error("expected key1 to be deleted")
	}
}

func TestSafeMapLen(t *testing.T) {
	sm := safesync.New[string, string]()

	if sm.Len() != 0 {
		t.Error("expected initial length 0")
	}

	sm.Set("a", "1")
	sm.Set("b", "2")
	if sm.Len() != 2 {
		t.Errorf("expected length 2, got %d", sm.Len())
	}

	sm.Delete("a")
	if sm.Len() != 1 {
		t.Errorf("expected length 1, got %d", sm.Len())
	}
}

func TestSafeMapRange(t *testing.T) {
	sm := safesync.New[string, int]()

	sm.Set("a", 1)
	sm.Set("b", 2)
	sm.Set("c", 3)

	count := 0
	sum := 0
	sm.Range(func(k string, v int) bool {
		count++
		sum += v
		return true
	})

	if count != 3 {
		t.Errorf("expected 3 items, got %d", count)
	}
	if sum != 6 {
		t.Errorf("expected sum 6, got %d", sum)
	}
}

func TestSafeMapRangeBreak(t *testing.T) {
	sm := safesync.New[int, int]()

	for i := 0; i < 10; i++ {
		sm.Set(i, i*10)
	}

	count := 0
	sm.Range(func(k, v int) bool {
		count++
		if count >= 3 {
			return false // Break early
		}
		return true
	})

	if count != 3 {
		t.Errorf("expected 3 iterations, got %d", count)
	}
}

func TestSafeMapClear(t *testing.T) {
	sm := safesync.New[string, string]()

	sm.Set("a", "1")
	sm.Set("b", "2")

	if sm.Len() != 2 {
		t.Error("expected 2 items before clear")
	}

	sm.Clear()

	if sm.Len() != 0 {
		t.Error("expected 0 items after clear")
	}
}

func TestSafeMapThreadSafety(t *testing.T) {
	sm := safesync.New[int, int]()
	var wg sync.WaitGroup

	// 50 writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sm.Set(id*100+j, id*100+j)
			}
		}(i)
	}

	// 50 readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				sm.Get(id*100 + j)
			}
		}(i)
	}

	wg.Wait()

	if sm.Len() != 5000 {
		t.Errorf("expected 5000 items, got %d", sm.Len())
	}
}

func TestSafeMapIntKeyStringValue(t *testing.T) {
	sm := safesync.New[int, string]()

	sm.Set(1, "one")
	sm.Set(2, "two")
	sm.Set(3, "three")

	val, ok := sm.Get(2)
	if !ok || val != "two" {
		t.Error("expected 'two'")
	}
}

func TestSafeMapStringKeyIntValue(t *testing.T) {
	sm := safesync.New[string, int]()

	sm.Set("a", 100)
	sm.Set("b", 200)

	val, ok := sm.Get("a")
	if !ok || val != 100 {
		t.Error("expected 100")
	}
}

func TestSafeMapConcurrentReadWrite(t *testing.T) {
	sm := safesync.New[string, int]()
	var wg sync.WaitGroup

	// Start with some data
	for i := 0; i < 10; i++ {
		sm.Set(fmt.Sprintf("key%d", i), i)
	}

	// 20 concurrent readers
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				sm.Get(fmt.Sprintf("key%d", j))
			}
		}(i)
	}

	// 10 concurrent writers
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				sm.Set(fmt.Sprintf("newkey%d_%d", id, j), j)
			}
		}(i)
	}

	wg.Wait()

	// Should have original 10 + 100 new items
	if sm.Len() != 110 {
		t.Errorf("expected 110 items, got %d", sm.Len())
	}
}
