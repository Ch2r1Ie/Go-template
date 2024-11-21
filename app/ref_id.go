package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RefIDMiddleware(headerKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		refID := c.Request.Header.Get(headerKey)
		if refID == "" {
			slog.WarnContext(c.Request.Context(), "no reference", slog.String("header-key", headerKey))
			refID = uuid.NewString()
		}

		c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
		c.Next()
	}
}

func ForwardRefIDOption(r *http.Request, ctxs ...context.Context) {
	ctx := context.Background()
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}

	if refID, ok := ctx.Value(refIDKey).(string); ok {
		r.Header.Add(string(refIDKey), refID)
	}
}

func newRefIDContext(ctx context.Context, refID string) context.Context {
	return context.WithValue(ctx, refIDKey, refID)
}

func SetRefID(c *gin.Context, refID string) {
	c.Request.Context()
	c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), refIDKey, refID))
}

func RefID(c *gin.Context) string {
	if v, ok := c.Request.Context().Value(refIDKey).(string); ok {
		return v
	}
	return ""
}
