package dto

import "time"

type GetOtpRequest struct {
	MobileNumber string `json:"mobile_number" binding:"required,mobile,min=11,max=11"`
}
type RegisterLoginByMobileRequest struct {
	MobileNumber string `json:"mobileNumber" binding:"required,mobile,min=11,max=11"`
	Otp          string `json:"otp" binding:"required,min=6,max=6"`
}

type UserInfo struct {
	ID           int       `json:"id"`
	MobileNumber string    `json:"mobile_number"`
	RegisteredAt time.Time `json:"registered_at"`
}
type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

type UserList struct {
	Users    []UserInfo `json:"users"`
	Total    int        `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}
