package lvcdb

import (
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type daoDso struct {
	db *gorm.DB
}

type Dao interface {
	Close()
	ListNodes(owner string) []NodeDro
	DeleteNode(owner string, mac string)
	SaveNode(owner string, mac string, name string)
	GetNode(owner string, mac string) NodeDro
}

func Dialector(driver string, source string) gorm.Dialector {
	switch driver {
	case "sqlite":
		return sqlite.Open(source)
	case "postgres":
		return postgres.Open(source)
	}
	PanicF("unknown driver %s", driver)
	return nil
}

func NewDao(driver string, source string) Dao {
	config := &gorm.Config{}
	dialector := Dialector(driver, source)
	db, err := gorm.Open(dialector, config)
	PanicIfError(err)
	err = db.AutoMigrate(&NodeDro{})
	PanicIfError(err)
	return &daoDso{db}
}

func (dao *daoDso) Close() {
	sqlDB, err := dao.db.DB()
	PanicIfError(err)
	sqlDB.Close()
}

func (dao *daoDso) ListNodes(owner string) (nodes []NodeDro) {
	//predictable order for testing purposes
	result := dao.db.Where("owner = ?", owner).Order("id").Find(&nodes)
	PanicIfError(result.Error)
	return
}

func (dao *daoDso) DeleteNode(owner string, mac string) {
	result := dao.db.Where("owner = ?", owner).Delete(&NodeDro{})
	PanicIfError(result.Error)
}

func (dao *daoDso) SaveNode(owner string, mac string, name string) {
	row := &NodeDro{Owner: owner, Name: name, Mac: mac}
	result := dao.db.Create(row)
	PanicIfError(result.Error)
}

func (dao *daoDso) GetNode(owner string, mac string) NodeDro {
	row := NodeDro{}
	result := dao.db.Where("owner = ? AND mac = ?", owner, mac).First(&row)
	PanicIfError(result.Error)
	return row
}
