package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNewTraeparent(t *testing.T) {
	traceparent := NewTraceParent().String()

	tp, err := Parse(traceparent)
	if err != nil {
		t.Errorf("valid traceparent %q, error %s\n", traceparent, err)
	}

	if tp.SpanID.String() == "" {
		t.Error("a new traceparent expect not empty spad-id\n")
	}

	if tp.TraceID.String() == "" {
		t.Error("a new traceparent expect not empty trace-id\n")
	}
}

func TestParseTraceparent(t *testing.T) {
	t.Run("valid traceparent", func(t *testing.T) {
		traceparent := "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"
		tp, err := Parse(traceparent)
		if err != nil {
			t.Errorf("valid traceparent %q, error %s\n", traceparent, err)
		}

		if tp.SpanID.String() != "b7ad6b7169203331" {
			t.Errorf("expect b7ad6b7169203331 but got %q\n", tp.SpanID.String())
		}
		if tp.TraceID.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect 0af7651916cd43dd8448eb211c80319c but got %q\n", tp.TraceID.String())
		}
	})
	t.Run("invalid pattern traceparent", func(t *testing.T) {
		traceparent := "5f2b8701-6c14-436f-8e69-92f1670f0ec2"
		tp, err := Parse(traceparent)
		if err == nil {
			t.Errorf("invalid pattern traceparent expect error but we've got trace-id %q and span-id %q\n\n", tp.TraceID, tp.SpanID)
		}
	})
	t.Run("empty traceparent", func(t *testing.T) {
		traceparent := ""
		tp, err := Parse(traceparent)
		if err == nil {
			t.Errorf("empty pattern traceparent expect error but we've got trace-id %q and span-id %q\n\n", tp.TraceID, tp.SpanID)
		}
	})
}

func TestDefaultIDGenerator(t *testing.T) {
	gen := defIDGenerator()
	spanID1 := gen.NewSpanID()
	traceID1 := gen.NewTraceID()
	if len(spanID1) != 8 {
		t.Errorf("span-id len expect 8 but actual %d\n", len(spanID1))
	}
	if len(traceID1) != 16 {
		t.Errorf("trace-id len expect 8 but actual %d\n", len(spanID1))
	}

	spanID2 := gen.NewSpanID()
	traceID2 := gen.NewTraceID()

	if spanID1 == spanID2 {
		t.Errorf("span should randon everytime")
	}
	if traceID1 == traceID2 {
		t.Errorf("trace should randon everytime")
	}
}

func TestTraceContextTraceIDMIddleware(t *testing.T) {
	t.Run("ref-id key as traceparent", func(t *testing.T) {
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

		req.Header.Add(traceparentHeaderKey, "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01")
		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(traceparentHeaderKey))

		e.POST("/api", handler)

		e.HandleContext(c)

		if w.Result().StatusCode != 200 {
			t.Errorf("expect http status 200 but actual %d\n", w.Result().StatusCode)
		}
		if w.Body.String() != "0af7651916cd43dd8448eb211c80319c" {
			t.Errorf("expect ref-id not empty\n")
		}
	})
	t.Run("ref-id key as traceparent but empty", func(t *testing.T) {
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

		c.Request = req

		e.Use(TraceContextTraceIDMiddleware(traceparentHeaderKey))

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
