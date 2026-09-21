package global

import (
	"testing"

	"fastduck/treasure-doc/module/user/config"
)

// 控制台输出场景下 Sync 必须保持安静。
//
// Windows 终端句柄不支持 FlushFileBuffers，直接 zapcore.AddSync(os.Stdout) 会让
// 进程退出时的 Zap.Sync()/Log.Sync() 报 "sync /dev/stdout: The handle is invalid"，
// 看起来像启动或重置密码失败，实际没有任何影响。
func TestLogWriterSyncIsQuietForConsoleOutput(t *testing.T) {
	prev := Conf
	Conf = &config.Config{Log: config.Log{
		Directory:     t.TempDir(),
		ShowInConsole: true,
		Format:        "console",
		Level:         "info",
	}}
	t.Cleanup(func() { Conf = prev })

	writer, err := getLogWriter()
	if err != nil {
		t.Fatalf("getLogWriter 失败: %v", err)
	}
	if err := writer.Sync(); err != nil {
		t.Fatalf("控制台输出下的 Sync 不应报错: %v", err)
	}
}

// 仅写文件时同样不应因为控制台句柄产生错误。
func TestLogWriterSyncIsQuietWithoutConsole(t *testing.T) {
	prev := Conf
	Conf = &config.Config{Log: config.Log{
		Directory:     t.TempDir(),
		ShowInConsole: false,
		Format:        "console",
		Level:         "info",
	}}
	t.Cleanup(func() { Conf = prev })

	writer, err := getLogWriter()
	if err != nil {
		t.Fatalf("getLogWriter 失败: %v", err)
	}
	if err := writer.Sync(); err != nil {
		t.Fatalf("文件输出的 Sync 不应报错: %v", err)
	}
}
