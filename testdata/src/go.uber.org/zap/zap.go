package zap

type Logger struct{}

type Field struct{}

func NewNop() *Logger { return &Logger{} }
func L() *Logger      { return &Logger{} }

func (l *Logger) Debug(msg string, fields ...Field)  {}
func (l *Logger) Info(msg string, fields ...Field)   {}
func (l *Logger) Warn(msg string, fields ...Field)   {}
func (l *Logger) Error(msg string, fields ...Field)  {}
func (l *Logger) Fatal(msg string, fields ...Field)  {}
func (l *Logger) Panic(msg string, fields ...Field)  {}
func (l *Logger) DPanic(msg string, fields ...Field) {}

func String(key, val string) Field { return Field{} }
