package main

import (
	"os"
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

const LogPrefix = "[ZIWI] "

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
}
