package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"os"
	"statstable/cmd/models"
)

func main() {
	fs := flag.NewFlagSet("stats-table", flag.ExitOnError)

	repository := fs.String("repository", ".", "Path to the Git repository")
	revision := fs.String("revision", "HEAD", "Git revision to analyze")
	orderBy := fs.String("order-by", "lines", "Order by: lines, commits, files")
	useCommitter := fs.Bool("use-committer", false, "Use committer instead of author")
	format := fs.String("format", "tabular", "Output format: tabular, csv, json, json-lines")
	extensions := fs.StringSlice("extensions", []string{}, "File extensions to include, e.g., .go,.md")
	languages := fs.StringSlice("languages", []string{}, "Programming languages to include, e.g., go,markdown")
	exclude := fs.StringSlice("exclude", []string{}, "Glob patterns to exclude files")
	restrictTo := fs.StringSlice("restrict-to", []string{}, "Glob patterns to include files")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, "Error parsing flags:", err)
		if err2 != nil {
			return
		}
		os.Exit(1)
	}

	context := &models.Context{
		Repository:   *repository,
		Revision:     *revision,
		Extensions:   *extensions,
		OrderBy:      *orderBy,
		Format:       *format,
		Exclude:      *exclude,
		RestrictTo:   *restrictTo,
		Languages:    *languages,
		UseCommitter: *useCommitter,
	}

	chain := &models.BaseHandler{}
	chain.
		SetNext(&models.ValidateHandler{}).
		SetNext(&models.FileCollectionHandler{}).
		SetNext(&models.FileFilterHandler{}).
		SetNext(&models.BlameCollectorHandler{}).
		SetNext(&models.SortHandler{}).
		SetNext(&models.OutputHandler{})

	if err = chain.Handle(context); err != nil {
		_, err2 := fmt.Fprintln(os.Stderr, err)
		if err2 != nil {
			return
		}
		os.Exit(2)
	}
}
