package single

import (
	"sync"
	"errors"
	"sync/atomic"
)

var (
	ErrInternal = errors.New("component internal error") // 组件内部出错
)

type Work struct {
	l        *sync.RWMutex
	result   []byte
	err      error
	reqCount int64
}

type SerialExecutor struct {
	m   sync.Map
	l   *sync.Mutex
	Fun func(key string) ([]byte, error)
}

func NewSerialExecutor(fun func(key string) ([]byte, error)) *SerialExecutor {
	return &SerialExecutor{sync.Map{}, &sync.Mutex{}, fun}
}

func (sw *SerialExecutor) Get(key string) ([]byte, error) {
	var work *Work
	var res []byte
	var err error
	func() {
		sw.l.Lock()
		defer sw.l.Unlock()
		if i, ok := sw.m.Load(key); ok {
			if v, ok1 := i.(*Work); ok1 {
				atomic.AddInt64(&v.reqCount, 1)
				work = v
			} else {
				work = nil
			}
		} else {
			l := &sync.RWMutex{}
			v := &Work{l, nil, nil, 1}
			sw.m.Store(key, v)
			work = v
		}
	}()
	
	if work == nil {
		sw.m.Delete(key)
		panic(ErrInternal)
	} else {
		func() {
			work.l.Lock()
			defer work.l.Unlock()
			
			// 当请求数量为0时删除字典键值
			defer func() {
				sw.l.Lock()
				defer sw.l.Unlock()
				if rc := atomic.LoadInt64(&work.reqCount); rc == 0 {
					sw.m.Delete(key)
				}
			}()
			
			if atomic.LoadInt64(&work.reqCount) > 1 {
				// 记录已经存在
				res = work.result
				err = work.err
				atomic.AddInt64(&work.reqCount, -1)
			} else {
				// 记录不存在 执行
				res, err = sw.Fun(key)
				work.result = res
				work.err = err
				atomic.AddInt64(&work.reqCount, -1)
			}
		}()
		return res, err
	}
}
