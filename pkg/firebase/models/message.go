package models

import (
	"context"
	"time"

	"github.com/emika-team/line-oa-manager/pkg/firebase"
)

type Message struct {
	Destination        string  `json:"destination"`
	UID                string  `json:"uid"`
	Type               string  `json:"type"`
	Text               string  `json:"text"`
	Emojis             []Emoji `json:"emojis"`
	PackageID          string  `json:"packageId"`
	StickerID          string  `json:"stickerId"`
	OriginalContentUrl string  `json:"originalContentUrl"`
	PreviewImageUrl    string  `json:"previewImageUrl"`
	TrackingId         string  `json:"trackingId"`
	Duration           int     `json:"duration"`
	IsRead             bool    `json:"isRead"`
	CreatedAt          int64   `json:"createdAt"`
	UpdatedAt          int64   `json:"updatedAt"`
}

func (m *Message) Create() error {
	m.IsRead = false
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	_, _, err := firebase.FirestoreClient.Collection("messages").Add(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}
