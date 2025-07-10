package smolcache_test

import (
	"fmt"
	"smol/smolcache"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Program testing can be used to show the presence of bugs, but never to show their absence!
// - Dijkstra

func TestCache_Parallel_goroutines(t *testing.T) {
	const parallelTasks = 10

	c := smolcache.New[int, string](parallelTasks, time.Millisecond*100)

	wg := sync.WaitGroup{}
	wg.Add(parallelTasks)

	for i := range parallelTasks {
		go func(j int) {
			defer wg.Done()
			c.Upsert(4, fmt.Sprint(j))

		}(i)
	}

	wg.Wait()
}

func TestCache_Parallel(t *testing.T) {
	c := smolcache.New[int, string](2, time.Millisecond*100)

	t.Run("write six", func(t *testing.T) {
		t.Parallel()
		c.Upsert(6, "six")
	})

	t.Run("write kuus", func(t *testing.T) {
		t.Parallel()
		c.Upsert(6, "kuus")
	})
}

func TestCache_TTL(t *testing.T) {
	t.Parallel()

	c := smolcache.New[string, string](5, time.Millisecond*100)
	c.Upsert("Norwegian", "Blue")

	got, found := c.Read("Norwegian")
	assert.True(t, found)
	assert.Equal(t, "Blue", got)

	time.Sleep(time.Millisecond * 200)

	got, found = c.Read("Norwegian")

	assert.False(t, found)
	assert.Equal(t, "", got)
}

func TestCache_MaxSize(t *testing.T) {
	t.Parallel()

	c := smolcache.New[int, int](3, time.Minute)

	c.Upsert(1, 1)
	c.Upsert(2, 2)
	c.Upsert(3, 3)

	got, found := c.Read(1)
	assert.True(t, found)
	assert.Equal(t, 1, got)

	c.Upsert(1, 10)

	c.Upsert(4, 4)

	got, found = c.Read(2)
	assert.False(t, found)
	assert.Equal(t, 0, got)
}
