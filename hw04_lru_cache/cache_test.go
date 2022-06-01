package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("test cache Clear", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		c.Clear()

		wasInCache = c.Set("aaa", 300)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 300)
		require.False(t, wasInCache)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic: capacity limit", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 4; i++ { // set: a->101, b->102, c->103, d->104
			r := 'a' + i
			c.Set(Key(rune(r)), 100+i)
		}

		// no (a->101), cache: b->102, c->103, d->104
		_, ok := c.Get("a")
		require.False(t, ok)

		val, ok := c.Get("b")
		require.True(t, ok)
		require.Equal(t, 101, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 102, val)

		val, ok = c.Get("d")
		require.True(t, ok)
		require.Equal(t, 103, val)
	})

	t.Run("purge logic: old items", func(t *testing.T) {
		c := NewCache(3)
		for i := 0; i < 3; i++ { // set: a->101, b->102, c->103
			r := 'a' + i
			c.Set(Key(rune(r)), 100+i)
		}

		// cache (by order: front - .. - back): 202(b) - 103(c) - 101(a)
		c.Set(Key('b'), 202)

		// cache (by order: front - .. - back): 203(c) - 202(b) - 101(a)
		c.Set(Key('c'), 203)

		// cache (by order: front - .. - back): 201(a) - 203(c) - 202(b)
		c.Set(Key('a'), 201)

		// cache (by order: front - .. - back): 104(d) - 201(a) - 203(c)  	no: (202(b))
		c.Set(Key('d'), 104)

		_, ok := c.Get("b")
		require.False(t, ok)

		val, ok := c.Get("d")
		require.True(t, ok)
		require.Equal(t, 104, val)

		val, ok = c.Get("a")
		require.True(t, ok)
		require.Equal(t, 201, val)

		val, ok = c.Get("c")
		require.True(t, ok)
		require.Equal(t, 203, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
