package dumper

import (
	"io"
	"strconv"

	"github.com/gobeam/stringy"
	"github.com/olekukonko/tablewriter"
)

type entry[K any] struct {
	name   string
	getter func(in *K) (string, tablewriter.Colors)
}

type entryType[T any] []entry[T]

type Table[K any] struct {
	entries entryType[K]
}

func (t *Table[K]) Field(name string, getter func(in *K) string) {
	entry := entry[K]{
		name: name,
		getter: func(in *K) (string, tablewriter.Colors) {
			return getter(in), tablewriter.Colors{}
		},
	}

	t.entries = append(t.entries, entry)
}

func (t *Table[K]) Rich(name string, getter func(in *K) (string, tablewriter.Colors)) {
	entry := entry[K]{
		name:   name,
		getter: getter,
	}

	t.entries = append(t.entries, entry)
}

func (t *Table[K]) Dump(out io.Writer, items []K) {
	if len(items) == 0 {
		return
	}

	headers := make([]string, 0, len(t.entries))
	for i := range t.entries {
		header := stringy.New(t.entries[i].name).SnakeCase().ToUpper()
		if i == 0 {
			header = header + " (" + strconv.Itoa(len(t.entries)) + ")"
		}

		headers = append(headers, header)

	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(headers)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(false)
	table.SetColumnSeparator(tablewriter.SPACE)
	table.SetCenterSeparator(tablewriter.SPACE)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	for _, i := range items {
		row := make([]string, 0, len(t.entries))
		col := make([]tablewriter.Colors, 0, len(t.entries))

		for _, f := range t.entries {
			v, c := f.getter(&i)

			row = append(row, v)
			col = append(col, c)
		}

		table.Rich(row, col)
	}

	table.Render()
}
