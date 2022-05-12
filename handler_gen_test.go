package handlergen

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/music-tribe/uuid"
)

var (
	srv      *echo.Echo
	testUUID = uuid.New()
)

type testObj struct {
	ID   uuid.UUID `json:"id" param:"id"`
	Name string    `json:"name,omitempty"`
	Age  uint      `json:"age,omitempty"`
}

func TestMain(m *testing.M) {
	srv = echo.New()
	exitVal := m.Run()

	os.Exit(exitVal)
}

func TestEchoHandleFuncGenerator(t *testing.T) {
	type req struct {
		method  string
		target  string
		body    io.Reader
		paramID string
	}
	type args struct {
		dFn DomainLogicHandler
		i   interface{}
	}
	tests := []struct {
		name           string
		req            req
		args           args
		wantErr        bool
		wantStatusCode int
	}{
		{
			"Error: cannot bind to nil value interface",
			req{
				method:  http.MethodPost,
				target:  "/",
				body:    nil,
				paramID: testUUID.String(),
			},
			args{
				func(i interface{}) (responseData interface{}, statusCode int, err error) {
					return
				},
				nil, // nil value interface
			},
			true,
			400,
		},
		{
			"Error: cannot bind - uuid param is invalid",
			req{
				method: http.MethodPost,
				target: "/",
				body:   nil,
				// id param is missing from request
			},
			args{
				func(i interface{}) (responseData interface{}, statusCode int, err error) {
					return
				},
				&testObj{}, // nil value interface
			},
			true,
			400,
		},
		{
			"Error: domain logic func has nil value - potential panic",
			req{
				method:  http.MethodPost,
				target:  "/",
				body:    nil,
				paramID: testUUID.String(),
			},
			args{
				nil,
				&testObj{}, // nil value interface
			},
			true,
			400,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := EchoHandleFuncGenerator(tt.args.dFn, tt.args.i)

			req := httptest.NewRequest(tt.req.method, tt.req.target, tt.req.body)
			rec := httptest.NewRecorder()
			ctx := srv.NewContext(req, rec)
			ctx.SetParamNames("id")
			ctx.SetParamValues(tt.req.paramID)

			err := handler(ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("wanted error to be %v but got %v", tt.wantErr, err)
			}

			t.Logf("%s error: %v\n", tt.name, err)

			// assert the error to use echo handler
			echErr := new(echo.HTTPError)
			if ok := errors.As(err, &echErr); ok {
				if echErr.Code != tt.wantStatusCode {
					t.Errorf("wanted status code to be %d but got %d", tt.wantStatusCode, echErr.Code)
				}
			}
		})
	}
}
