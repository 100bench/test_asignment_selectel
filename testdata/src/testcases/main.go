package testcases

import (
	"log"
	"log/slog"

	"go.uber.org/zap"
)

func valid() {
	log.Print("everything looks fine")
	log.Printf("processing item %d", 1)
	log.Println("done")
	slog.Info("server started")
	slog.Debug("request received")
	slog.Warn("retrying connection")
	slog.Error("connection failed")

	logger := zap.NewNop()
	logger.Info("server started")
	logger.Debug("handling request")

	zap.L().Info("server started")
	zap.L().Debug("handling request")
}

func badLowercase() {
	log.Print("Everything looks fine")   // want `log message must start with a lowercase letter`
	log.Printf("Processing item %d", 1) // want `log message must start with a lowercase letter`
	slog.Info("Server started")          // want `log message must start with a lowercase letter`
	slog.Error("Connection failed")      // want `log message must start with a lowercase letter`

	logger := zap.NewNop()
	logger.Warn("Something went wrong") // want `log message must start with a lowercase letter`

	zap.L().Info("Server started") // want `log message must start with a lowercase letter`
}

func badEnglish() {
	log.Print("всё в порядке")  // want `log message must be in English`
	slog.Info("сервер запущен") // want `log message must be in English`

	logger := zap.NewNop()
	logger.Error("ошибка подключения") // want `log message must be in English`
}

func badSpecialChars() {
	log.Print("everything looks fine!") // want `log message must not contain special characters or emoji`
	log.Print("is this working?")       // want `log message must not contain special characters or emoji`
	log.Print("wait for it...")         // want `log message must not contain special characters or emoji`
	log.Print("done;")                  // want `log message must not contain special characters or emoji`
	log.Print("starting:")              // want `log message must not contain special characters or emoji`
	slog.Info("hello 🎉")               // want `log message must not contain special characters or emoji`
}

func badSensitive() {
	log.Print("password: admin123")      // want `log message must not contain potentially sensitive data`
	slog.Warn("secret= hidden")         // want `log message must not contain potentially sensitive data`
	slog.Error("api_key=production_key") // want `log message must not contain potentially sensitive data`

	logger := zap.NewNop()
	logger.Info("token: abc123") // want `log message must not contain potentially sensitive data`
}

func badSensitiveConcat() {
	var password, apiKey, token string
	log.Print("user password: " + password) // want `log message must not contain potentially sensitive data`
	log.Print("api_key=" + apiKey)          // want `log message must not contain potentially sensitive data`
	slog.Info("token: " + token)            // want `log message must not contain potentially sensitive data`
}

func validSensitiveContext() {
	log.Print("user authenticated successfully")
	log.Print("token validated")
	log.Print("credential check passed")
	slog.Info("api request completed")
}
