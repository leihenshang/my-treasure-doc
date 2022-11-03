package config

const MODE_RELEASE = "prod"
const MODE_DEV = "dev"

type App struct {
	Host    string
	Port    int
	Name    string
	RunMode string
}

func (app *App) IsRelease() bool {
	return app.RunMode == MODE_RELEASE
}

func (app *App) IsDev() bool {
	return app.RunMode == MODE_DEV
}
