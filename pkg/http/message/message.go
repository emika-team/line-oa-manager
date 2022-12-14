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

		if t == "message" {
			message := models.Message{
				Destination:        content["destination"].(string),
				UID:                userId,
				Type:               mtype,
				Text:               mtext,
				Emojis:             []models.Emoji{},
				PackageID:          "",
				StickerID:          "",
				OriginalContentUrl: "",
				PreviewImageUrl:    "",
				TrackingId:         "",
				Duration:           0,
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
