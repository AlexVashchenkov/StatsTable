package models

import (
	"fmt"
	"os/exec"
	"strings"
)

type BlameCollectorHandler struct {
	next Handler
}

func (h *BlameCollectorHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *BlameCollectorHandler) Handle(ctx *Context) error {
	statsMap := make(map[string]*AuthorStats)

	if len(ctx.Files) == 0 || len(ctx.Files) == 1 && ctx.Files[0] == "" {
		return h.next.Handle(ctx)
	}

	for _, file := range ctx.Files {
		cmd := exec.Command("git", "blame", "--line-porcelain", ctx.Revision, "--", file)
		cmd.Dir = ctx.Repository

		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("failed to run git blame on %s: %w", file, err)
		}
		var currentAuthor string
		var currentCommit string

		for _, line := range strings.Split(string(output), "\n") {
			if ctx.UseCommitter && strings.HasPrefix(line, "committer ") {
				currentAuthor = strings.TrimSpace(strings.TrimPrefix(line, "committer "))
			} else if !ctx.UseCommitter && strings.HasPrefix(line, "author ") {
				currentAuthor = strings.TrimSpace(strings.TrimPrefix(line, "author "))
			}

			if _, exists := statsMap[currentAuthor]; !exists {
				statsMap[currentAuthor] = &AuthorStats{
					Name:    currentAuthor,
					Lines:   0,
					Commits: make(map[string]struct{}),
					Files:   make(map[string]struct{}),
				}
			}

			if strings.HasPrefix(line, "committer-time ") {
				commit := strings.Fields(line)[1]
				currentCommit = commit
				statsMap[currentAuthor].Commits[commit] = struct{}{}
			}

			if strings.HasPrefix(line, "\t") {
				statsMap[currentAuthor].Lines++
				statsMap[currentAuthor].Files[file] = struct{}{}
			}
		}

		if currentCommit == "" {
			cmd = exec.Command("git", "log", "-1", "--format=%ct", ctx.Revision, "--", file)
			cmd.Dir = ctx.Repository

			logOutput, err := cmd.Output()
			if err != nil {
				return fmt.Errorf("failed to get last commit time for %s: %w", file, err)
			}

			commit := strings.TrimSpace(string(logOutput))
			if currentAuthor == "" {
				cmd = exec.Command("git", "log", "-1", "--format=%an", ctx.Revision, "--", file)
				if ctx.UseCommitter {
					cmd = exec.Command("git", "log", "-1", "--format=%cn", ctx.Revision, "--", file)
				}
				cmd.Dir = ctx.Repository

				logOutput, err := cmd.Output()
				if err != nil {
					return fmt.Errorf("failed to get last committer for %s: %w", file, err)
				}

				author := strings.TrimSpace(string(logOutput))
				if author == "" {
					continue
				}

				if _, exists := statsMap[author]; !exists {
					statsMap[author] = &AuthorStats{
						Name:    author,
						Lines:   0,
						Commits: make(map[string]struct{}),
						Files:   make(map[string]struct{}),
					}
				}
				statsMap[author].Files[file] = struct{}{}
				statsMap[author].Commits[commit] = struct{}{}
			}
		}
	}

	ctx.Stats = make([]AuthorStats, 0, len(statsMap))
	for name, stat := range statsMap {
		if name == "" || (stat.Lines == 0 && len(stat.Commits) == 0 && len(stat.Files) == 0) {
			continue
		}
		ctx.Stats = append(ctx.Stats, *stat)
	}

	return h.next.Handle(ctx)
}
