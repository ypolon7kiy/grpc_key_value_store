package in_memory_store

import (
	"testing"

	in_memory_store "server/stores"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	store := in_memory_store.NewStore()

	// Test setting a value
	store.Set("key1", "value1")
	value, exists := store.Get("key1")

	assert.True(t, exists)
	assert.Equal(t, "value1", value)
}

func TestGet(t *testing.T) {
	store := in_memory_store.NewStore()

	// Test getting a value that exists
	store.Set("key1", "value1")
	value, exists := store.Get("key1")

	assert.True(t, exists)
	assert.Equal(t, "value1", value)

	// Test getting a value that does not exist
	value, exists = store.Get("nonexistent_key")
	assert.False(t, exists)
	assert.Equal(t, "", value)
}

func TestDelete(t *testing.T) {
	store := in_memory_store.NewStore()

	store.Set("key1", "value1")
	store.Delete("key1")
	value, exists := store.Get("key1")

	assert.False(t, exists)
	assert.Equal(t, "", value)

	// Test deleting a key that does not exist
	store.Delete("nonexistent_key")
}

func TestConcurrentAccess(t *testing.T) {
	store := in_memory_store.NewStore()

	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			store.Set("key", "value")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			store.Get("key")
		}
		done <- true
	}()

	<-done
	<-done
}

func TestSetAndGetConcurrently(t *testing.T) {
	store := in_memory_store.NewStore()

	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			store.Set("key", "value")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			_, _ = store.Get("key")
		}
		done <- true
	}()

	<-done
	<-done

	// Verify final value
	value, exists := store.Get("key")
	assert.True(t, exists)
	assert.Equal(t, "value", value)
}

func TestSetDeleteAndGetConcurrentlyDeadLocks(t *testing.T) {
	store := in_memory_store.NewStore()

	done := make(chan bool)

	go func() {
		for i := 0; i < 1000; i++ {
			store.Set("key", "value")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			store.Delete("key")
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			_, _ = store.Get("key")
		}
		done <- true
	}()

	<-done
	<-done
	<-done
}
