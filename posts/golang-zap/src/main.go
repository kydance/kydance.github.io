package main

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// Custom color
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

func dev() {
	logger, _ := zap.NewDevelopment()
	logger.Info("dev this is info")
	logger.Warn("dev this is warn")
	logger.Error("dev this error")
}

func test() {
	logger := zap.NewExample()
	logger.Info("test this is info")
	logger.Warn("test this is warn")
	logger.Error("test this is error")
}

func prod() {
	logger, _ := zap.NewProduction()
	logger.Info("prod this is info")
	logger.Warn("prod this is warn")
	logger.Error("prod this error")
}

func devWithConfig() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)

	logger, _ := cfg.Build()
	logger.Info("info this is info")
	logger.Warn("warn this is warn")
	logger.Error("error this is error")
}

func devWithConfigTimeFormat() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger.Info("info this is info")
	logger.Warn("warn this is warn")
	logger.Error("error this is error")
}

func devWithField() {
	logger, _ := zap.NewDevelopment()
	logger.Info(
		"This is info",
		zap.String("Name", "kyden"),
		zap.Int("Age", 18),
		zap.Bool("Cool", true),
	)
}

func devWithColor() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString(ColorBlue + level.CapitalString() + ColorReset)
		case zapcore.InfoLevel:
			enc.AppendString(ColorGreen + level.CapitalString() + ColorReset)
		case zapcore.WarnLevel:
			enc.AppendString(ColorYellow + level.CapitalString() + ColorReset)
		case zapcore.ErrorLevel:
			enc.AppendString(ColorRed + level.CapitalString() + ColorReset)
		default:
			enc.AppendString(level.String()) // default behavior
		}
	}
	logger, _ := cfg.Build()

	logger.Info("Dev::This is info")
	logger.Warn("Dev::This is warn")
	logger.Error("Dev::This is error")
}

type PrefixEncoder struct{ zapcore.Encoder }

func (enc *PrefixEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	// 先调用原始的 EncodeEntry 方法生成日志行
	buf, err := enc.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// 将日志行前缀添加到日志行
	logLine := buf.String()
	buf.Reset()
	buf.AppendString(LogPrefix + logLine)

	return buf, nil
}

func devWithCustomEncoder() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString(ColorBlue + level.String() + ColorReset)
		case zapcore.InfoLevel:
			enc.AppendString(ColorGreen + level.String() + ColorReset)
		case zapcore.WarnLevel:
			enc.AppendString(ColorYellow + level.String() + ColorReset)
		case zapcore.ErrorLevel:
			enc.AppendString(ColorRed + level.String() + ColorReset)
		default:
			enc.AppendString(level.String()) // default behavior
		}
	}

	// Create a custom encoder
	encoder := &PrefixEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig), // Use the default console encoder
	}

	// Create Core
	core := zapcore.NewCore(
		encoder,                    // Custom encoder
		zapcore.AddSync(os.Stdout), // Write logs to stdout
		zapcore.DebugLevel,         // Set log level
	)

	// Create Logger
	logger := zap.New(core, zap.AddCaller())

	logger.Info("Dev::This is info")
	logger.Warn("Dev::This is warn")
	logger.Error("Dev::This is error")
}

// 初始化全局日志
func initLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString(ColorBlue + level.String() + ColorReset)
		case zapcore.InfoLevel:
			enc.AppendString(ColorGreen + level.String() + ColorReset)
		case zapcore.WarnLevel:
			enc.AppendString(ColorYellow + level.String() + ColorReset)
		case zapcore.ErrorLevel:
			enc.AppendString(ColorRed + level.String() + ColorReset)
		default:
			enc.AppendString(level.String()) // default behavior
		}
	}

	logger, _ := cfg.Build()
	zap.ReplaceGlobals(logger)
}

func devWithGlobal() {
	zap.L().Info("Global::This is info")
	zap.L().Warn("Global::This is warn")
	zap.L().Error("Global::This is error")

	zap.S().Infof("Global::This is info %s", "kytedance")
	zap.S().Warnf("Global::This is warn %v", "kytedance")
	zap.S().Errorf("Global::This is error %+v", "kytedance")
}

func InitLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")

	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)

	file, _ := os.OpenFile("ziwi.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	fileCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.AddSync(file),
		zapcore.DebugLevel,
	)

	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(logger)

	zap.S().Infof("Global::This is info %s", "kytedance")
}

func ZapMultiWriteSyncer() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	file, _ := os.OpenFile("ziwi.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	writeSyncer := zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(os.Stdout),
		zapcore.AddSync(file),
	)

	zap.ReplaceGlobals(zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		writeSyncer,
		zapcore.DebugLevel,
	), zap.AddCaller()))

	// Use
	zap.S().Infof("Global::This is info %s", "kytedance")
}

// --- 日志切片 ---

type DynamicLogWriter struct {
	mtx     sync.Mutex //
	currDay string     // current day
	file    *os.File
	logDir  string
}

func (w *DynamicLogWriter) Write(b []byte) (n int, err error) {
	w.mtx.Lock()
	defer w.mtx.Unlock()

	// Check: If the current day has changed, create a new file
	today := time.Now().Format(time.DateOnly)
	if today != w.currDay {
		if w.file != nil { // Close non-nil log file
			w.file.Close()
		}

		// Create new log file
		if err := os.MkdirAll(w.logDir, 0755); err != nil {
			return 0, err
		}
		filePath := path.Join(w.logDir, "ziwi-"+today+".log")
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}

		// Update log writer
		w.file = file
		w.currDay = today
	}

	// Write log into file
	return w.file.Write(b)
}

func InitGolbalLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)
	cfg.EncoderConfig.EncodeLevel = func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch level {
		case zapcore.DebugLevel:
			enc.AppendString(ColorBlue + level.String() + ColorReset)
		case zapcore.InfoLevel:
			enc.AppendString(ColorGreen + level.String() + ColorReset)
		case zapcore.WarnLevel:
			enc.AppendString(ColorYellow + level.String() + ColorReset)
		case zapcore.ErrorLevel:
			enc.AppendString(ColorRed + level.String() + ColorReset)
		default:
			enc.AppendString(level.String()) // default behavior
		}
	}

	logger := zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(&DynamicLogWriter{
				logDir: "logs",
			}),

			os.Stdout,
		),
		zapcore.DebugLevel,
	), zap.AddCaller())
	zap.ReplaceGlobals(logger)
}

func SliceLog() {
	InitGolbalLogger()
	zap.S().Infof("Global::SliceLog::This is info %s", "kytedance")
	zap.S().Debugf("Global::SliceLog::This is debug %s", "kytedance")
	zap.S().Warnf("Global::SliceLog::This is warn %s", "kytedance")
	zap.S().Errorf("Global::SliceLog::This is error %s", "kytedance")
}

// --- 日志切片：根据 level ---
type LevelEncoder struct {
	zapcore.Encoder

	errFile *os.File // Error Log File
}

func (e *LevelEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	//
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	switch entry.Level {
	case zapcore.ErrorLevel:
		if e.errFile == nil {
			file, err := os.OpenFile(path.Join("logs", "err.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
			if err != nil {
				return nil, fmt.Errorf("open error log file: %w", err)
			}

			e.errFile = file
		}

		_, _ = e.errFile.WriteString(buf.String())
	}

	return buf, nil
}

func InitGolbalLevelLogger() {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime)

	logger := zap.New(zapcore.NewCore(
		&LevelEncoder{
			Encoder: zapcore.NewConsoleEncoder(cfg.EncoderConfig),
		},
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	), zap.AddCaller())

	zap.ReplaceGlobals(logger)
}

func SliceLevelLog() {
	InitGolbalLevelLogger()

	zap.S().Infof("Global::SliceLevelLog::This is info %s", "kytedance")
	zap.S().Debugf("Global::SliceLevelLog::This is debug %s", "kytedance")
	zap.S().Warnf("Global::SliceLevelLog::This is warn %s", "kytedance")
	zap.S().Errorf("Global::SliceLevelLog::This is error %s", "kytedance")
}
