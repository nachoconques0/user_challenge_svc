package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TraceIDKey string

const traceIDKey TraceIDKey = "trace_id"

// TraceIDMiddleware generates a trace_id per request and injects it into Gin context and Go context
func TraceIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.New().String()
		c.Set(string(traceIDKey), traceID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), traceIDKey, traceID))
		c.Header("X-Trace-ID", traceID)
		c.Next()
	}
}
