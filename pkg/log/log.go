package log

import (
	"colatiger/pkg/helper/path"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

type Logger struct {
	*zap.Logger
}

func NewLog(conf *viper.Viper) *Logger {
	return initZap(conf)
}

func initZap(conf *viper.Viper) *Logger {
	rootPath := path.RootPath()

	// 日志路径
	logFileDir := conf.GetString("log.root_dir")
	if !filepath.IsAbs(logFileDir) {
		logFileDir = filepath.Join(rootPath, logFileDir)
	}

	if ok, _ := path.PathExists(logFileDir); !ok {
		_ = os.Mkdir(logFileDir, os.ModePerm)
	}

	lv := conf.GetString("log.level")

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
		encoder.AppendString(conf.GetString("app.env") + "." + l.String())
	}

	loggerWriter := &lumberjack.Logger{
		Filename:   filepath.Join(logFileDir, conf.GetString("log.filename")),
		MaxSize:    conf.GetInt("log.max_size"),
		MaxBackups: conf.GetInt("log.max_backups"),
		MaxAge:     conf.GetInt("log.max_age"),
		Compress:   conf.GetBool("log.compress"),
	}

	return &Logger{
		zap.New(
			zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig),
				zapcore.AddSync(loggerWriter), level), options...),
	}
}
