package entity

import "time"

type Auth struct {
	UID         string    `json:"uid"`
	UserID      string    `json:"userId"`
	AccessToken string    `json:"accessToken"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	IsLogout    bool      `json:"isLogout"`
}
