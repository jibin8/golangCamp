package main

import (
	"container/list"
	"log"
	"sync"
	"time"
)

// 滑动窗口：时间分格，每格一个计数器

type SlidingWindowLimiter struct {
	sync.Mutex
	Cap      int           // 总容量
	Segments int           // 时间格
	Interval time.Duration // 窗口总大小 时间单位
	ticker   *time.Ticker
	Queue    *list.List // 维护每个时间格内的数量
}

func NewSlidingWindowLimiter(cap, segs int, interval time.Duration) *SlidingWindowLimiter {
	limiter := &SlidingWindowLimiter{
		Cap:      cap,
		Segments: segs,
		Interval: interval,
		ticker:   time.NewTicker(time.Duration(int64(interval) / int64(segs))),
		Queue:    list.New(),
	}
	for i := 0; i < segs; i++ {
		limiter.Queue.PushBack(0)
	}
	go func() {
		for {
			<-limiter.ticker.C
			limiter.Lock()
			limiter.Queue.Remove(limiter.Queue.Front())
			limiter.Queue.PushBack(0)
			limiter.Unlock()
		}
	}()
	return limiter
}

func (limiter *SlidingWindowLimiter) Increase() {
	limiter.Lock()
	defer limiter.Unlock()
	limiter.Queue.Back().Value = limiter.Queue.Back().Value.(int) + 1
}

func (limiter *SlidingWindowLimiter) IsAvailable() bool {
	limiter.Lock()
	defer limiter.Unlock()
	return limiter.cur() < limiter.Cap
}

func (limiter SlidingWindowLimiter) cur() int {
	sum := 0
	for e := limiter.Queue.Front(); e != nil; e = e.Next() {
		if i, ok := e.Value.(int); ok {
			sum += i
		}
	}
	return sum
}

func main() {
	a := NewSlidingWindowLimiter(2, 5, time.Duration(1*time.Second))
	wg := sync.WaitGroup{}
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			a.Increase()
			if a.IsAvailable() {
				log.Printf("i: %d allow", i)
			} else {
				log.Printf("i: %d over", i)
			}
		}(i)
		time.Sleep(time.Millisecond * 400)
	}

	wg.Wait()
}
