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
}

func (app *App) IsRelease() bool {
	return app.RunMode == GinModeRelease
}

func (app *App) IsDev() bool {
	return app.RunMode == GinModeDev
}
