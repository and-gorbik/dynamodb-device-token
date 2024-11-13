package model

type TokenKind string

const (
	TokenKindAndroidGeneral  TokenKind = "android_general"
	TokenKindIOSGeneral      TokenKind = "ios_general"
	TokenKindIOSLiveActivity TokenKind = "ios_live_activity"
	TokenKindIOSVoip         TokenKind = "ios_voip"
)

type DeviceToken struct {
	UserID      int64     `json:"user_id"`
	ModifiedAt  int64     `json:"modified_at"`
	Kind        TokenKind `json:"kind"`
	Token       string    `json:"token"`
	AppVersion  string    `json:"app_version"`
	DeviceModel string    `json:"device_model"`
}
