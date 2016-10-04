package helper

import "sync"

type ConcurrentWork struct {
	workChan      chan interface{}
	out           chan interface{}
	numGoRoutines int
	mutex         sync.Mutex
}

func NewConcurrentWork(num int) ConcurrentWork {
	return ConcurrentWork{
		workChan:      make(chan interface{}, 0),
		numGoRoutines: num,
		mutex:         sync.Mutex{},
	}
}

func (c ConcurrentWork) Add(units []interface{}) ConcurrentWork {
	go func() {
		for _, unit := range units {
			c.workChan <- unit
		}
		close(c.workChan)
	}()
	return c
}

func (c ConcurrentWork) AddChan(ch chan interface{}) ConcurrentWork {
	go func() {
		for unit := range ch {
			c.workChan <- unit
		}
		close(c.workChan)
	}()
	return c
}

func (c ConcurrentWork) Do(callback func(interface{}) error) []error {
	return c.do(callback, false)
}

func (c ConcurrentWork) SafeDo(callback func(interface{}) error) []error {
	return c.do(callback, true)
}

func (c ConcurrentWork) do(callback func(interface{}) error, isSafe bool) []error {
	wg := sync.WaitGroup{}
	errors := []error{}
	for n := 0; n < c.numGoRoutines; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for unit := range c.workChan {
				if isSafe {
					c.mutex.Lock()
				}
				err := callback(unit)
				if isSafe {
					c.mutex.Unlock()
				}
				errors = append(errors, err)
			}
		}()
	}
	wg.Wait()
	return errors
}