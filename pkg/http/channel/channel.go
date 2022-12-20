package channel

import (
	"github.com/emika-team/line-oa-manager/pkg/firebase"
	"github.com/labstack/echo/v4"
)

func GetChannel(c echo.Context) error {
	col := firebase.FirestoreClient.Collection("channel")
	docs := col.Documents(c.Request().Context())
	result := []map[string]interface{}{}
	for {
		doc, err := docs.Next()
		if err != nil {
			break
		}
		c := doc.Data()
		c["id"] = doc.Ref.ID
		result = append(result, c)
	}
	return c.JSON(200, result)
}
