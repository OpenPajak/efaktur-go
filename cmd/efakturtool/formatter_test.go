package main

import (
	"os"
	"testing"
)

func TestFormatter(t *testing.T) {
	fmts := []FormatFunc{
		formatFuncJSON,
		formatFuncYAML,
		formatFuncWIDE,
	}
	for _, ff := range fmts {
		f := NewFormatter(ff)
		f.SetHeader([]FormatHeaderItem{{Name: "ID"}, {Name: "Name"}})
		f.Add(map[string]any{
			"id":   "111111111111111",
			"name": "AAAAAAAAAAAAAAA",
		})
		f.Add(map[string]int{
			"id":   0xDEAD_BEEF,
			"name": 0xF00D_BABE,
		})
		f.Add(struct {
			ID   string
			Name string
		}{
			ID:   "222222222222222",
			Name: "BBBBBBBBBBBBBBB",
		})

		var vmm = map[string]any{}
		f.Add(vmm)
		vmm["id"] = "333333333333333"
		vmm["name"] = "BBBBBBBBBBBBBBB"

		f.Add([]string{"444444444444444"})
		f.Add([]string{"444444444444444", "444444444444444"})

		f.Add([]int{11111, 22222})

		f.WriteTo(os.Stdout)
	}
}
