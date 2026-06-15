package model

type TokenInfo struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Token     string `json:"token"`
	Enabled   bool   `json:"enabled"`
	CreatedAt int    `json:"createdAt"`
}

type XUIResponse[T any] struct {
	Success bool   `json:"success"`
	Msg     string `json:"msg"`
	Obj     T      `json:"obj"`
}

type XUIInbound struct {
	ID       int                `json:"id"`
	Remark   string             `json:"remark"`
	Enable   bool               `json:"enable"`
	Listen   string             `json:"listen"`
	Port     int                `json:"port"`
	Settings XUIInboundSettings `json:"settings"`
}

type XUIInboundClient struct {
	ID    any    `json:"id"`
	Email string `json:"email"`
	Flow  string `json:"flow"`
	TgID  int    `json:"tgId"`
	SubID string `json:"subId"`
}

type XUIInboundSettings struct {
	Cleints []XUIInboundClient `json:"clients"`
}
