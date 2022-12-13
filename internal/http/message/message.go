package message

import (
	"fmt"
	"net/http"

	"github.com/emika-team/line-oa-manager/internal/firebase/models"
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
			message := models.Message{
				Destination:        content["destination"].(string),
				UID:                v.(map[string]interface{})["source"].(map[string]interface{})["userId"].(string),
				Type:               v.(map[string]interface{})["message"].(map[string]interface{})["type"].(string),
				Text:               v.(map[string]interface{})["message"].(map[string]interface{})["text"].(string),
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
				return c.JSON(http.StatusInternalServerError, "internal server error")
			}
		}
	}
	return c.NoContent(http.StatusNoContent)
}
