package main

import (
	"context"
	"database/sql"
	"fastduck/treasure-doc/service/mall/global"
	"flag"

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
	ctx := context.Background()
	//全局初始化
	global.GlobalInit(ctx, "cli")

	//同步写入日志
	defer global.Zap.Sync()
	defer global.ZapSugar.Sync()

	//关闭mysql
	db, _ := global.DbIns.DB()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	if *genOrm {
		global.ZapSugar.Info("genGormFile start!")
		genGormFile()
	}

	global.ZapSugar.Info("cli is end!")
}

func genGormFile() {
	g := gen.NewGenerator(gen.Config{
		OutPath: *genOrmPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode

	})

	g.UseDB(global.DbIns) // reuse your gorm db

	// Generate basic type-safe DAO API for struct `model.User` following conventions
	g.ApplyBasic(g.GenerateAllTable()...)

	// Generate the code
	g.Execute()
}
