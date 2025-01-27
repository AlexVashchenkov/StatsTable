package models

type Handler interface {
	Handle(context *Context) error
	SetNext(handler Handler) Handler
}

type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(handler Handler) Handler {
	b.next = handler
	return handler
}

func (b *BaseHandler) Handle(context *Context) error {
	if b.next != nil {
		return b.next.Handle(context)
	}
	return nil
}
