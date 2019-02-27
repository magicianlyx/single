package single
//
// import (
// 	"sync"
// 	"errors"
// 	"sync/atomic"
// )
//
// var (
// 	ErrInternal = errors.New("internal data error")
// )
//
// type Work struct {
// 	l *sync.RWMutex
// 	//
// 	result   []byte
// 	err      error
// 	reqCount int64
// }
//
// type Swork struct {
// 	m sync.Map
// }
//
// func NewSWork() *Swork {
// 	return &Swork{sync.Map{}}
// }
//
// func (sw *Swork) Get(key string, fun func(key string) ([]byte, error)) ([]byte, error) {
//
// 	if i, ok := sw.m.Load(key); ok {
// 		if v, ok1 := i.(*Work); ok1 {
// 			// 拿到数据且数据格式正确
// 			defer func() {
// 				if v != nil {
// 					if atomic.LoadInt64(&v.reqCount) == 0 {
// 						sw.m.Delete(key)
// 					}
// 				}
// 			}()
//
// 			v.l.Lock()
// 			atomic.AddInt64(&v.reqCount, 1)
// 			res := v.result
// 			err := v.err
// 			atomic.AddInt64(&v.reqCount, -1)
// 			v.l.Unlock()
// 			return res, err
// 		} else {
// 			sw.m.Delete(key)
// 			return nil, ErrInternal
// 		}
// 	} else {
// 		// 没有拿到数据
// 		l := &sync.RWMutex{}
// 		v := &Work{l, nil, nil, 0}
// 		sw.m.Store(key, v)
//
// 		defer func() {
// 			if v != nil {
// 				if atomic.LoadInt64(&v.reqCount) == 0 {
// 					sw.m.Delete(key)
// 				}
// 			}
// 		}()
// 		return func() ([]byte, error) {
// 			l.Lock()
// 			defer l.Unlock()
// 			if atomic.AddInt64(&v.reqCount, 1) >= 2 {
//
// 			}
// 			res, err := fun(key)
// 			v.result = res
// 			v.err = err
// 			atomic.AddInt64(&v.reqCount, -1)
// 			return res, err
// 		}()
// 	}
// }
