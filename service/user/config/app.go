package config

const GinModeRelease = "release"
const GinModeDev = "dev"

type App struct {
	Host    string
	Port    int
	Name    string
	RunMode string
}

func (app *App) IsRelease() bool {
	return app.RunMode == GinModeRelease
}

func (app *App) IsDev() bool {
	return app.RunMode == GinModeDev
}
