package logger

import (
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
)

var (
	ErrorJsonMarshal   = "Error Json Marshal"
	ErrorJsonUnmarshal = "Error Json Unmarshal"

	ErrorSendMessage = "Error send message to tg"
	ErrorAtoi        = "Error parse string to int"
)

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции InitLogger, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
type Logger struct {
	*zap.Logger
	originalError error
}

var Log = &Logger{
	Logger: zap.NewNop()}

func (l *Logger) AddOriginalError(err error) CustomErorrer {
	l.originalError = err
	return l
}

type CustomErorrer interface {
	CustomInfo(msg string, fields map[string]interface{})
	CustomError(msg string, fields map[string]interface{})
	CustomWarn(msg string, fields map[string]interface{})
}

func (l *Logger) CustomInfo(msg string, fields map[string]interface{}) {

	field := make([]zapcore.Field, 0)
	field = append(field, zap.Any("fields", fields))
	if l.originalError != nil {
		field = append(field, zap.Error(l.originalError))
	}
	l.Info(msg, field...)
}

func (l *Logger) CustomError(msg string, fields map[string]interface{}) {

	field := make([]zapcore.Field, 0)
	field = append(field, zap.Any("fields", fields))
	if l.originalError != nil {
		field = append(field, zap.Error(l.originalError))
	}
	l.Error(msg, field...)
}

func (l *Logger) CustomWarn(msg string, fields map[string]interface{}) {

	field := make([]zapcore.Field, 0)
	field = append(field, zap.Any("fields", fields))
	if l.originalError != nil {
		field = append(field, zap.Error(l.originalError))
	}
	l.Warn(msg, field...)
}

func (l *Logger) HttpError(statusCode int, msg string, fields ...zapcore.Field) {
	field := append(fields, zap.Int("statusCode", statusCode))
	l.Error(msg, field...)
}

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()

	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	// устанавливаем синглтон
	Log.Logger = zl

	return nil
}

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize2(level string, logFilePath string) error {
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()

	// устанавливаем уровень
	cfg.Level = lvl

	//
	fileEnconder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	fileEnconder.AddString("idLogger", uuid.New().String())

	core := zapcore.NewCore(fileEnconder, zapcore.AddSync(file), lvl)

	// создаём логер на основе конфигурации
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	//создаем мультилог
	logger := zap.New(zapcore.NewTee(zl.Core(), core), zap.AddCaller())

	// устанавливаем синглтон
	Log.Logger = logger

	return nil
}

//	func Initialize(level string, logFilePath string) error {
//		// ... (existing code)
//
//		// создаем конфигурацию для записи в файл
//		fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
//		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
//		if err != nil {
//			return err
//		}
//		fileCore := zapcore.NewCore(fileEncoder, zapcore.AddSync(file), lvl)
//
//		// создаем логер для записи в файл
//		fileLogger := zap.New(fileCore, zap.AddCaller())
//
//		// устанавливаем мультилоггер для записи в несколько мест
//		multiLogger := zap.New(zapcore.NewTee(zl.Core(), fileCore), zap.AddCaller())
//
//		// устанавливаем синглтон
//		Log = multiLogger
//
//		return nil
//	}
//
// RequestLogger — middleware-логер для входящих HTTP-запросов.
func RequestLogger(h http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("got incoming HTTP request",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
		)
		h(w, r)
	})
}
