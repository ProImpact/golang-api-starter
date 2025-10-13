package request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"apistarter/internal/server/response"
	"apistarter/pkg/model"

	"github.com/gin-gonic/gin"
)

type Validable interface {
	Validate() error
}

// IsValidRequest unmarchasll into data and checks if the request is valid and send back the proper error messaje
func IsValidRequest[T Validable](ctx *gin.Context, data T) bool {
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {
		case errors.As(err, &syntaxError):
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.INVALID_FORMAT,
				fmt.Sprintf("body contains badly-formed JSON (at character %d)", syntaxError.Offset),
				nil,
			)
			return false
		case errors.Is(err, io.ErrUnexpectedEOF):
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.MALFORMED_REQUEST,
				"body contains badly-formed JSON",
				nil,
			)
			return false
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				response.Error(
					ctx,
					http.StatusBadRequest,
					model.FIELD_VALIDATION_ERROR,
					fmt.Sprintf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field),
					nil,
				)
				return false
			}
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.INVALID_FORMAT,
				fmt.Sprintf("body contains badly-formed JSON (at character %d)", syntaxError.Offset),
				nil,
			)
			return false
		case errors.Is(err, io.EOF):
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.INVALID_REQUEST,
				"body must not be empty",
				nil,
			)
			return false
		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			response.Error(
				ctx,
				http.StatusBadRequest,
				model.INVALID_REQUEST,
				err.Error(),
				nil,
			)
			return false
		}
	}
	err = data.Validate()
	if err != nil {
		response.ValidationError(ctx, err)
		return false
	}
	return true
}
