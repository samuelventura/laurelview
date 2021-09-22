package lvcdb

import (
	"gorm.io/gorm"
)

type AccountDro struct {
	gorm.Model
	Email  string `gorm:"index:idx_email,unique"`
	Name   string
	Avatar string
}

type NodeDro struct {
	gorm.Model
	Mac   string `gorm:"index:idx_mac,unique"`
	Owner string //Owner Email
	Name  string
	Pip   string //Proxy IP
	Sshp  int    //SSH port
	Httpp int    //HTTP port
}

type SessionDro struct {
	gorm.Model
	Token  string `gorm:"index:idx_token,unique"`
	Email  string
	Name   string
	Avatar string
	Origin string //IP and Browser profile
}

type RequestDro struct {
	gorm.Model
	Token   string
	Address string //URL + query
}
