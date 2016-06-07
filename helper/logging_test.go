package helper_test

import . "gopkg.in/check.v1"

type LoggingSuite struct{}

var (
	_ = Suite(&LoggingSuite{})
)

type IOStruct struct {
	Func   func(string)
	Input  string
	Output string
}

//func (s *LoggingSuite) TestLogging(c *C) {
//	iolist := []IOStruct{
//		{
//			func(str string) { helper.Info(str) },
//			"Hello World",
//			"INFO[0005] " + time.Now().Format(time.StampMilli) + ": Hello World",
//		},
//		{
//			func(str string) { helper.Error(errors.New("error"), str) },
//			"Hello World",
//			"ERRO[0005] " + time.Now().Format(time.StampMilli) + ": Hello World, Error: error",
//		},
//		{
//			func(str string) { helper.Debug(str, "test obj") },
//			"Hello World",
//			"DEBU[0005] " + time.Now().Format(time.StampMilli) + ": Hello World, \"test obj\"",
//		},
//	}
//	for _, io := range iolist {
//		outout := tests.CaptureStdout(func() { io.Func(io.Input) })
//		c.Assert(output, Equals, io.Output)
//	}
//}
