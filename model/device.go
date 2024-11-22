package model

import (
	"fmt"
	"strconv"
	"strings"
)

type TokenKind string

const (
	TokenKindAndroidGeneral  TokenKind = "android_general"
	TokenKindIOSGeneral      TokenKind = "ios_general"
	TokenKindIOSLiveActivity TokenKind = "ios_live_activity"
	TokenKindIOSVoip         TokenKind = "ios_voip"
)

type Device struct {
	UserID      int64     `json:"user_id"`
	Kind        TokenKind `json:"kind"`
	DeviceModel string    `json:"device_model"`
	ModifiedAt  int64     `json:"modified_at"`
	Token       string    `json:"token"`
	AppVersion  string    `json:"app_version"`
	Locale      string    `json:"locale"`
}

func (d *Device) PartitionKey() string {
	return strconv.FormatInt(d.UserID, 10)
}

func (d *Device) SortKey() string {
	return fmt.Sprintf("%s#%s", d.Kind, d.DeviceModel)
}

func (d *Device) SetSortKey(sortKey string) {
	parts := strings.Split(sortKey, "#")
	d.Kind = TokenKind(parts[0])
	d.DeviceModel = parts[1]
}
