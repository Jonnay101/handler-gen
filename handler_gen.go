package handlergen

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type DomainLogicHandler func(i interface{}) (responseData interface{}, statusCode int, err error)

func EchoHandleFuncGenerator(dFn DomainLogicHandler, i interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		if i == nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				errors.New("Provided binding interface has nil value"),
			)
		}

		if err := c.Bind(i); err != nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				fmt.Errorf("Error while binding request >>>> \n%w\n", err),
			)
		}

		if dFn == nil {
			return echo.NewHTTPError(
				http.StatusBadRequest,
				errors.New("Provided DomainLogicHandler has nil value"),
			)
		}

		return nil
	}
}
