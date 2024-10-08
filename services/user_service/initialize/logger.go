package initialize

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/zhaohaihang/user_service/global"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var dest io.Writer
var logFileWriter io.Writer

// InitLogger
func InitLogger() {
	writeSyncer := getLogWriter()
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.InfoLevel)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)
	zap.S().Infow("logger init success")
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter() zapcore.WriteSyncer {
	logFileWriter = &lumberjack.Logger{
		Filename:   createLogFileName(),
		MaxSize:    1,
		MaxAge:     5,
		MaxBackups: 30,
	}
	dest = io.MultiWriter(logFileWriter, os.Stdout)
	return zapcore.AddSync(dest)
}

func createLogFileName() string {
	today := time.Now().Format("2006-01-02")
	return fmt.Sprintf("%s/%s.log", global.FilePath.LogFile, today)
}
