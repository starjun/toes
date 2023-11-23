package global

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var (
	logger *zap.Logger
)

func InitLog(_log *Log) {
	mu.Lock()
	defer mu.Unlock()
	getLogger(_log)
}

func getLogger(_log *Log) {
	var level zapcore.Level
	switch _log.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "ts"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller_line"
	encoderConfig.FunctionKey = zapcore.OmitKey
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	//encoderConfig.LineEnding = "\r"
	encoderConfig.EncodeLevel = cEncodeLevel
	//encoderConfig.EncodeTime = cEncodeTime
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = cEncodeCaller

	/*
		LineEnding:     zapcore.DefaultLineEnding,     //输出的分割符
		EncodeLevel:    zapcore.LowercaseLevelEncoder, //序列化字符串的大小写
		//EncodeTime:          zapcore.ISO8601TimeEncoder,     //时间的编码格式
		EncodeTime:          EncodeTime,                     //时间自定义的
		EncodeDuration:      zapcore.SecondsDurationEncoder, //时间显示的位数
		EncodeCaller:        zapcore.ShortCallerEncoder,     //输出的运行文件路径长度
		EncodeName:          zapcore.FullNameEncoder,        //可选的
		NewReflectedEncoder: nil,
		ConsoleSeparator:    "", //控制台格式时，每个字段间的分割符,不配置默认即可
	*/
	var Encoder zapcore.Encoder
	if _log.Format == "console" {
		Encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		Encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	WriteSyncer := &lumberjack.Logger{
		Filename:   _log.Path,
		MaxSize:    300,
		MaxBackups: 3,
		MaxAge:     _log.Days,
	}

	writes := []zapcore.WriteSyncer{zapcore.AddSync(WriteSyncer)}
	if _log.Console {
		writes = append(writes, zapcore.AddSync(os.Stdout))
	}
	core := zapcore.NewCore(Encoder,
		zapcore.NewMultiWriteSyncer(writes...),
		level)

	//6.构造日志
	//设置为开发模式会记录panic
	development := zap.Development()
	//caller := zap.WithCaller(true)
	//构造一个字段
	//zap.Fields(zap.String("appName", "demozap"))
	//通过传入的配置实例化一个日志
	logger = zap.New(core, development, zap.AddCaller())

	// 替换全局 zap log
	zap.ReplaceGlobals(logger)
	// 全局使用 eg
	// zap.S().Info("hello")

	// 把标准库的 log.Logger 的 info 级别的输出重定向到 zap.Logger
	zap.RedirectStdLog(logger)
}

// cEncodeLevel 自定义日志级别显示
func cEncodeLevel(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

// cEncodeTime 自定义时间格式显示
func cEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + t.Format(LogTmFmt) + "]")
}

// cEncodeCaller 自定义行号显示
func cEncodeCaller(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + caller.TrimmedPath() + "]")
}

func LogSync() {
	logger.Sync()
}

// Debugw 输出 debug 级别的日志.
func LogDebugw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Debugw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Infow 输出 info 级别的日志.
func LogInfow(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Infow(msg, keysAndValues...)
	//defer logger.Sync()
}

// Warnw 输出 warning 级别的日志.
func LogWarnw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Warnw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Errorw 输出 error 级别的日志.
func LogErrorw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Errorw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Panicw 输出 panic 级别的日志.
func LogPanicw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Panicw(msg, keysAndValues...)
	//defer logger.Sync()
}

// Fatalw 输出 fatal 级别的日志.
func LogFatalw(msg string, keysAndValues ...interface{}) {
	logger.Sugar().Fatalw(msg, keysAndValues...)
	//defer logger.Sync()
}

// -- web中间件记录日志
func LogGin(ctx context.Context) *zap.Logger {

	_logger := logger.With(zap.Any("WEB", "GIN"))
	if requestID := ctx.Value(Cfg.Header.Requestid); requestID != nil {
		_logger = _logger.With(zap.Any("Traceid", requestID))
	}
	return _logger
}
