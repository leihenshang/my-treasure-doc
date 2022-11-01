package main

import (
	"database/sql"
	"fastduck/treasure-doc/service/mall/global"
	"flag"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gen"
)

var genOrm = flag.Bool("gen", false, "gen gin orm")
var genOrmPath = flag.String("genOrmPath", "../data/query", "gen Orm Path")

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {

	flag.Parse()

	//全局初始化
	global.GlobalInit()

	//同步写入日志
	defer global.ZAP.Sync()
	defer global.ZAPSUGAR.Sync()

	//关闭mysql
	db, _ := global.DB.DB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	if *genOrm {
		global.ZAPSUGAR.Info("genGormFile start!")
		genGormFile()
	}

	global.ZAPSUGAR.Info("cli is end!")
}

func genGormFile() {
	//设置运行模式
	if global.CONFIG.App.IsRelease() {
		fmt.Println("设置模式为", gin.ReleaseMode)
		gin.SetMode(gin.ReleaseMode)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: *genOrmPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(global.DB) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(g.GenerateAllTable()...)

	// Generate the code
	g.Execute()
}
