package dumper

import (
	"encoding/csv"
	"io"
	"strconv"

	"github.com/gobeam/stringy"
	"github.com/olekukonko/tablewriter"
)

type Row struct {
	Value  string
	Colors tablewriter.Colors
}
type column[K any] struct {
	name   string
	wide   bool
	getter func(in *K) Row
}

type columnType[T any] []column[T]

type TableConfig struct {
	CSV  bool
	Wide bool
}
type Table[K any] struct {
	Config  TableConfig
	entries columnType[K]
}

func (t *Table[K]) Column(name string, wide bool, getter func(in *K) Row) {
	entry := column[K]{
		name:   name,
		wide:   wide,
		getter: getter,
	}

	t.entries = append(t.entries, entry)
}

func (t *Table[K]) Dump(out io.Writer, items []K) error {
	if len(items) == 0 {
		return nil
	}

	headers := make([]string, 0, len(t.entries))
	for i := range t.entries {
		if !t.Config.Wide && t.entries[i].wide {
			continue
		}

		header := stringy.New(t.entries[i].name).SnakeCase().ToUpper()
		headers = append(headers, header)

	}

	if t.Config.CSV {
		w := csv.NewWriter(out)

		w.Write(headers)

		for _, i := range items {
			row := make([]string, 0, len(t.entries))

			for _, f := range t.entries {
				r := f.getter(&i)

				row = append(row, r.Value)
			}

			w.Write(row)
		}

		w.Flush()

		return w.Error()
	} else {
		headers[0] = headers[0] + " (" + strconv.Itoa(len(t.entries)) + ")"

		table := tablewriter.NewWriter(out)
		table.SetHeader(headers)
		table.SetBorder(false)
		table.SetAutoFormatHeaders(false)
		table.SetRowLine(false)
		table.SetColumnSeparator(tablewriter.SPACE)
		table.SetCenterSeparator(tablewriter.SPACE)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAlignment(tablewriter.ALIGN_LEFT)

		for _, item := range items {
			row := make([]string, 0, len(t.entries))
			col := make([]tablewriter.Colors, 0, len(t.entries))

			for i, f := range t.entries {
				if !t.Config.Wide && t.entries[i].wide {
					continue
				}

				r := f.getter(&item)

				row = append(row, r.Value)
				col = append(col, r.Colors)
			}

			table.Rich(row, col)
		}

		table.Render()
	}

	return nil
}
