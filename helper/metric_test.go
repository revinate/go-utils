package helper_test

import (
	"testing"

	"github.com/revinate/go-utils/helper"

	"time"

	. "gopkg.in/check.v1"
)

type MetricSuite struct{}

func Test(t *testing.T) {
	TestingT(t)
}

var (
	_ = Suite(&MetricSuite{})
)

func (s *MetricSuite) TestMetricIncr(c *C) {
	metric := helper.NewMetric("test", 3)
	metric.Incr(1)
	metric.Incr(-1)
	metric.Incr(1)

	c.Assert(metric.GetPerSec().Total(), Equals, 1)
	c.Assert(metric.GetPerMin().Total(), Equals, 1)
	c.Assert(metric.GetPerHour().Total(), Equals, 1)
}

func (s *MetricSuite) TestMetricTotal(c *C) {
	metric := helper.NewMetric("test", 3)
	for i := 0; i < 10; i++ {
		metric.Incr(1)
	}

	c.Assert(metric.GetPerSec().Total(), Equals, 10)
	c.Assert(metric.GetPerMin().Total(), Equals, 10)
	c.Assert(metric.GetPerHour().Total(), Equals, 10)
}

func (s *MetricSuite) TestNewMetric(c *C) {
	metric := helper.NewMetric("test", 3)
	c.Assert(metric, NotNil)
	c.Assert(metric.GetPerMin().GetSize(), Equals, 3)
	c.Assert(metric.GetPerSec().GetSize(), Equals, 3)
	c.Assert(metric.GetPerHour().GetSize(), Equals, 3)
}

func (s *MetricSuite) TestCleanOld(c *C) {
	metric := helper.NewMetric("test", 3)
	for i := 0; i < 5; i++ {
		metric.Incr(1)
		time.Sleep(time.Second)
	}

	c.Assert(len(metric.GetPerSec().GetWindows()), Equals, 3)
	c.Assert(len(metric.GetPerMin().GetWindows()) >= 1 && len(metric.GetPerMin().GetWindows()) <= 2, Equals, true)
	c.Assert(len(metric.GetPerHour().GetWindows()) >= 1 && len(metric.GetPerHour().GetWindows()) <= 2, Equals, true)
}

func (s *MetricSuite) TestMetricPublicView(c *C) {
	metric := helper.NewMetric("test", 1)
	metric.Incr(1)

	public, ok := metric.GetPublic().(map[string]interface{})
	c.Assert(ok, Equals, true)
	c.Assert(public["name"], Equals, "test")

	perSec, ok := public["perSec"].(map[string]interface{})
	c.Assert(ok, Equals, true)
	windows, ok := perSec["windows"].(map[string]*helper.Window)
	c.Assert(ok, Equals, true)
	c.Assert(len(windows), Equals, 1)
}
