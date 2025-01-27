package models

import (
	"fmt"
	"os/exec"
	"strings"
)

type FileCollectionHandler struct {
	next Handler
}

func (h *FileCollectionHandler) Handle(context *Context) error {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", context.Revision)
	cmd.Dir = context.Repository

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to list files: %w", err)
	}

	context.Files = strings.Split(strings.TrimSpace(string(output)), "\n")
	return h.next.Handle(context)
}

func (h *FileCollectionHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}
