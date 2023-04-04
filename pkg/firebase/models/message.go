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
	ChannelID         string    `json:"channelId" firestore:"channelId"`
	IsRead            bool      `json:"isRead" firestore:"isRead"`
	RecentMessageType string    `json:"recentMessageType" firestore:"recentMessageType"`
	RecentMessage     string    `json:"recentMessage" firestore:"recentMessage"`
	RecentAt          time.Time `json:"recentAt" firestore:"recentAt"`
}

func (m *Message) Create() error {
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	c := Chat{
		ChannelID:         m.Destination,
		IsRead:            false,
		RecentMessageType: m.Type,
		RecentMessage:     m.Text,
		RecentAt:          time.Now(),
	}
	chatCol := firebase.FirestoreClient.Collection("chat")
	_, err := chatCol.Doc(m.UID).Set(context.Background(), c)
	if err != nil {
		return err
	}
	_, err = chatCol.Doc(m.UID).Collection("messages").Doc(m.ID).Set(context.Background(), m)
	if err != nil {
		return err
	}
	return nil
}

func (m *Message) CreateWithTransaction(tx *firestore.Transaction) error {
	m.CreatedAt = time.Now().Unix()
	m.UpdatedAt = time.Now().Unix()
	isRead := false
	if m.Destination == "" {
		isRead = true
	}
	c := Chat{
		IsRead:            isRead,
		RecentMessageType: m.Type,
		RecentMessage:     m.Text,
		RecentAt:          time.Now(),
	}
	chatCol := firebase.FirestoreClient.Collection("chat")
	err := tx.Set(chatCol.Doc(m.UID), c)
	if err != nil {
		return err
	}
	fmt.Println(m)
	err = tx.Set(chatCol.Doc(m.UID).Collection("messages").Doc(m.ID), m)
	if err != nil {
		return err
	}
	return nil
}
