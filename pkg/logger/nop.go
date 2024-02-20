package logger

type NopLogger struct {
}

func (n *NopLogger) Debug(msg string, args ...Field) {

}

func (n *NopLogger) INFO(msg string, args ...Field) {

}

func (n *NopLogger) WARN(msg string, args ...Field) {

}

func (n *NopLogger) ERROR(msg string, args ...Field) {

}
