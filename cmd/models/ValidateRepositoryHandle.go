package models

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type ValidateHandler struct {
	BaseHandler
}

func (h *ValidateHandler) Handle(ctx *Context) error {
	if _, err := os.Stat(ctx.Repository); os.IsNotExist(err) {
		return fmt.Errorf("repository does not exist: %s", ctx.Repository)
	}

	cmd := exec.Command("git", "rev-parse", ctx.Revision)
	cmd.Dir = ctx.Repository
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to resolve revision: %w", err)
	}

	ctx.ResolvedRevision = strings.TrimSpace(string(output))
	return h.next.Handle(ctx)
}

func (h *ValidateHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}
