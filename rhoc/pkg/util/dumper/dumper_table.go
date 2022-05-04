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
type Column[K any] struct {
	Name   string
	Wide   bool
	Getter func(in *K) Row
}

type columnType[T any] []Column[T]

type TableConfig struct {
	CSV  bool
	Wide bool
}
type Table[K any] struct {
	Config  TableConfig
	Columns columnType[K]
}

func (t *Table[K]) Dump(out io.Writer, items []K) error {
	if len(items) == 0 {
		return nil
	}

	headers := make([]string, 0, len(t.Columns))
	for i := range t.Columns {
		if !t.Config.Wide && t.Columns[i].Wide {
			continue
		}

		header := stringy.New(t.Columns[i].Name).SnakeCase().ToUpper()
		headers = append(headers, header)

	}

	if t.Config.CSV {
		w := csv.NewWriter(out)

		w.Write(headers)

		for _, i := range items {
			row := make([]string, 0, len(t.Columns))

			for _, f := range t.Columns {
				r := f.Getter(&i)

				row = append(row, r.Value)
			}

			w.Write(row)
		}

		w.Flush()

		return w.Error()
	} else {
		headers[0] = headers[0] + " (" + strconv.Itoa(len(items)) + ")"

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
			row := make([]string, 0, len(t.Columns))
			col := make([]tablewriter.Colors, 0, len(t.Columns))

			for i, f := range t.Columns {
				if !t.Config.Wide && t.Columns[i].Wide {
					continue
				}

				r := f.Getter(&item)

				row = append(row, r.Value)
				col = append(col, r.Colors)
			}

			table.Rich(row, col)
		}

		table.Render()
	}

	return nil
}
