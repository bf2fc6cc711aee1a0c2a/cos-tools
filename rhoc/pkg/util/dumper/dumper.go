package dumper

import (
	"io"

	"github.com/gobeam/stringy"
	"github.com/olekukonko/tablewriter"
	"github.com/redhat-developer/app-services-cli/pkg/core/ioutil/dump"
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

func (t *Table[K]) Dump(items []K, out io.Writer) {
	if len(items) == 0 {
		return
	}

	h := make([]string, 0, len(t.entries))
	for i := range t.entries {
		h = append(h, stringy.New(t.entries[i].name).SnakeCase().ToUpper())
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(h)
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

type Formatted[K any] struct {
	format string
}

func (f *Formatted[K]) Dump(data K, out io.Writer) error {
	return dump.Formatted(out, f.format, data)
}
