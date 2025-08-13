package logger

type Level int8

const (
	Debug Level = iota
	Info
	Warn
	Error
	Fatal
)

func (level Level) Key() string {
	return "level"
}

func (level Level) String() string {
	switch level {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return ""
	}
}
