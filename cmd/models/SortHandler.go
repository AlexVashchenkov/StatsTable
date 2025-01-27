package models

import (
	"errors"
	"sort"
)

type SortHandler struct {
	next Handler
}

func (h *SortHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *SortHandler) Handle(context *Context) error {
	switch context.OrderBy {
	case "lines", "files", "commits":

	default:
		return errors.New("bad sort-by")
	}

	context.StatsToPrint = make([]AuthorStatsForPrint, 0, len(context.Stats))

	for _, stat := range context.Stats {
		context.StatsToPrint = append(context.StatsToPrint, AuthorStatsForPrint{
			Name:    stat.Name,
			Lines:   stat.Lines,
			Commits: len(stat.Commits),
			Files:   len(stat.Files),
		})
	}

	sort.Slice(context.StatsToPrint, func(i, j int) bool {
		switch context.OrderBy {
		case "lines":
			if context.StatsToPrint[i].Lines == context.StatsToPrint[j].Lines {
				if context.StatsToPrint[i].Files == context.StatsToPrint[j].Files {
					if context.StatsToPrint[i].Commits == context.StatsToPrint[j].Commits {
						return context.StatsToPrint[i].Name < context.StatsToPrint[j].Name
					}
					return context.StatsToPrint[i].Commits > context.StatsToPrint[j].Commits
				}
				return context.StatsToPrint[i].Files > context.StatsToPrint[j].Files
			}
			return context.StatsToPrint[i].Lines > context.StatsToPrint[j].Lines
		case "commits":
			if context.StatsToPrint[i].Commits == context.StatsToPrint[j].Commits {
				if context.StatsToPrint[i].Lines == context.StatsToPrint[j].Lines {
					if context.StatsToPrint[i].Files == context.StatsToPrint[j].Files {
						return context.StatsToPrint[i].Name < context.StatsToPrint[j].Name
					}
					return context.StatsToPrint[i].Files > context.StatsToPrint[j].Files
				}
				return context.StatsToPrint[i].Lines > context.StatsToPrint[j].Lines
			}
			return context.StatsToPrint[i].Commits > context.StatsToPrint[j].Commits
		case "files":
			if context.StatsToPrint[i].Files == context.StatsToPrint[j].Files {
				if context.StatsToPrint[i].Lines == context.StatsToPrint[j].Lines {
					if context.StatsToPrint[i].Commits == context.StatsToPrint[j].Commits {
						return context.StatsToPrint[i].Name < context.StatsToPrint[j].Name
					}
					return context.StatsToPrint[i].Commits > context.StatsToPrint[j].Commits
				}
				return context.StatsToPrint[i].Lines > context.StatsToPrint[j].Lines
			}
			return context.StatsToPrint[i].Files > context.StatsToPrint[j].Files
		default:
			return context.StatsToPrint[i].Name < context.StatsToPrint[j].Name
		}
	})

	if h.next != nil {
		return h.next.Handle(context)
	} else {
		return nil
	}
}
