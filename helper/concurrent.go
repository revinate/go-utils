package helper

import "sync"

type ConcurrentWork struct {
	workChan      chan interface{}
	out           chan interface{}
	numGoRoutines int
}

func NewConcurrentWork(num int) ConcurrentWork {
	return ConcurrentWork{
		workChan:      make(chan interface{}, 0),
		numGoRoutines: num,
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
	wg := sync.WaitGroup{}
	errors := []error{}
	for n := 0; n < c.numGoRoutines; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for unit := range c.workChan {
				err := callback(unit)
				errors = append(errors, err)
			}
		}()
	}
	wg.Wait()
	return errors
}
