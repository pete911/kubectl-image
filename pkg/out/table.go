package out

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"text/tabwriter"
)

type Table struct {
	logger     *slog.Logger
	writer     *tabwriter.Writer
	maxColSize int
}

func NewTable(logger *slog.Logger, maxColSize int) Table {
	return Table{
		logger:     logger,
		writer:     tabwriter.NewWriter(os.Stdout, 1, 1, 2, ' ', 0),
		maxColSize: maxColSize,
	}
}

func (t Table) AddRow(columns ...string) {
	if _, err := fmt.Fprintln(t.writer, strings.Join(t.resizeColumns(columns), "\t")); err != nil {
		t.logger.Error(fmt.Sprintf("table: add row: %v", err))
	}
}

func (t Table) resizeColumns(columns []string) []string {
	if t.maxColSize < 1 {
		return columns
	}
	var out []string
	for i := range columns {
		if len(columns[i]) > t.maxColSize {
			column := columns[i][:t.maxColSize-3]
			out = append(out, fmt.Sprintf("%s..", column))
			continue
		}
		out = append(out, columns[i])
	}
	return out
}

func (t Table) Print() {
	if err := t.writer.Flush(); err != nil {
		t.logger.Error(fmt.Sprintf("table: print: %v", err))
	}
}
