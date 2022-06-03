package gl_logging

type GRPCLogger struct{}

func (g *GRPCLogger) Fatal(args ...interface{}) {
	error_(args...)
}

func (g *GRPCLogger) Fatalf(format string, args ...interface{}) {
	errorf(format, args...)
}

func (g *GRPCLogger) Fatalln(args ...interface{}) {
	errorln(args...)
}

func (g *GRPCLogger) Print(args ...interface{}) {
	info(args...)
}

func (g *GRPCLogger) Printf(format string, args ...interface{}) {
	infof(format, args...)
}

func (g *GRPCLogger) Println(args ...interface{}) {
	infoln(args...)
}
