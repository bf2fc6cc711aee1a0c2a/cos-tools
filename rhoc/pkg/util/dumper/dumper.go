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
	table.SetAlignment(tablewriter.ALIGN_LEFT)

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
			case int64:
				data = strconv.FormatInt(v, 10)
			case int:
				data = strconv.Itoa(v)
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

type Tbl[K any] struct {
	fields []TblEntry[K]
}

type TblEntry[K any] struct {
	Name   string
	Getter func(in *K) (string, tablewriter.Colors)
}

func (t *Tbl[K]) Field(name string, getter func(in *K) string) {
	var e []TblEntry[K]
	e = t.fields

	entry := TblEntry[K]{
		Name: name,
		Getter: func(in *K) (string, tablewriter.Colors) {
			return getter(in), tablewriter.Colors{}
		},
	}

	e = append(e, entry)
	t.fields = e
}

func (t *Tbl[K]) Rich(name string, getter func(in *K) (string, tablewriter.Colors)) {
	var e []TblEntry[K]
	e = t.fields

	entry := TblEntry[K]{
		Name:   name,
		Getter: getter,
	}

	e = append(e, entry)
	t.fields = e
}

func (t *Tbl[K]) Dump(items []K, out io.Writer) {
	if len(items) == 0 {
		return
	}

	var e []TblEntry[K]
	e = t.fields

	h := make([]string, 0, len(e))
	for i := range t.fields {
		h = append(h, stringy.New(t.fields[i].Name).SnakeCase().ToUpper())
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
		row := make([]string, 0, len(e))
		col := make([]tablewriter.Colors, 0, len(e))

		for _, f := range t.fields {
			v, c := f.Getter(&i)

			row = append(row, v)
			col = append(col, c)
		}

		table.Rich(row, col)
	}

	table.Render()
}
