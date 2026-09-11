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
	// CorsAllowOrigins 允许跨域访问后台接口的来源白名单。
	// 留空表示仅允许同源访问；开发态前端独立运行在其它端口时需把其地址加入。
	CorsAllowOrigins []string
}

func (app *App) IsRelease() bool {
	return app.RunMode == GinModeRelease
}

func (app *App) IsDev() bool {
	return app.RunMode == GinModeDev
}
