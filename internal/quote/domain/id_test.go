package domain

import (
	"sync"
	"testing"
)

func TestIDGenerator_Generate(t *testing.T) {
	gen := &idGenerator{}

	t.Run("生成的 ID 在有效范围（17-19 位数字）", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			id := gen.Generate(1)
			if id < 10000000000000000 {
				t.Fatalf("ID %d 太小（小于 17 位）", id)
			}
			if id < 0 {
				t.Fatalf("ID %d 为负数", id)
			}
		}
	})

	t.Run("ID 单调递增", func(t *testing.T) {
		gen := &idGenerator{}
		var prev int64
		for i := 0; i < 100; i++ {
			id := gen.Generate(1)
			if i > 0 && id <= prev {
				t.Fatalf("ID 非单调递增: prev=%d, current=%d", prev, id)
			}
			prev = id
		}
	})

	t.Run("ID 互不重复", func(t *testing.T) {
		gen := &idGenerator{}
		ids := make(map[int64]bool, 1000)
		for i := 0; i < 1000; i++ {
			id := gen.Generate(1)
			if ids[id] {
				t.Fatalf("重复 ID: %d", id)
			}
			ids[id] = true
		}
	})

	t.Run("并发生成不 panic", func(t *testing.T) {
		gen := &idGenerator{}
		var wg sync.WaitGroup
		ids := make(chan int64, 1000)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					ids <- gen.Generate(1)
				}
			}()
		}

		go func() {
			wg.Wait()
			close(ids)
		}()

		uniqueIDs := make(map[int64]bool)
		for id := range ids {
			if uniqueIDs[id] {
				t.Fatalf("并发场景下重复 ID: %d", id)
			}
			uniqueIDs[id] = true
		}
		if len(uniqueIDs) != 1000 {
			t.Fatalf("期望 1000 个唯一 ID, 得到 %d", len(uniqueIDs))
		}
	})
}
