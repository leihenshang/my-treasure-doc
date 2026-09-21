package global

import (
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"fastduck/treasure-doc/module/user/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

// stdoutSyncer 包装控制台输出，把 Sync 变成无副作用的空操作。
//
// Windows 下终端句柄不支持 FlushFileBuffers，zap 在进程退出时调用 Sync 会返回
// ERROR_INVALID_HANDLE（"The handle is invalid"）；而标准输出本来就不需要 fsync，
// 这个错误没有任何实际影响，吞掉它可以避免退出日志里出现两条看起来像故障的错误。
type stdoutSyncer struct{ writer io.Writer }

func (s stdoutSyncer) Write(p []byte) (int, error) { return s.writer.Write(p) }

func (s stdoutSyncer) Sync() error { return nil }

func initZapLogger() error {
	writeSyncyer, err := getLogWriter()
	if err != nil {
		return err
	}

	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncyer, zapcore.InfoLevel)

	var level zapcore.Level

	switch Conf.Log.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		Zap = zap.New(core, zap.AddStacktrace(level))
	} else {
		Zap = zap.New(core)
	}

	//显示行数
	Zap = Zap.WithOptions(zap.AddCaller())

	//Zap SUGAR
	Log = Zap.Sugar()

	return nil
}

func getEncoder() zapcore.Encoder {
	if Conf.Log.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConf())
	}
	return zapcore.NewConsoleEncoder(getEncoderConf())
}

func getEncoderConf() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	// switch {
	// case global.Conf.Zap.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
	// 	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	// case global.Conf.Zap.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
	// 	config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	// case global.Conf.Zap.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
	// 	config.EncodeLevel = zapcore.CapitalLevelEncoder
	// case global.Conf.Zap.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
	// 	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// default:
	// 	config.EncodeLevel = zapcore.LowercaseLevelEncoder
	// }
	return config
}

func getLogWriter() (zapcore.WriteSyncer, error) {
	logDir := Conf.Log.Directory
	var err error
	if isExists, _ := utils.PathExists(logDir); !isExists {
		fmt.Printf("创建日志目录 %v \n", logDir)
		err = os.Mkdir(logDir, os.ModePerm)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   path.Join(logDir, "app.log"),
		MaxSize:    2,  //最大2M
		MaxBackups: 5,  //最多保留5个备份
		MaxAge:     30, //最多保留30天
		Compress:   false,
	}

	if Conf.Log.ShowInConsole {
		return zapcore.NewMultiWriteSyncer(stdoutSyncer{writer: os.Stdout}, zapcore.AddSync(lumberJackLogger)), err
	}

	return zapcore.AddSync(lumberJackLogger), err
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 - 15:04:05.000"))
}
