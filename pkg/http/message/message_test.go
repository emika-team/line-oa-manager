package message_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/emika-team/line-oa-manager/pkg/http/message"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

var (
	webhookRawRequestJson = `{
		"destination": "U96264f5facec3f4901e9cfc1339fe030",
		"events": [
		  {
			"deliveryContext": {
			  "isRedelivery": false
			},
			"message": {
			  "id": "17445705109139",
			  "text": "รับเครดิตฟรี",
			  "type": "text"
			},
			"mode": "active",
			"replyToken": "6c9689c27c4e4d249946a99e41ae1d4f",
			"source": {
			  "type": "user",
			  "userId": "testing"
			},
			"timestamp": 1673331459900,
			"type": "message",
			"webhookEventId": "01GPD57J11B7PGSBM7024Y3T7S"
		  }
		]
	}`
)

func TestReceiveMessage(t *testing.T) {
	fmt.Println("test on function")
	fmt.Println(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))
	e := echo.New()

	req := httptest.NewRequest(http.MethodPost, "/webhook", strings.NewReader(webhookRawRequestJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, message.ReceiveMessage(c)) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
	}
}
