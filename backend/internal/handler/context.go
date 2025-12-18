package handler

import (
	"searchav/internal/constants"
	"searchav/internal/dto"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

// Context is the enhanced request context
type Context struct {
	*fiber.Ctx
	Logger *zerolog.Logger
	Resp   *dto.Response
}

// Handler is the handler function type
type Handler func(*Context) error

// ContextHandler wraps handlers with enhanced context
type ContextHandler struct {
	logger *zerolog.Logger
}

// NewContextHandler creates a new context handler
func NewContextHandler(logger *zerolog.Logger) *ContextHandler {
	return &ContextHandler{
		logger: logger,
	}
}

// Wrap wraps a handler function with context
func (h *ContextHandler) Wrap(fn Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := &Context{
			Ctx:    c,
			Logger: h.logger,
			Resp:   &dto.Response{},
		}
		return fn(ctx)
	}
}

// Success returns a success response
func (ctx *Context) Success() error {
	code := constants.Success
	return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(code.Message))
}

// SuccessWithData returns a success response with data
func (ctx *Context) SuccessWithData(data interface{}) error {
	code := constants.Success
	return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(code.Message).WithData(data))
}

// SuccessWithList returns a success response with list
func (ctx *Context) SuccessWithList(list interface{}) error {
	code := constants.Success
	return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(code.Message).WithList(list))
}

// BadRequest returns a bad request response
func (ctx *Context) BadRequest(msg string) error {
	code := constants.InvalidParams
	if msg != "" {
		return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(msg))
	}
	return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(code.Message))
}

// InternalError returns an internal error response
func (ctx *Context) InternalError(err error) error {
	code := constants.InternalError
	ctx.Logger.Error().Err(err).Msg("internal error")
	return ctx.JSON(ctx.Resp.WithCode(code.Code).WithMessage(code.Message))
}
