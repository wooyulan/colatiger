package bootstrap

import (
	"colatiger/config"
	"colatiger/pkg/helper/path"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

func NewLog(conf *config.Configuration) (*zap.Logger, *lumberjack.Logger) {
	rootPath := path.RootPath()

	// 日志路径
	logFileDir := conf.Log.RootDir
	if !filepath.IsAbs(logFileDir) {
		logFileDir = filepath.Join(rootPath, logFileDir)
	}

	if ok, _ := path.PathExists(logFileDir); !ok {
		_ = os.Mkdir(logFileDir, os.ModePerm)
	}

	lv := conf.Log.Level

	var level zapcore.Level  // zap 日志等级
	var options []zap.Option // zap 配置项

	switch lv {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("2006-01-02 15:04:05.000"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(conf.App.Env + "." + l.String())
	}

	loggerWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logFileDir, conf.Log.Filename),
		MaxSize:    conf.Log.MaxSize,
		MaxBackups: conf.Log.MaxBackups,
		MaxAge:     conf.Log.MaxAge,
		Compress:   conf.Log.Compress,
	}

	var syncer zapcore.WriteSyncer
	// 判断是否控制台输出日志
	if conf.Log.LogInConsole {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(loggerWriter))
	} else {
		syncer = zapcore.AddSync(loggerWriter)
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		syncer,
		level,
	)

	logger := zap.New(core)

	// 判断是否显示代码行号
	if conf.Log.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return logger, loggerWriter

}
