package models

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/emika-team/line-oa-manager/pkg/firebase"
)

type Emoji struct {
	Index     int    `json:"index"`
	ProductID string `json:"productId"`
	EmojiID   string `json:"emojiId"`
}

type Message struct {
	ID                 string  `json:"-" firestore:"-"`
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
	Sender             string  `json:"sender" firestore:"sender"`
	CreatedAt          int64   `json:"createdAt" firestore:"createdAt"`
	UpdatedAt          int64   `json:"updatedAt" firestore:"updatedAt"`
}

type Chat struct {
	ChannelUserID     string    `json:"channelUserId" firestore:"channelUserId"`
	IsRead            bool      `json:"isRead" firestore:"isRead"`
	RecentMessageType string    `json:"recentMessageType" firestore:"recentMessageType"`
	RecentMessage     string    `json:"recentMessage" firestore:"recentMessage"`
	RecentAt          time.Time `json:"recentAt" firestore:"recentAt"`
}

func (m *Message) Create(tx *firestore.Transaction) error {
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	isRead := false
	if m.Destination == "" {
		isRead = true
	}
	c := Chat{
		ChannelUserID:     m.Destination,
		IsRead:            isRead,
		RecentMessageType: m.Type,
		RecentMessage:     m.Text,
		RecentAt:          time.Now(),
	}
	fmt.Println(m)
	channelCol := firebase.FirestoreClient.Collection("channels").Where("channelUserId", "==", m.Destination).Limit(1)
	chatDoc, err := channelCol.Documents(context.Background()).Next()
	if err != nil {
		return err
	}

	if chatDoc.Exists() {
		chatCol := chatDoc.Ref.Collection("chats")
		err := tx.Set(chatCol.Doc(m.UID), c)
		if err != nil {
			return err
		}
		err = tx.Set(chatCol.Doc(m.UID).Collection("messages").Doc(m.ID), m)
		if err != nil {
			return err
		}
	}
	return nil
}
