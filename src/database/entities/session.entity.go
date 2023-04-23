package entities

import "gopkg.in/guregu/null.v4"

type SessionEntity struct {
	BaseEntity
	UserId    string      `json:"UserId" db:"user_id"`
	Token     string      `json:"token" db:"token"`
	IpAddress null.String `json:"ipAddress" db:"ip_address"`
	Device    null.String `json:"device" db:"device"`
	Platform  null.String `json:"platform" db:"platform"`
	Latitude  null.String `json:"latitude" db:"latitude"`
	Longitude null.String `json:"longitude" db:"longitude"`
}
