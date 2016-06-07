package helper_test

import (
	"errors"

	"github.com/revinate/go-utils/helper"
	. "gopkg.in/check.v1"
)

type ErrorSuite struct{}

var _ = Suite(&ErrorSuite{})

func (s *ErrorSuite) TestErrors(c *C) {
	iolist := []struct {
		Err1, Err2, Err3 error
		Output           string
	}{
		{nil, nil, nil, ""},
		{errors.New("error1"), nil, nil, "error1"},
		{errors.New("error1"), errors.New("error2"), errors.New("error3"), "error1, error2, error3"},
		{errors.New("error1"), nil, errors.New("error2"), "error1, error2"},
	}
	for _, io := range iolist {
		c.Assert(helper.ErrorsToString(io.Err1, io.Err2, io.Err3), Equals, io.Output)
	}
}
