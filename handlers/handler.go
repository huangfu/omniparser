package handlers

import (
	"io"

	"github.com/jf-tech/omniparser/customfuncs"
	"github.com/jf-tech/omniparser/errs"
	"github.com/jf-tech/omniparser/header"
	"github.com/jf-tech/omniparser/transformctx"
)

// HandlerCtx is a context object for a schema handler.
type HandlerCtx struct {
	Name          string
	Header        header.Header
	Content       []byte
	CustomFuncs   customfuncs.CustomFuncs
	HandlerParams interface{}
}

// CreateHandlerFunc is a func that checks if a given schema is supported by its associated
// handler, and, if yes, parses the schema content, creates and initializes a new instance of
// its associated handler.
// If the given schema is not supported, errs.ErrSchemaNotSupported should be returned.
// Any other error returned will cause omniparser to fail entirely.
// Note, any non errs.ErrSchemaNotSupported error returned here should be errs.CtxAwareErr
// formatted (i.e. error should contain schema name and if possible error line number).
type CreateHandlerFunc func(ctx *HandlerCtx) (SchemaHandler, error)

// SchemaHandler is an interface representing a schema handler responsible for ingesting,
// processing and transforming input stream based on its given schema.
type SchemaHandler interface {
	// NewIngester returns an Ingester for an input stream.
	// Omniparser will not call NewIngester unless CreateHandler has returned supported.
	// Omniparser calls NewIngester when client supplies an input stream and is ready
	// for the parser to ingest/process/transform the input.
	NewIngester(ctx *transformctx.Ctx, input io.Reader) (Ingester, error)
}

// Ingester is an interface responsible for a given input stream ingestion, and transformation.
type Ingester interface {
	// Read is called repeatedly during the processing of an input stream. Each call it should return
	// one result object, called record. It's entirely up to the implementation of this interface/method
	// to decide whether internally it does all the processing all at once (such as in the very first call
	// of `Read()`) and only hands out one record object at a time, OR, processes and returns one record
	// for each call. However, the overall design principle of omniparser is to have streaming processing
	// capability so memory won't be a constraint when dealing with large input file. All built-in handler
	// and ingesters are done this way.
	Read() ([]byte, error)

	// IsContinuableError is called to determine if the error returned by Read is fatal or not. After Read
	// is called, the result record or error will be returned to caller. After caller consumes record or
	// error, omniparser needs to decide whether to continue the transform operation or not, based on
	// whether the last err is "continuable" or not.
	IsContinuableError(error) bool

	// CtxAwareErr interface is embedded to provide omniparser and custom functions a way to provide
	// context aware (such as input file name + line number) error formatting.
	errs.CtxAwareErr
}
