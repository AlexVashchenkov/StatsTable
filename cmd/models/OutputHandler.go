package models

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"text/tabwriter"
)

type OutputHandler struct {
	next Handler
}

func (h *OutputHandler) SetNext(handler Handler) Handler {
	h.next = handler
	return handler
}

func (h *OutputHandler) Handle(context *Context) error {
	switch context.Format {
	case "tabular":
		return h.printTabular(context.StatsToPrint)
	case "json":
		return h.printJSON(context.StatsToPrint)
	case "csv":
		return h.printCSV(context.StatsToPrint)
	case "json-lines":
		return h.printJSONLines(context.StatsToPrint)
	default:
		return errors.New("bad format")
	}
}

func (h *OutputHandler) printTabular(data []AuthorStatsForPrint) error {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', 0)
	if _, err := fmt.Fprintf(w, "Name\tLines\tCommits\tFiles\n"); err != nil {
		return err
	}

	for _, stat := range data {
		_, err := fmt.Fprintf(w, "%s\t%d\t%d\t%d\n", stat.Name, stat.Lines, stat.Commits, stat.Files)
		if err != nil {
			return err
		}
	}
	return w.Flush()
}

func (h *OutputHandler) printCSV(data []AuthorStatsForPrint) error {
	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	if err := writer.Write([]string{"Name", "Lines", "Commits", "Files"}); err != nil {
		return err
	}

	for _, stat := range data {
		row := []string{
			stat.Name,
			strconv.Itoa(stat.Lines),
			strconv.Itoa(stat.Commits),
			strconv.Itoa(stat.Files),
		}

		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

func (h *OutputHandler) printJSON(data []AuthorStatsForPrint) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (h *OutputHandler) printJSONLines(data []AuthorStatsForPrint) error {
	var buffer bytes.Buffer

	for _, stat := range data {
		line, err := json.Marshal(stat)
		if err != nil {
			return err
		}
		buffer.Write(line)
		buffer.WriteString("\n")
	}

	_, err := os.Stdout.Write(buffer.Bytes())
	return err
}
