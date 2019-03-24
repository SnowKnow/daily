package syncBenchmark

import (
	"sync"
	"testing"
	"time"
)

func TestSyncDo(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			addSync()
			wg.Done()
		}()
	}
	time.Sleep(100 * time.Millisecond)
	wg.Wait()
	if i != 1000 {
		t.Errorf("result is wrong,want 1000,res :%d", i)
	}
}

func TestGoroutineDo(t *testing.T) {
	var wg sync.WaitGroup
	c = make(chan bool)
	go func() {
		for {
			select {
			case <-c:
				{
					add()
				}
			}
		}
	}()
	for i := 0; i < 1000; i++ {
		go func() {
			wg.Add(1)
			c <- true
			wg.Done()

		}()
	}
	time.Sleep(100 * time.Millisecond)
	wg.Wait()
	if i != 1000 {
		t.Errorf("result is wrong,want 1000,res :%d", i)
	}
}

func BenchmarkSyncDo(b *testing.B) {
	var wg sync.WaitGroup
	c = make(chan bool)

	for i := 0; i < b.N; i++ {
		go func() {
			wg.Add(1)
			addSync()
			wg.Done()
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)

}

//BenchmarkSyncDo-16    	 2000000	       506 ns/op	       0 B/op	       0 allocs/op
//BenchmarkGoroutineDo-16    	  200000	      8918 ns/op	     423 B/op	       1 allocs/op
//
func BenchmarkGoroutineDo(b *testing.B) {
	var wg sync.WaitGroup

	go func() {
		for {
			select {
			case <-c:
				{
					wg.Add(1)
					add()
					wg.Done()
				}
			}
		}
	}()
	for i := 0; i < b.N; i++ {
		go func() {
			c <- true
		}()
	}
	wg.Wait()
	time.Sleep(100 * time.Millisecond)
}

var i int
var m sync.Mutex
var c chan bool

func add() {
	i = i + 1
}
func addSync() {
	m.Lock()
	add()
	m.Unlock()
}
