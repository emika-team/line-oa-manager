package message

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/emika-team/line-oa-manager/pkg/dto"
	"github.com/emika-team/line-oa-manager/pkg/firebase"
	"github.com/emika-team/line-oa-manager/pkg/firebase/models"
	httpclient "github.com/emika-team/line-oa-manager/pkg/http/httpclient"
	response "github.com/emika-team/line-oa-manager/pkg/http/response"
	"github.com/emika-team/line-oa-manager/pkg/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func ReceiveMessage(c echo.Context) error {
	content := map[string]interface{}{}
	err := c.Bind(&content)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}
	err = firebase.FirestoreClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		for _, v := range content["events"].([]interface{}) {
			t := v.(map[string]interface{})["type"]

			if t == "message" {
				message := utils.BuildMessage(v, content["destination"].(string))
				err := message.Create(tx)
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
	documentRef := firebase.FirestoreClient.Collection("channels").Doc(channelID)
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
	url := fmt.Sprintf("%s/v2/bot/message/%s/content", dto.LineAPIDataEndPoint, messageID)
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
	documentRef := firebase.FirestoreClient.Collection("channels").Doc(channelID)
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
	fmt.Println(content)
	messages := []map[string]interface{}{}
	jsonMessages, err := json.Marshal(content["messages"].([]interface{}))
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	err = json.Unmarshal([]byte(jsonMessages), &messages)
	if err != nil {
		return response.ReturnInternalServerError(c, err)
	}
	message := messages[0]
	message["id"] = uuid.New().String()
	event := interface{}(map[string]interface{}{
		"sender": channel.ChannelID,
		"source": map[string]interface{}{
			"userId": content["to"],
		},
		"type":    "message",
		"message": message,
	})
	m := utils.BuildMessage(event, "")
	err = firebase.FirestoreClient.RunTransaction(context.Background(), func(ctx context.Context, tx *firestore.Transaction) error {
		err := m.Create(tx)
		if err != nil {
			return err
		}
		h := map[string]interface{}{
			"Authorization":    fmt.Sprintf("Bearer %s", channel.AccessToken),
			"Content-Type":     "application/json",
			"X-Line-Retry-Key": message["id"].(string),
		}
		url := fmt.Sprintf("%s/v2/bot/message/push", dto.LineAPIEndPoint)
		data := map[string]interface{}{
			"to":       content["to"],
			"messages": []interface{}{message},
		}
		_, err = httpclient.HttpRequest("POST", url, data, nil, h)
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
