package app

import (
	crand "crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const traceparentHeaderKey = "traceparent"
const supportedVersion = 0
const TraceFlagsSampled = 3

type TraceID [16]byte
type SpanID [8]byte
type TraceFlags byte

func (id SpanID) String() string {
	return hex.EncodeToString(id[:])
}

func (id TraceID) String() string {
	return hex.EncodeToString(id[:])
}

type TraceParent struct {
	TraceID    TraceID
	SpanID     SpanID
	TraceFlags TraceFlags
}

var (
	ErrEmptyTraceParent = errors.New("empty traceparent")
	ErrWrongTraceParent = errors.New("wrong traceparent")
)

func NewTraceParent() TraceParent {
	sp := TraceParent{}
	id := defIDGenerator()
	sp.TraceID = id.NewTraceID()
	sp.SpanID = id.NewSpanID()

	return sp
}

func (tp TraceParent) String() string {
	return fmt.Sprintf("%.2x-%s-%.16x-%.2x",
		supportedVersion,
		tp.TraceID,
		tp.SpanID,
		tp.TraceFlags&TraceFlagsSampled)
}

var emptyTrace = TraceParent{}

func Parse(parent string) (TraceParent, error) {
	if parent == "" {
		return emptyTrace, ErrEmptyTraceParent
	}

	token := strings.Split(parent, "-")
	if len(token) != 4 {
		return emptyTrace, ErrWrongTraceParent
	}

	traceIDBytes, err := hex.DecodeString(token[1])
	if err != nil {
		return emptyTrace, err
	}

	spanIDBytes, err := hex.DecodeString(token[2])
	if err != nil {
		return emptyTrace, err
	}

	traceFlagsBytes, err := hex.DecodeString(token[3])
	if err != nil {
		return emptyTrace, err
	}

	var sp TraceParent
	copy(sp.TraceID[:], traceIDBytes[:16])
	copy(sp.SpanID[:], spanIDBytes[:8])
	sp.TraceFlags = TraceFlags(traceFlagsBytes[0])

	return sp, nil
}

func defIDGenerator() *defaultIDGenerator {
	gen := &defaultIDGenerator{}
	b := make([]byte, 32)
	_, err := crand.Read(b)
	if err != nil {
		slog.Error(err.Error())
	}
	gen.randSource = rand.NewChaCha8([32]byte(b))
	return gen
}

type defaultIDGenerator struct {
	sync.Mutex
	randSource *rand.ChaCha8
}

func (gen *defaultIDGenerator) NewSpanID() SpanID {
	gen.Lock()
	defer gen.Unlock()
	sid := SpanID{}
	_, err := gen.randSource.Read(sid[:])
	if err != nil {
		slog.Error(err.Error())
	}
	return sid
}

// NewTraceID returns a non-zero trace ID from a randomly-chosen sequence.
// mu should be held while this function is called.
func (gen *defaultIDGenerator) NewTraceID() TraceID {
	gen.Lock()
	defer gen.Unlock()
	tid := TraceID{}
	_, err := gen.randSource.Read(tid[:])
	if err != nil {
		slog.Error(err.Error())
	}
	return tid
}

func TraceContextTraceIDMiddleware(headerKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		refID := c.Request.Header.Get(headerKey)
		if refID == "" {
			slog.WarnContext(c.Request.Context(), "no reference", slog.String("header-key", headerKey))
			refID = NewTraceParent().TraceID.String()
		}

		tp, err := Parse(refID)
		if err != nil {
			refID = NewTraceParent().TraceID.String()
		} else {
			refID = tp.TraceID.String()
		}

		c.Request = c.Request.WithContext(newRefIDContext(c.Request.Context(), refID))
		c.Next()
	}
}
