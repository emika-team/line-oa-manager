package models

import (
	"context"
	"time"

	"github.com/emika-team/line-oa-manager/pkg/firebase"
)

type Message struct {
	Destination        string  `json:"destination" firestore:"destination"`
	UID                string  `json:"-" firestore:"-"`
	Type               string  `json:"type" firestore:"type"`
	Text               string  `json:"text" firestore:"text"`
	Emojis             []Emoji `json:"emojis" firestore:"emojis"`
	PackageID          string  `json:"packageId" firestore:"packageId"`
	StickerID          string  `json:"stickerId" firestore:"stickerId"`
	OriginalContentUrl string  `json:"originalContentUrl" firestore:"originalContentUrl"`
	PreviewImageUrl    string  `json:"previewImageUrl" firestore:"previewImageUrl"`
	TrackingId         string  `json:"trackingId" firestore:"trackingId"`
	Duration           int     `json:"duration" firestore:"duration"`
	IsRead             bool    `json:"isRead" firestore:"isRead"`
	CreatedAt          int64   `json:"createdAt" firestore:"createdAt"`
	UpdatedAt          int64   `json:"updatedAt" firestore:"updatedAt"`
}

func (m *Message) Create() error {
	m.IsRead = false
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	_, err := firebase.FirestoreClient.Collection("messages").Doc(m.UID).Set(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}
