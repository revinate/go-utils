package helper

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Window struct {
	Count     *int64
	timestamp int64
}

type Count struct {
	sync.Mutex
	size     int
	duration time.Duration
	windows  map[string]*Window
}

type Metric struct {
	name    string `json:"name"`
	perSec  *Count `json:"per_sec"`
	perMin  *Count `json:"per_min"`
	perHour *Count `json:"per_hour"`
}

func (m *Metric) GetPublic() interface{} {
	return map[string]interface{}{
		"name":    m.name,
		"perSec":  m.perSec.GetPublic(),
		"perMin":  m.perMin.GetPublic(),
		"perHour": m.perHour.GetPublic(),
	}
}

func NewMetric(name string, size int) *Metric {
	return &Metric{
		name:    name,
		perSec:  &Count{windows: map[string]*Window{}, size: size, duration: time.Second},
		perMin:  &Count{windows: map[string]*Window{}, size: size, duration: time.Minute},
		perHour: &Count{windows: map[string]*Window{}, size: size, duration: time.Hour},
	}
}

func (m Metric) Incr(by int64) {
	now := time.Now()
	secKey := fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	minKey := fmt.Sprintf("%d-%d-%d %d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	hourKey := fmt.Sprintf("%d-%d-%d %d", now.Year(), now.Month(), now.Day(), now.Hour())
	m.perSec.Incr(now, secKey, by)
	m.perMin.Incr(now, minKey, by)
	m.perHour.Incr(now, hourKey, by)
}

func (m *Metric) GetPerSec() *Count {
	return m.perSec
}

func (m *Metric) GetPerMin() *Count {
	return m.perMin
}

func (m *Metric) GetPerHour() *Count {
	return m.perHour
}

func (c *Count) Incr(now time.Time, key string, by int64) {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.windows[key]; !ok {
		var init int64 = 0
		c.windows[key] = &Window{Count: &init, timestamp: now.Unix()}
	}
	atomic.AddInt64(c.windows[key].Count, by)
	c.cleanOld()
}

func (c *Count) GetCount(now time.Time) int64 {
	key := ""
	if c.duration == time.Second {
		key = fmt.Sprintf("%d-%d-%d %d:%d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	} else if c.duration == time.Minute {
		key = fmt.Sprintf("%d-%d-%d %d:%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	} else if c.duration == time.Hour {
		key = fmt.Sprintf("%d-%d-%d %d", now.Year(), now.Month(), now.Day(), now.Hour())
	}
	if window, ok := c.windows[key]; ok {
		return *window.Count
	}
	return 0
}

func (c *Count) cleanOld() {
	timeLimit := time.Now().Add(-c.duration * time.Duration(c.size)).Unix()
	for key, val := range c.windows {
		if val.timestamp <= timeLimit {
			delete(c.windows, key)
		}
	}
}

func (c *Count) GetPublic() interface{} {
	return map[string]interface{}{"windows": c.windows}
}

func (c *Count) Total() int {
	var total int64 = 0
	for _, w := range c.windows {
		total += *w.Count
	}
	return int(total)
}

func (c *Count) GetSize() int {
	return c.size
}

func (c *Count) GetWindows() map[string]*Window {
	return c.windows
}
