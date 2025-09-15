package models

import (
	"database/sql"
	"time"
)

type User struct {
	Id int `gorm:"primarykey"`

	MobileNumber string    `gorm:"type:string;size:11;null;unique;default:null"`
	Enabled      bool      `gorm:"default:true"`
	RegisteredAt time.Time `gorm:"type:TIMESTAMP with time zone;not null"`
	Password     string    `gorm:"type:string;size:255;not null"`
	UserRoles    *[]UserRole

	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

type Role struct {
	Id         int    `gorm:"primarykey"`
	Name       string `gorm:"type:string;size:10;not null,unique"`
	UserRoles  *[]UserRole
	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

type UserRole struct {
	Id         int  `gorm:"primarykey"`
	User       User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	Role       Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	UserId     int
	RoleId     int
	CreatedAt  time.Time    `gorm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  sql.NullTime `gorm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  int            `gorm:"not null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}
