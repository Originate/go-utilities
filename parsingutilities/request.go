package parsingutilities

import (
	"encoding/json"
	"io"

	"github.com/Originate/go-utilities/errorutilities"
	"github.com/gin-gonic/gin"
)

func UnmarshalBody[T any](ctx *gin.Context) (T, error) {
	var result T

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return result, errorutilities.NewErrorf("error reading body from request: %s", err.Error())
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return result, errorutilities.NewErrorf("error unmarshaling json body: %s", err.Error())
	}

	return result, nil
}
