package models

import "gorm.io/gorm"

type UserAuth struct {
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	Password string `json:"password,omitempty"`
	First    string `json:"first,omitempty"`
	Last     string `json:"last,omitempty"`
	Role     string `gorm:"default:user" json:"role"` // user or admin
}

type User struct {
	gorm.Model
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	First      string `json:"first,omitempty"`
	Last       string `json:"last,omitempty"`
	Role       string `gorm:"default:user" json:"role"`
	Birthday   string `json:"birthday,omitempty"`
	City       string `json:"city,omitempty"`
	Street     string `json:"street,omitempty"`
	PostalCode string `json:"postalCode,omitempty"`
	TaxID      string `json:"taxID,omitempty"`
	Phone      string `json:"phone,omitempty"`
	PhotoURL   string `json:"photo_url,omitempty"`
	XLink      string `json:"xLink,omitempty"`
	GitHubLink string `json:"gitHubLink,omitempty"`
	FBLink     string `json:"fbLink,omitempty"`
}
