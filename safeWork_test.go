package single

import (
	"testing"
	"fmt"
	"time"
	"sync"
)

var i = 0
var ib = 0

func TestNewSWork(t *testing.T) {
	f := func(key string) ([]byte, error) {
		if key == "A" {
			i += 1
		} else {
			ib += 1
		}
		time.Sleep(time.Millisecond)
		return nil, nil
	}
	
	sw := NewSerialExecutor(f)
	key := "A"
	wg := sync.WaitGroup{}
	for a := 0; a < 20000; a++ {
		go func() {
			wg.Add(1)
			bs, err := sw.Get(key)
			_ = bs
			_ = err
			wg.Done()
		}()
		go func() {
			wg.Add(1)
			bs, err := sw.Get("B")
			_ = bs
			_ = err
			wg.Done()
		}()
	}
	wg.Wait()
	
	fmt.Println(i)
	fmt.Println(ib)
	v, ok := sw.m.Load(key)
	fmt.Println("ok : ", ok)
	fmt.Println("v  : ", v)
	time.Sleep(time.Hour)
}
