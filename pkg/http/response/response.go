package response

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

var InternalServerError = map[string]interface{}{
	"message": "internal server error",
}

func ReturnInternalServerError(c echo.Context, err error) error {
	fmt.Println(err)
	return c.JSON(http.StatusInternalServerError, InternalServerError)
}
