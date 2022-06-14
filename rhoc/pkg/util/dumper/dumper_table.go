package dumper

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"github.com/gobeam/stringy"
	"github.com/olekukonko/tablewriter"
)

type TableStyle int

const (
	TableStyleDefault TableStyle = iota
	TableStyleCSV
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

type ColumnList[T any] []Column[T]

type TableConfig[T any] struct {
	Style   TableStyle
	Wide    bool
	Columns ColumnList[T]
}

func DumpTable[T any](config TableConfig[T], out io.Writer, items []T) error {
	if len(items) == 0 {
		return nil
	}

	headers := make([]string, 0, len(config.Columns))
	for i := range config.Columns {
		if !config.Wide && config.Columns[i].Wide {
			continue
		}

		header := stringy.New(config.Columns[i].Name).SnakeCase().ToUpper()
		headers = append(headers, header)

	}

	switch config.Style {
	case TableStyleCSV:
		w := csv.NewWriter(out)

		w.Write(headers)

		for _, i := range items {
			row := make([]string, 0, len(config.Columns))

			for _, f := range config.Columns {
				r := f.Getter(&i)

				row = append(row, r.Value)
			}

			w.Write(row)
		}

		w.Flush()

		return w.Error()
	case TableStyleDefault:
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
		table.SetAutoWrapText(false)

		for _, item := range items {
			row := make([]string, 0, len(config.Columns))
			col := make([]tablewriter.Colors, 0, len(config.Columns))

			for i, f := range config.Columns {
				if !config.Wide && config.Columns[i].Wide {
					continue
				}

				r := f.Getter(&item)

				row = append(row, r.Value)
				col = append(col, r.Colors)
			}

			table.Rich(row, col)
		}

		table.Render()

		return nil
	default:
		return fmt.Errorf("unsupported table style %d", config.Style)
	}
}
