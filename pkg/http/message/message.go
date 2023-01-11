package message

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/emika-team/line-oa-manager/pkg/firebase"
	"github.com/emika-team/line-oa-manager/pkg/firebase/models"
	httpclient "github.com/emika-team/line-oa-manager/pkg/http/httpclient"
	response "github.com/emika-team/line-oa-manager/pkg/http/response"
	"github.com/labstack/echo/v4"
)

var LineAPIDataEndPoint = "https://api-data.line.me"
var LineAPIEndPoint = "https://api.line.me"

func buildMessage(event interface{}, destination string) models.Message {
	m := event.(map[string]interface{})["message"].(map[string]interface{})
	userId, ok := event.(map[string]interface{})["source"].(map[string]interface{})["userId"].(string)
	if !ok {
		userId = ""
	}
	mtype, ok := m["type"].(string)
	if !ok {
		mtype = ""
	}
	mtext, ok := m["text"].(string)
	if !ok {
		mtext = ""
	}
	emojis := []models.Emoji{}
	e, ok := m["emojis"].([]interface{})
	if !ok {
		e = []interface{}{}
	}
	for _, e := range e {
		emojis = append(emojis, models.Emoji{
			Index:     int(e.(map[string]interface{})["index"].(float64)),
			ProductID: e.(map[string]interface{})["productId"].(string),
			EmojiID:   e.(map[string]interface{})["emojiId"].(string),
		})
	}
	packageId, ok := m["packageId"].(string)
	if !ok {
		packageId = ""
	}
	stickerId, ok := m["stickerId"].(string)
	if !ok {
		stickerId = ""
	}
	originalContentUrl, ok := m["originalContentUrl"].(string)
	if !ok {
		originalContentUrl = ""
	}
	previewImageUrl, ok := m["previewImageUrl"].(string)
	if !ok {
		previewImageUrl = ""
	}
	trackingId, ok := m["trackingId"].(string)
	if !ok {
		trackingId = ""
	}
	duration, ok := m["duration"].(float64)
	if !ok {
		duration = 0
	}
	sender, ok := event.(map[string]interface{})["sender"].(string)
	if !ok {
		sender = userId
	}
	return models.Message{
		ID:                 m["id"].(string),
		Destination:        destination,
		UID:                userId,
		Type:               mtype,
		Text:               mtext,
		Emojis:             emojis,
		PackageID:          packageId,
		StickerID:          stickerId,
		OriginalContentUrl: originalContentUrl,
		PreviewImageUrl:    previewImageUrl,
		TrackingId:         trackingId,
		Duration:           int(duration),
		Sender:             sender,
	}
}

func ReceiveMessage(c echo.Context) error {
	content := map[string]interface{}{}
	err := c.Bind(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println(content)
	err = firebase.FirestoreClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		for _, v := range content["events"].([]interface{}) {
			t := v.(map[string]interface{})["type"]

			if t == "message" {
				message := buildMessage(v, content["destination"].(string))
				err := message.CreateWithTransaction(tx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func GetContent(c echo.Context) error {
	messageID := c.Param("messageId")
	channelID := c.Param("channelId")
	documentRef := firebase.FirestoreClient.Collection("channel").Doc(channelID)
	channel := models.Channel{}
	documentSnapshot, err := documentRef.Get(context.Background())
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	err = documentSnapshot.DataTo(&channel)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	h := map[string]interface{}{
		"Authorization": fmt.Sprintf("Bearer %s", channel.AccessToken),
	}
	url := fmt.Sprintf("%s/v2/bot/message/%s/content", LineAPIDataEndPoint, messageID)
	result, err := httpclient.HttpRequest("GET", url, nil, nil, h)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": base64.StdEncoding.EncodeToString(result),
	})
}

func SendMessage(c echo.Context) error {
	channelID := c.Param("channelId")
	documentRef := firebase.FirestoreClient.Collection("channel").Doc(channelID)
	channel := models.Channel{}
	documentSnapshot, err := documentRef.Get(context.Background())
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	err = documentSnapshot.DataTo(&channel)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	content := map[string]interface{}{}
	err = c.Bind(&content)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	event := interface{}(map[string]interface{}{
		"sender": channel.ChannelID,
		"source": map[string]interface{}{
			"userId": content["to"],
		},
		"type":    "message",
		"message": content["message"],
	})
	m := buildMessage(event, "")
	err = firebase.FirestoreClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		err := m.CreateWithTransaction(tx)
		if err != nil {
			return err
		}
		h := map[string]interface{}{
			"Authorization": fmt.Sprintf("Bearer %s", channel.AccessToken),
			"Content-Type":  "application/json",
		}
		url := fmt.Sprintf("%s/v2/bot/message/push", LineAPIEndPoint)
		_, err = httpclient.HttpRequest("POST", url, content, nil, h)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}
