# Go 日志最佳实践：Zap 从入门到实战


{{< admonition type=abstract title="导语" open=true >}}
在现代微服务架构中，一个优秀的日志系统是保障应用可观测性的关键。Zap 作为 Go 生态中最受欢迎的日志库之一，以其卓越的性能和灵活的配置闻名。本文将带你深入了解 Zap 的实践应用，从基础配置到容器化环境下的最佳实践，帮助你构建一个既高效又易于维护的日志系统。无论是构建新项目还是优化现有系统，这都是一份不可或缺的实战指南。
{{< /admonition >}}

<!--more-->

如果你的应用采用容器化部署，其实更建议将日志输出到标准输出。
容器平台一般都具有采集容器日志的能力。
采集日志时，可以选择从标准输出采集或者容器中的日志文件采集，如果是从日志文件进行采集，通常需要配置日志采集路径，但如果是从标准输出采集，则不用。
所以，如果将日志直接输出到标准输出，则可以不加配置直接复用容器平台已有的能力，做到记录日志和采集日志完全解耦。

定制开发步骤分为以下几步：

创建一个封装了 zap.Logger 的自定义 Logger；

编写创建函数，创建 zapLogger 对象；

创建 *zap.Logger 对象；

实现日志接口。

## zap 核心结构

- `Logger`: 主日志记录器，用于写入日志
- `Core`: 日志的核心组件，定义日志的输出目标、格式和日志级别
- `Encoder`: 日志格式化器，决定日志的最终表现形式（如 `json` 或 `console` 格式）

  - `zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())`: JSON 编码器，结构化输出，适合生产环境
  - `zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())`: Console 编码器，人类可读格式，适合开发调试

- `Sync`: 用于确保日志缓冲区内容被写入，通过 `zapcore.AddSync` 设置日志写入目标，例如`zapcore.AddSync(os.Stdout)`、`zapcore.AddSync(file)`
- `Field`: 结构化日志的键值对，用于附加额外的上下文信息

## zap 基本使用

```go
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
```

- Development 模式: zap 输出 text 格式的日志，warn 和 error 带有战信息
- Example 模式: zap 输出 JSON 格式的日志，只有 `level` 和 `msg` 字段
- Production 模式: zap 输出 JSON 格式的日志，具有 `level`、`msg`、`ts`、`caller`、`stacktrace` 字段

## 设置日志级别

```go
func devWithConfig() {
 // 使用 zap 的 NewDevelopmentConfig 快速配置
 cfg := zap.NewDevelopmentConfig()
 cfg.Level = zap.NewAtomicLevelAt(zap.WarnLevel)

 // Create logger
 logger, _ := cfg.Build()
 logger.Info("info this is info")
 logger.Warn("warn this is warn")
 logger.Error("error this is error")
}
```

## 时间格式化

zap 日志库默认时间要么是带时区，要么就是时间戳，美观程度不足。

```Go
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
```

## 标准日志和 Sugar 日志

标准的 `*zap.Logger` 具有以下方法：

```Go
type Logger
    func L() *Logger
    func Must(logger *Logger, err error) *Logger
    func New(core zapcore.Core, options ...Option) *Logger
    func NewDevelopment(options ...Option) (*Logger, error)
    func NewExample(options ...Option) *Logger
    func NewNop() *Logger
    func NewProduction(options ...Option) (*Logger, error)
    func (log *Logger) Check(lvl zapcore.Level, msg string) *zapcore.CheckedEntry
    func (log *Logger) Core() zapcore.Core
    func (log *Logger) DPanic(msg string, fields ...Field)
    func (log *Logger) Debug(msg string, fields ...Field)
    func (log *Logger) Error(msg string, fields ...Field)
    func (log *Logger) Fatal(msg string, fields ...Field)
    func (log *Logger) Info(msg string, fields ...Field)
    func (log *Logger) Level() zapcore.Level
    func (log *Logger) Log(lvl zapcore.Level, msg string, fields ...Field)
    func (log *Logger) Name() string
    func (log *Logger) Named(s string) *Logger
    func (log *Logger) Panic(msg string, fields ...Field)
    func (log *Logger) Sugar() *SugaredLogger
    func (log *Logger) Sync() error
    func (log *Logger) Warn(msg string, fields ...Field)
    func (log *Logger) With(fields ...Field) *Logger
    func (log *Logger) WithLazy(fields ...Field) *Logger
    func (log *Logger) WithOptions(opts ...Option) *Logger
```

因此，可以使用 Sugar 方法得到一个加强版实例，常用方法如下：

```Go
type SugaredLogger
    func S() *SugaredLogger
    func (s *SugaredLogger) DPanic(args ...interface{})
    func (s *SugaredLogger) DPanicf(template string, args ...interface{})
    func (s *SugaredLogger) DPanicln(args ...interface{})
    func (s *SugaredLogger) DPanicw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Debug(args ...interface{})
    func (s *SugaredLogger) Debugf(template string, args ...interface{})
    func (s *SugaredLogger) Debugln(args ...interface{})
    func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Desugar() *Logger
    func (s *SugaredLogger) Error(args ...interface{})
    func (s *SugaredLogger) Errorf(template string, args ...interface{})
    func (s *SugaredLogger) Errorln(args ...interface{})
    func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Fatal(args ...interface{})
    func (s *SugaredLogger) Fatalf(template string, args ...interface{})
    func (s *SugaredLogger) Fatalln(args ...interface{})
    func (s *SugaredLogger) Fatalw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Info(args ...interface{})
    func (s *SugaredLogger) Infof(template string, args ...interface{})
    func (s *SugaredLogger) Infoln(args ...interface{})
    func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Level() zapcore.Level
    func (s *SugaredLogger) Log(lvl zapcore.Level, args ...interface{})
    func (s *SugaredLogger) Logf(lvl zapcore.Level, template string, args ...interface{})
    func (s *SugaredLogger) Logln(lvl zapcore.Level, args ...interface{})
    func (s *SugaredLogger) Logw(lvl zapcore.Level, msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Named(name string) *SugaredLogger
    func (s *SugaredLogger) Panic(args ...interface{})
    func (s *SugaredLogger) Panicf(template string, args ...interface{})
    func (s *SugaredLogger) Panicln(args ...interface{})
    func (s *SugaredLogger) Panicw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) Sync() error
    func (s *SugaredLogger) Warn(args ...interface{})
    func (s *SugaredLogger) Warnf(template string, args ...interface{})
    func (s *SugaredLogger) Warnln(args ...interface{})
    func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})
    func (s *SugaredLogger) With(args ...interface{}) *SugaredLogger
    func (s *SugaredLogger) WithLazy(args ...interface{}) *SugaredLogger
    func (s *SugaredLogger) WithOptions(opts ...Option) *SugaredLogger
```

## 结构化日志

zap 支持通过 `Field` 的形式记录结构化日志，方便分析和查阅

```Go
func devWithField() {
 logger, _ := zap.NewDevelopment()
 logger.Info(
  "This is info",
  zap.String("Name", "kyden"),
  zap.Int("Age", 18),
  zap.Bool("Cool", true),
 )
}
```

## 输出美化

`info`, `warn`，`error` 显示不同颜色，方便查阅

```Go
func devWithColor() {
 cfg := zap.NewDevelopmentConfig()
 cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
 cfg.EncoderConfig.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
  switch l{
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
```

## 日志前缀

如果一个项目中存在多个服务，可以使用前缀区分不同服务的日志，可以使用 `With` 方法来实现，示例如下：

```Go
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
```

## 全局日志

如果想在项目中的任何地方使用日志，那么就可以使用**全局日志**：

```Go
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
```

- `zap.L()`: 获取到的是标准的 Zap 示例
- `zap.S()`: 获取到的是 SugarZap 示例

## 日志双写

常见的：控制台和日志文件双写

方式一：使用 `zapcore.NewTee` 可以组合多个 core 示例

```Go
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

 // Usage
 zap.S().Infof("Global::This is info %s", "kytedance")
}
```

方式二：使用 `zapcore.NewMultiWriteSyncer`

```Go
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
```

## 日志切片

一般情况下，日志会根据`时间`和`日志等级`进行切片，将切好的日志放到单独的一个文件中，这样可以方便查阅

### 按时间分片

```Go
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
```

### 按日志等级分片

当需要根据日志等级进行分片时，通常是把`Error`等级的日志单独放到单独的文件中，而`Debug`、`Info`等等级的日志可以放到一个文件中，这样可以方便查阅

```Go
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
```

## 总结

```Go
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
```


---

> Author: [kyden](https://github.com/kydance)  
> URL: http://kydance.github.io/posts/golang-zap/  

