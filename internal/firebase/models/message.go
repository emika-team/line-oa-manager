package models

import (
	"context"

	"github.com/emika-team/line-oa-manager/internal/firebase"
)

type emoji struct {
	Index     int    `json:"index"`
	ProductID string `json:"productId"`
	EmojiID   string `json:"emojiId"`
}

type Message struct {
	Destination        string  `json:"destination"`
	UID                string  `json:"uid"`
	Type               string  `json:"type"`
	Text               string  `json:"text"`
	Emojis             []emoji `json:"emojis"`
	PackageID          string  `json:"packageId"`
	StickerID          string  `json:"stickerId"`
	OriginalContentUrl string  `json:"originalContentUrl"`
	PreviewImageUrl    string  `json:"previewImageUrl"`
	TrackingId         string  `json:"trackingId"`
	Duration           int     `json:"duration"`
}

func (m *Message) Create() error {
	_, _, err := firebase.FirestoreClient.Collection("messages").Add(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}