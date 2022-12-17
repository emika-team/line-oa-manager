package models

import (
	"context"
	"time"

	"github.com/emika-team/line-oa-manager/pkg/firebase"
)

type Emoji struct {
	Index     int    `json:"index"`
	ProductID string `json:"productId"`
	EmojiID   string `json:"emojiId"`
}

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
	Sender             string  `json:"sender" firestore:"sender"`
	CreatedAt          int64   `json:"createdAt" firestore:"createdAt"`
	UpdatedAt          int64   `json:"updatedAt" firestore:"updatedAt"`
}

type Chat struct {
	RecentMessage string `json:"recentMessage" firestore:"recentMessage"`
}

func (m *Message) Create() error {
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	c := Chat{
		RecentMessage: m.Text,
	}
	chatCol := firebase.FirestoreClient.Collection("chat")
	_, _, err := chatCol.Doc(m.UID).Parent.Add(context.Background(), c)
	if err != nil {
		return err
	}
	_, _, err = chatCol.Parent.Collection("messages").Add(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}
