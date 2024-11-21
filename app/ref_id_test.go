package app

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRefIDMiddleware(t *testing.T) {
	t.Run("middleware keep ref ID in context", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)

		givenRefID := "111222333444"

		req.Header.Add("X-REF-ID", givenRefID)
		c.Request = req

		e.Use(RefIDMiddleware("X-REF-ID"))

		e.POST("/api", handler)

		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() != givenRefID {
			t.Errorf("expect ref-id as result %q but actual %q\n", givenRefID, w.Body.String())
		}
	})
	t.Run("known ref-id key but empty value", func(t *testing.T) {
		handler := func(c *gin.Context) {
			if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
				c.String(200, v)
				return
			}
			c.String(500, "not found")
		}

		w := httptest.NewRecorder()
		c, e := gin.CreateTestContext(w)
		req := httptest.NewRequest(http.MethodPost, "/api", nil)

		req.Header.Add("X-REF-ID", "")
		c.Request = req

		e.Use(RefIDMiddleware("X-REF-ID"))

		e.POST("/api", handler)

		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() == "" {
			t.Errorf("expect ref-id not empty\n")
		}
	})
}

func TestForwardRefOption(t *testing.T) {
	givenRedID := "12345"

	ctx := context.WithValue(context.Background(), refIDKey, givenRedID)
	req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/", nil)
	ForwardRefIDOption(req, ctx)

	if givenRedID != req.Header.Get(string(refIDKey)) {
		t.Errorf("%s is expected but got %s\n", givenRedID, req.Header.Get(string(refIDKey)))
	}
}

func TestSetGetRefID(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api", nil)
	ctx.Request = req
	SetRefID(ctx, "1234")

	refID := RefID(ctx)

	if refID != "1234" {
		t.Errorf("1234 was expected but get %q", refID)
	}
}

func TestRefID(t *testing.T) {
	t.Run("generally get ref-id", func(t *testing.T) {
		givenRefID := "111222333444"

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		ctx := newRefIDContext(context.Background(), givenRefID)

		req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api", nil)

		c.Request = req.WithContext(ctx)

		result := RefID(c)
		if givenRefID != result {
			t.Errorf("expect ref-id %q but actual %q\n", givenRefID, result)
		}
	})
	t.Run("can not get ref-id", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		req := httptest.NewRequest(http.MethodPost, "http://0.0.0.0/api", nil)

		c.Request = req

		result := RefID(c)
		if result != "" {
			t.Errorf("expect empty string but actual %q\n", result)
		}
	})
}
