package main

import (
	"log"

	"github.com/samuelventura/go-tools"
	"github.com/samuelventura/go-tree"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type daoDso struct {
	db *gorm.DB
}

func dialector(node tree.Node) gorm.Dialector {
	driver := node.GetValue("driver").(string)
	source := node.GetValue("source").(string)
	switch driver {
	case "sqlite":
		return sqlite.Open(source)
	case "postgres":
		return postgres.Open(source)
	}
	log.Panicf("unknown driver %s", driver)
	return nil
}

func newDao(node tree.Node) *daoDso {
	mode := logger.Default.LogMode(logger.Silent)
	//mode := logger.Default.LogMode(logger.Info)
	config := &gorm.Config{Logger: mode}
	dialector := dialector(node)
	db, err := gorm.Open(dialector, config)
	tools.PanicIfError(err)
	err = db.AutoMigrate(&AccountDro{}, &SessionDro{})
	tools.PanicIfError(err)
	return &daoDso{db}
}

func (dso *daoDso) close() {
	sqlDB, err := dso.db.DB()
	tools.PanicIfError(err)
	err = sqlDB.Close()
	tools.PanicIfError(err)
}

func (dso *daoDso) create(dro interface{}) error {
	result := dso.db.Create(dro)
	return result.Error
}

func (dso *daoDso) update(dro interface{}) error {
	//gorm logs confirmed that this updates only
	result := dso.db.Save(dro)
	return result.Error
}

func (dso *daoDso) delete(dro interface{}) error {
	result := dso.db.Delete(dro)
	return result.Error
}

func (dso *daoDso) getAccount(id string) (*AccountDro, error) {
	dro := &AccountDro{}
	result := dso.db.First(dro, "aid = ?", id)
	return dro, result.Error
}

func (dso *daoDso) getSession(id string) (*SessionDro, error) {
	dro := &SessionDro{}
	result := dso.db.First(dro, "sid = ?", id)
	return dro, result.Error
}
