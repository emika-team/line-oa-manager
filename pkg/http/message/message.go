package message

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/emika-team/line-oa-manager/pkg/firebase"
	"github.com/emika-team/line-oa-manager/pkg/firebase/models"
	httpclient "github.com/emika-team/line-oa-manager/pkg/http/httpclient"
	response "github.com/emika-team/line-oa-manager/pkg/http/response"
	"github.com/labstack/echo/v4"
)

var LineEndPoint = "https://api-data.line.me"

func ReceiveMessage(c echo.Context) error {
	content := map[string]interface{}{}
	err := c.Bind(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	fmt.Println(content)
	for _, v := range content["events"].([]interface{}) {
		t := v.(map[string]interface{})["type"]

		if t == "message" {
			userId, ok := v.(map[string]interface{})["source"].(map[string]interface{})["userId"].(string)
			if !ok {
				userId = ""
			}
			m := v.(map[string]interface{})["message"].(map[string]interface{})
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
			message := models.Message{
				ID:                 m["id"].(string),
				Destination:        content["destination"].(string),
				UID:                userId,
				Type:               mtype,
				Text:               mtext,
				Emojis:             emojis,
				PackageID:          packageId,
				StickerID:          stickerId,
				OriginalContentUrl: "",
				PreviewImageUrl:    "",
				TrackingId:         "",
				Duration:           0,
				Sender:             userId,
			}
			err := message.Create()
			if err != nil {
				return response.ReturnInternalServerError(c, err)
			}
		}
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
	url := fmt.Sprintf("%s/v2/bot/message/%s/content", LineEndPoint, messageID)
	result, err := httpclient.HttpRequest("GET", url, nil, nil, h)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": base64.StdEncoding.EncodeToString(result),
	})
}
