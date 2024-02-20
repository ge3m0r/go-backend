package logger

type Logger interface {
	Debug(msg string, args ...Field)
	INFO(msg string, args ...Field)
	WARN(msg string, args ...Field)
	ERROR(msg string, args ...Field)
}

type Field struct {
	Key   string
	Value any
}
