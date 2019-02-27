package single

import (
	"sync"
	"errors"
	"sync/atomic"
)

var (
	ErrInternal = errors.New("internal data error")
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
				// 拿到数据且数据格式正确
				work = v
			} else {
				// 拿到的数据格式不正确
				work = nil
			}
		} else {
			// 没有拿到数据，创建数据
			l := &sync.RWMutex{}
			v := &Work{l, nil, nil, 0}
			sw.m.Store(key, v)
			work = v
		}
	}()
	
	if work == nil {
		sw.m.Delete(key)
		return nil, ErrInternal
	} else {
		func() {
			atomic.AddInt64(&work.reqCount, 1)
			work.l.Lock()
			defer work.l.Unlock()
			
			// 当请求数量为0时删除字典键值
			defer func() {
				if work != nil {
					if rc := atomic.LoadInt64(&work.reqCount); rc == 0 {
						sw.l.Lock()
						defer sw.l.Unlock()
						sw.m.Delete(key)
					}
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
