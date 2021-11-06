package lvndb

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Dao interface {
	Close()
	All() []ItemDro
	Delete(id uint)
	Create(name string, json string) ItemDro
	Update(id uint, name string, json string) ItemDro
}

type ItemDro struct {
	gorm.Model
	Name string
	Json string
}

type SyncDro struct {
	ID int64
}

type daoDso struct {
	db *gorm.DB
}

func NewDao(path string) Dao {
	dialect := sqlite.Open(path)
	config := &gorm.Config{}
	db, err := gorm.Open(dialect, config)
	PanicIfError(err)
	err = db.AutoMigrate(&ItemDro{}, &SyncDro{})
	PanicIfError(err)
	return &daoDso{db}
}

//suden reboot of linux sbc
//looses last created item
//if two are created only
//last one gets lost
//this is a sync shim
func (dao *daoDso) sync() {
	row := &SyncDro{}
	dao.db.Create(row)
	dao.db.Unscoped().Delete(&SyncDro{}, row.ID)
}

func (dao *daoDso) Close() {
	sqlDB, err := dao.db.DB()
	PanicIfError(err)
	sqlDB.Close()
}

func (dao *daoDso) All() (items []ItemDro) {
	//predictable order for testing purposes
	result := dao.db.Order("id").Find(&items)
	PanicIfError(result.Error)
	return
}

func (dao *daoDso) Create(name string, json string) ItemDro {
	row := &ItemDro{Name: name, Json: json}
	result := dao.db.Create(row)
	PanicIfError(result.Error)
	dao.sync()
	return *row
}

func (dao *daoDso) Delete(id uint) {
	result := dao.db.Delete(&ItemDro{}, id)
	PanicIfError(result.Error)
	dao.sync()
}

func (dao *daoDso) Update(id uint, name string, json string) ItemDro {
	row := &ItemDro{}
	result := dao.db.First(row, id)
	PanicIfError(result.Error)
	row.Name = name
	row.Json = json
	result = dao.db.Save(row)
	PanicIfError(result.Error)
	dao.sync()
	return *row
}
