package message

import (
	"fmt"
	"net/http"

	"github.com/emika-team/line-oa-manager/pkg/firebase/models"
	"github.com/labstack/echo/v4"
)

func GetContent(c echo.Context) error {
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
				fmt.Println(err)
				return c.JSON(http.StatusInternalServerError, "internal server error")
			}
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func GetImageContent(c echo.Context) error {
	return nil
}
