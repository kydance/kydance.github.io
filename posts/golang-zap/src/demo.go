package main

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	BuleColor   = "\033[34m"
	YellowColor = "\033[33m"
	GreenColor  = "\033[32m"
	RedColor    = "\033[31m"
	ResetColor  = "\033[0m"

	LogPrefix = "[ZIWI] "
)

type LogEncoder struct {
	zapcore.Encoder

	logDir      string
	file        *os.File // 普通日志文件，包含 error
	errFile     *os.File // error 日志文件
	currentDate string
}

// CusEncodeLevel beautify the log
func CusEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch level {
	case zapcore.DebugLevel:
		enc.AppendString(BuleColor + level.String() + ResetColor)
	case zapcore.InfoLevel:
		enc.AppendString(GreenColor + level.String() + ResetColor)
	case zapcore.WarnLevel:
		enc.AppendString(YellowColor + level.String() + ResetColor)
	case zapcore.ErrorLevel:
		enc.AppendString(RedColor + level.String() + ResetColor)
	default:
		enc.AppendString(level.String()) // default behavior
	}
}

// EncodeEntry add log prefix, split error log and split log files by date

func (e *LogEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, fmt.Errorf("EncodeEntry error: %w", err)
	}

	data := buf.String()
	buf.Reset()
	buf.AppendString(LogPrefix + data)
	data = buf.String()

	// split log files by date
	now := time.Now().Format(time.DateOnly)
	if e.currentDate != now {
		if err := os.MkdirAll(e.logDir, 0o666); err != nil {
			return nil, fmt.Errorf("Create log dir error: %w", err)
		}

		fileName := now + ".log"
		file, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
		if err != nil {
			return nil, fmt.Errorf("Create log file error: %w", err)
		}
		e.file = file
		e.currentDate = now
	}

	// split error log
	switch entry.Level {
	case zapcore.ErrorLevel:
		if e.errFile == nil {
			file, err := os.OpenFile(now+"_err.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o666)
			if err != nil {
				return nil, fmt.Errorf("Create err log file error: %w", err)
			}

			e.errFile = file
		}

		_, _ = e.errFile.WriteString(buf.String())
	}

	if e.currentDate == now {
		_, _ = e.file.WriteString(data)
	}

	return buf, nil
}

func InitGolbalLog(logdir string) *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = CusEncodeLevel

	logger := zap.New(
		zapcore.NewCore(
			&LogEncoder{
				Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig),
				logDir:  logdir,
			},
			zapcore.AddSync(os.Stdout),
			zapcore.DebugLevel,
		),
		zap.AddCaller(),
	)
	zap.ReplaceGlobals(logger)

	return logger
}

func main() {
	logger := InitGolbalLog("logs")
	logger.Sugar().Infof("Global::This is info %s", "kytedance")
	logger.Warn("Warn log test case")
	zap.L().Error("Error log test case")
	zap.S().Debugln("Panic log test case")
}
