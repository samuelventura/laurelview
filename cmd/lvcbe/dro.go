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
	Enabled bool
}
