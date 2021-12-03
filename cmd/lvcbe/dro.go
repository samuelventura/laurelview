package main

import "time"

type AccountDro struct {
	Aid      string `gorm:"primaryKey"`
	Password string
	Recover  string
	Enabled  bool
	Created  time.Time
}

type SessionDro struct {
	Sid     string `gorm:"primaryKey"`
	Aid     string `gorm:"index"`
	Created time.Time
	Expires time.Time
}

type NodeDro struct {
	Nid     string `gorm:"primaryKey"`
	Aid     string `gorm:"index"`
	Name    string
	Created time.Time
}

type CreditDro struct {
	Cid     string `gorm:"primaryKey"`
	Aid     string `gorm:"index"`
	Nid     string `gorm:"index"`
	Months  int
	Created time.Time
	Expires time.Time
}
