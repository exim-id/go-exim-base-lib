package models

import (
	"time"

	"github.com/exim-id/go-exim-base-lib/utils"
	"github.com/google/uuid"
)

type Session struct {
	Jti         string     `bson:"jti,omitempty"`
	Mode        string     `bson:"mode,omitempty"`
	Id          *uuid.UUID `bson:"id"`
	CreatedDate time.Time  `bson:"createdDate"`
	UpdatedDate *time.Time `bson:"updatedDate"`
}

type Client struct {
	Id        uuid.UUID     `bson:"id"`
	FrontName string        `bson:"frontName,omitempty"`
	MidName   string        `bson:"midName,omitempty"`
	LastName  string        `bson:"lastName,omitempty"`
	NickName  string        `bson:"nickName,omitempty"`
	Company   ClientCompany `bson:"company,inline"`
	Active    bool          `bson:"active"`
}

type ClientCompany struct {
	Name      string     `bson:"name,omitempty"`
	Address   GeoAddress `bson:"address,inline"`
	Logo      Photo      `bson:"logo,inline"`
	Galleries []Photo    `bson:"galleries,inline"`
}

type Employer struct {
	Id          uuid.UUID          `bson:"id"`
	FrontName   string             `bson:"frontName,omitempty"`
	MidName     string             `bson:"midName,omitempty"`
	LastName    string             `bson:"lastName,omitempty"`
	NickName    string             `bson:"nickName,omitempty"`
	Photos      []Photo            `bson:"photos,inline"`
	Contacts    []Contact          `bson:"contacts,inline"`
	Auth        EmployerAuth       `bson:"auth,inline"`
	Roles       []EmployerRole     `bson:"roles,inline"`
	Active      bool               `bson:"active"`
	CreateDate  time.Time          `bson:"createdDate"`
	UpdatedDate *time.Time         `bson:"updatedDate"`
	Contracts   []EmployerContract `bson:"contracts,inline"`
}

type EmployerContract struct {
	Start       time.Time           `bson:"start"`
	PromiseEnd  time.Time           `bson:"promiseEnd"`
	End         *time.Time          `bson:"end"`
	Salary      utils.Money         `bson:"salary,inline"`
	Status      string              `bson:"status,omitempty"`
	Supports    []EmployerSupport   `bson:"supports,inline"`
	Active      bool                `bson:"active"`
	Liabilities []EmployerLiability `bson:"liabilities,inline"`
}

type EmployerLiability struct {
	Amount utils.Money `bson:"amount,inline"`
	Name   string      `bson:"name,omitempty"`
	Active bool        `bson:"active"`
}

type EmployerSupport struct {
	Amount utils.Money `bson:"amount,inline"`
	Name   string      `bson:"name,omitempty"`
	Active bool        `bson:"active"`
}

type Photo struct {
	Path   string    `bson:"path,omitempty"`
	Active bool      `bson:"active"`
	Start  time.Time `bson:"start"`
}

type Contact struct {
	Phone       string     `bson:"phone,omitempty"`
	Fax         string     `bson:"fax,omitempty"`
	Wa          string     `bson:"wa,omitempty"`
	Telegram    string     `bson:"telegram,omitempty"`
	Address     GeoAddress `bson:"address,inline"`
	TiktokId    string     `bson:"tiktokId,omitempty"`
	InstagramId string     `bson:"instagramId,omitempty"`
	LinkedIn    string     `bson:"linkedIn,omitempty"`
}

type GeoAddress struct {
	Detail       string `bson:"detail,omitempty"`
	StreetNumber string `bson:"streetNumber,omitempty"`
	Block        string `bson:"block,omitempty"`
	Village      string `bson:"village,omitempty"`
	SubDistrict  string `bson:"subDistrict,omitempty"`
	District     string `bson:"district,omitempty"`
	City         string `bson:"city,omitempty"`
	Province     string `bson:"province,omitempty"`
	State        string `bson:"state,omitempty"`
	ZipCode      string `bson:"zipCode,omitempty"`
	Country      string `bosn:"country,omitempty"`
}

type EmployerAuth struct {
	Username           string    `bson:"username,omitempty"`
	Password           string    `bson:"password,omitempty"`
	LastChangePassword time.Time `bson:"lastChangePassword"`
}

type EmployerRole struct {
	Name   string    `bson:"name,omitempty"`
	Active bool      `bson:"active"`
	Start  time.Time `bson:"start"`
}
