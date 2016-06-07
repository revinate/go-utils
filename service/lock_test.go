package service_test

import (
	"testing"
	"time"

	"golang.org/x/net/context"

	. "gopkg.in/check.v1"
	"github.com/revinate/go-utils/service"
)

const LockName = "test"

type LockSuite struct{}

var (
	_ = Suite(&LockSuite{})
)

func Test(t *testing.T) {
	TestingT(t)
}

func getLockService() (*service.LockService, error) {
	return service.NewLockService("127.0.0.1:2379", context.Background())
}

func (s *LockSuite) TestAcquireLockSuccess(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestAcquireLockFail(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	_, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(err, NotNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestReleaseLockSuccess(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	_, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(err, NotNil)

	err = serv.ReleaseLock(LockName, value)
	c.Assert(err, IsNil)
}

func (s *LockSuite) TestReacquireLockSuccess(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	_, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(err, NotNil)

	err = serv.ReleaseLock(LockName, value)
	c.Assert(err, IsNil)

	value, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestReleaseLockFail(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	_, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(err, NotNil)

	err = serv.ReleaseLock(LockName, value+"2")
	c.Assert(err, NotNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestAcquireLockAfterTTLSuccess(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	_, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(err, NotNil)

	time.Sleep(time.Second * 2)

	value, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestRenewLockSuccess(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second*2, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	time.Sleep(time.Second)

	err = serv.RenewLock(LockName, value, time.Second*2)
	c.Assert(err, IsNil)

	time.Sleep(time.Second)

	value, err = serv.AcquireLock(LockName, time.Second, false)
	c.Assert(len(value), Equals, 0)
	c.Assert(err, NotNil)

	serv.ReleaseLock(LockName, value)
}

func (s *LockSuite) TestRenewLockFail(c *C) {
	serv, serr := getLockService()
	c.Assert(serr, IsNil)
	value, err := serv.AcquireLock(LockName, time.Second*2, false)
	c.Assert(len(value), Not(Equals), 0)
	c.Assert(err, IsNil)

	time.Sleep(time.Second)

	err = serv.RenewLock(LockName, value+"2", time.Second)
	c.Assert(err, NotNil)

	serv.ReleaseLock(LockName, value)
}
