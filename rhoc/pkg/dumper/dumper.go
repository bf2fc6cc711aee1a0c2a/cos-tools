package dumper

import (
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gobeam/stringy"
	"github.com/olekukonko/tablewriter"
)

type Table struct {
	customizers map[string]func(string) tablewriter.Colors
}

func NewTable(customizers map[string]func(string) tablewriter.Colors) Table {
	return Table{
		customizers: customizers,
	}
}

func (t *Table) Dump(items []interface{}, out io.Writer) {
	if len(items) == 0 {
		return
	}

	r := reflect.ValueOf(items[0])

	h := make([]string, 0, r.NumField())
	for i := 0; i < r.NumField(); i++ {
		field := r.Type().Field(i)
		name := t.tag(field, "header")

		if name == "" {
			name = t.tag(field, "json")
		}
		if name == "" {
			name = stringy.New(field.Name).SnakeCase().ToUpper()
		}

		if i == 0 {
			name = name + " (" + strconv.Itoa(len(items)) + ")"
		}
		h = append(h, name)
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader(h)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowLine(false)
	table.SetColumnSeparator(tablewriter.SPACE)
	table.SetCenterSeparator(tablewriter.SPACE)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)

	for _, i := range items {
		structVal := reflect.ValueOf(i)

		row := make([]string, 0, r.NumField())
		col := make([]tablewriter.Colors, 0, r.NumField())

		for i := 0; i < r.NumField(); i++ {
			fieldName := structVal.Type().Field(i).Name
			fieldVal := structVal.Field(i).Interface()

			var data string

			switch v := fieldVal.(type) {
			case time.Time:
				data = v.Format(time.RFC3339)
			default:
				data = fmt.Sprintf("%v", fieldVal)
			}

			if c, ok := t.customizers[strings.ToLower(fieldName)]; ok {
				col = append(col, c(data))
			} else {
				col = append(col, tablewriter.Colors{})
			}

			row = append(row, data)
		}

		table.Rich(row, col)
	}

	table.Render()
}

func (t *Table) tag(field reflect.StructField, name string) string {
	tag := field.Tag.Get(name)
	if tag == "" {
		return ""
	}

	return strings.Split(tag, ",")[0]
}
