package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

var registeredFormatter = map[string]func() *Formatter{
	"wide": func() *Formatter {
		return NewFormatter(formatFuncWIDE)
	},
	"json": func() *Formatter {
		return NewFormatter(formatFuncJSON)
	},
	"yaml": func() *Formatter {
		return NewFormatter(formatFuncYAML)
	},
	"csv": func() *Formatter {
		return NewFormatter(formatFuncCSV)
	},
}

func GetFormatter(fmtName string) *Formatter {
	return registeredFormatter[fmtName]()
}

type FormatFunc func(f *Formatter, wr *bufio.Writer) (err error)
type FormatHeaderItem struct {
	Name string
	// map key hint
	KeyHint string
}
type Formatter struct {
	headers    []FormatHeaderItem
	ents       []any
	formatFunc FormatFunc
}

func NewFormatter(formatFunc FormatFunc) *Formatter {
	if formatFunc == nil {
		formatFunc = formatFuncWIDE
	}
	return &Formatter{
		headers:    []FormatHeaderItem{},
		ents:       []any{},
		formatFunc: formatFunc,
	}
}

func (f *Formatter) SetHeader(headers []FormatHeaderItem) {
	f.headers = headers
}

func (f *Formatter) SetHeaderStrings(headers []string) {
	var hdr []FormatHeaderItem
	for _, h := range headers {
		hdr = append(hdr, FormatHeaderItem{
			Name: h,
		})
	}
	f.SetHeader(hdr)
}

func (f *Formatter) Add(element any) {
	f.ents = append(f.ents, element)
}

func (f *Formatter) WriteTo(w io.Writer) (n int64, err error) {
	brw := bufio.NewWriter(w)
	if err = f.formatFunc(f, brw); err != nil {
		return
	}
	n = int64(brw.Size())
	err = brw.Flush()
	return
}

func zipSlice2Map(keys, values []any) map[string]any {
	if len(keys) > len(values) {
		return nil
	}
	m := map[string]any{}
	for i, key := range keys {
		val := values[i]
		m[key.(string)] = val
	}
	return m
}

func formatFuncCSV(f *Formatter, wr *bufio.Writer) (err error) {
	return formatFuncBuilder(f, wr, true, func(wr *bufio.Writer, data [][]any) (err error) {
		writer := csv.NewWriter(wr)
		writer.Comma = ';'
		header, body := data[0], data[1:]
		var headerStr []string
		for _, hdr := range header {
			headerStr = append(headerStr, hdr.(string))
		}
		if err = writer.Write(headerStr); err != nil {
			return
		}

		for _, cols := range body {
			var colStrs []string
			for _, col := range cols {
				if col == nil {
					colStrs = append(colStrs, "")
					continue
				}
				colStrs = append(colStrs, fmt.Sprintf("%v", col))
			}
			if err = writer.Write(colStrs); err != nil {
				return
			}
		}
		return
	})
}

func formatFuncJSON(f *Formatter, wr *bufio.Writer) (err error) {
	return formatFuncBuilder(f, wr, true, func(wr *bufio.Writer, data [][]any) (err error) {
		codec := json.NewEncoder(wr)
		codec.SetIndent("", "  ")
		header, body := data[0], data[1:]

		var eData = []map[string]any{}
		for _, cols := range body {
			m := zipSlice2Map(header, cols)
			eData = append(eData, m)
		}
		if err = codec.Encode(eData); err != nil {
			return
		}
		return
	})
}

func formatFuncYAML(f *Formatter, wr *bufio.Writer) (err error) {
	return formatFuncBuilder(f, wr, true, func(wr *bufio.Writer, data [][]any) (err error) {
		codec := yaml.NewEncoder(wr)
		codec.SetIndent(4)
		header, body := data[0], data[1:]

		var eData = []map[string]any{}
		for _, cols := range body {
			m := zipSlice2Map(header, cols)
			eData = append(eData, m)
		}
		if err = codec.Encode(eData); err != nil {
			return
		}
		return
	})
}

func formatFuncWIDE(f *Formatter, wr *bufio.Writer) (err error) {
	return formatFuncBuilder(f, wr, true, func(wr *bufio.Writer, data [][]any) (err error) {
		tab := tablewriter.NewWriter(wr)
		header, body := data[0], data[1:]
		var (
			headerStr []string
			bodyStr   [][]string
		)
		for _, hdr := range header {
			headerStr = append(headerStr, hdr.(string))
		}
		for _, cols := range body {
			var colStrs = make([]string, 0, len(cols))
			for _, col := range cols {
				colStrs = append(colStrs, fmt.Sprintf("%v", col))
			}
			bodyStr = append(bodyStr, colStrs)
		}

		tab.SetHeader(headerStr)
		tab.AppendBulk(bodyStr)
		tab.Render()
		return
	})
}

func formatFuncBuilder(f *Formatter, wr *bufio.Writer, withHeader bool, writerFn func(wr *bufio.Writer, data [][]any) (err error)) (err error) {
	var rows [][]any
	var hcols []any
	var hmapref []string // TODO: guidance for accessing map entries.
	if hdr := f.headers; hdr != nil {
		for _, hc := range hdr {
			hcols = append(hcols, hc.Name)
			hmapref = append(hmapref, hc.KeyHint)
		}
		rows = append(rows, hcols)
	}
loopEnts:
	for _, ent := range f.ents {
		var cols []any

		vt := reflect.ValueOf(ent)
		tt := vt.Type()
		for tt.Kind() == reflect.Pointer {
			if vt.IsNil() {
				continue loopEnts
			}
			vt = vt.Elem()
			tt = vt.Type()
		}

		// only struct and map.
		if k := tt.Kind(); k != reflect.Slice && k != reflect.Struct && k != reflect.Map {
			continue loopEnts
		}

		switch tt.Kind() {
		case reflect.Slice:
			for i := 0; i < vt.Len() && i < len(hcols); i++ {
				cols = append(cols, vt.Index(i).Interface())
			}
		case reflect.Struct:
			for i := 0; i < tt.NumField() && i < len(hcols); i++ {
				fieldT := tt.Field(i)
				fieldV := vt.Field(i)
				_, _ = fieldT, fieldV

				cols = append(cols, fieldV.Interface())
			}
		case reflect.Map:
			it := vt.MapRange()
			for i := 0; it.Next() && i < len(hcols); i++ {
				ekt := it.Key()
				evt := it.Value()
				_, _ = ekt, evt

				cols = append(cols, evt.Interface())
			}
		}

		// pad
		pads := make([]any, len(hcols)-len(cols))
		// fillEmptyStr(&pads)
		cols = append(cols, pads...)

		rows = append(rows, cols)
	}

	err = writerFn(wr, rows)
	return
}

func fillEmptyStr(s *[]any) {
	for i := range *s {
		(*s)[i] = "null"
	}
}
