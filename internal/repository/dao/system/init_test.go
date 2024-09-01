package system

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"testing"
)

func TestInit(t *testing.T) {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../dal", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3307)/dahuang"))
	if err != nil {
		panic(err)
	}
	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	g.ApplyBasic(g.GenerateAllTable()...)
	// Generate default DAO interface for those specified structs
	g.Execute()
}
