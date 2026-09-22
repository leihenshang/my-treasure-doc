package config

const GinModeRelease = "release"
const GinModeDev = "dev"

const WebPath = "web"
const FilesPath = "files"

type App struct {
	Host            string
	Port            int
	Name            string
	RunMode         string
	RegisterEnabled bool
	// TrustedProxies 信任的反代 IP/CIDR 列表（如 127.0.0.1、10.0.0.0/8）。
	// gin 仅在这些 IP 击中时才采信 X-Forwarded-For，从而让 ClientIP 返回真实客户端 IP
	// （IP 白名单 / 限流依赖它）。为空 = 不信任任何反代，直接取 RemoteAddr。
	TrustedProxies []string
}

func (app *App) IsRelease() bool {
	return app.RunMode == GinModeRelease
}

func (app *App) IsDev() bool {
	return app.RunMode == GinModeDev
}
