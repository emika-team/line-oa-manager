package models

import (
	"context"

	"github.com/emika-team/line-oa-manager/pkg/firebase"
)

type Channel struct {
	ChannelID     string `json:"-" firestore:"-"`
	ChannelName   string `json:"channelName" firestore:"channelName"`
	ChannelSecret string `json:"channelSecret" firestore:"channelSecret"`
	AccessToken   string `json:"accessToken" firestore:"accessToken"`
}

func (c *Channel) CreateOrUpdate() error {
	col := firebase.FirestoreClient.Collection("channel")
	_, err := col.Doc(c.ChannelID).Set(context.Background(), c)
	if err != nil {
		return err
	}
	return nil
}
