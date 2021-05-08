package ratelimiter

type Logger interface {
	Errorf(format string, v ...interface{})
	Printf(format string, v ...interface{})
}

func SafeLog(logger Logger, do func(logger Logger)) {
	if logger == nil {
		return
	}
	do(logger)
}
