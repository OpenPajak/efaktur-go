package main

import (
	"context"
	"encoding/csv"
	"flag"
	"io"
	"os"
	"strconv"

	"github.com/OpenPajak/efaktur-go/pkg/provider/web"

	vd "github.com/bytedance/go-tagexpr/v2/validator"
)

type (
	cmdPrepopulatedT    = cmd[dataCmdPrepopulated]
	dataCmdPrepopulated struct {
		Masa         int    `vd:"$>=1 && $<=12"`
		Tahun        int    `vd:"$>=2000"`
		JenisDokumen string `vd:""`

		OutputType string `vd:"in($,'wide','yaml','json','csv')"`
		Annually   bool
	}
)

var cmdPrepopulated = &cmdPrepopulatedT{
	name: "prepopulated",
	data: dataCmdPrepopulated{
		Masa:         1,
		Tahun:        2000,
		JenisDokumen: "FPM",
		Annually:     false,
		OutputType:   "wide",
	},
	setup: func(fs *flag.FlagSet, data *dataCmdPrepopulated) {

		fs.IntVar(&data.Masa, "m", data.Masa, "Masa pajak")
		fs.IntVar(&data.Tahun, "y", data.Tahun, "Tahun pajak")
		fs.StringVar(
			&data.JenisDokumen,
			"t",
			data.JenisDokumen,
			"Document type")
		fs.BoolVar(&data.Annually, "annually", data.Annually, "Annually mode")
		fs.StringVar(&data.OutputType, "o", data.OutputType, "Output type")

	},
	callback: func(ctx context.Context, c *cmdPrepopulatedT) (err error) {
		if err = vd.Validate(c.data); err != nil {
			return
		}
		var docType web.PrepopulatedJenisDokumen
		if docType, err = web.PrepopulatedJenisDokumenString(c.data.JenisDokumen); err != nil {
			return
		}
		defer cmdPrerun(ctx)()
		client := GetClientFromContext(ctx)

		var masaPajak = [][2]int{}
		if !c.data.Annually {
			masaPajak = append(masaPajak, [2]int{c.data.Tahun, c.data.Masa})
		} else {
			for i := 1; i <= 12; i++ {
				masaPajak = append(masaPajak, [2]int{c.data.Tahun, i})
			}
		}

		type Meta struct {
			*csv.Reader
			Masa [2]int
		}

		var respMeta []Meta
		for _, masa := range masaPajak {
			var csvReader *csv.Reader
			year, month := masa[0], masa[1]
			if csvReader, err = client.Prepopulated.Download(ctx, web.PrepopulatedDownloadRequest{
				MasaPajak:    web.PrepopulatedMasaPajak(month),
				TahunPajak:   year,
				JenisDokumen: docType,
			}); err != nil {
				return
			}
			respMeta = append(respMeta, Meta{
				Reader: csvReader,
				Masa:   masa,
			})
		}

		var rowIdx int
		var header []string

		ff := GetFormatter(c.data.OutputType)

		for _, meta := range respMeta {
			var rec []string
			year, month := meta.Masa[0], meta.Masa[1]
			yearA, monthA := strconv.Itoa(year), strconv.Itoa(month)

			for i := 0; ; func() { rowIdx++; i++ }() {
				if rec, err = meta.Read(); err == io.EOF {
					err = nil
					break
				} else if err != nil {
					return
				}
				if i == 0 && header == nil {
					rec = append([]string{
						"PAJAK_TAHUN", "PAJAK_MASA",
					}, rec...)
					header = rec
					ff.SetHeaderStrings(header)
					continue
				} else if i == 0 {
					continue
				}

				rec = append([]string{
					yearA, monthA,
				}, rec...)

				ff.Add(rec)
			}
		}

		if _, err = ff.WriteTo(os.Stdout); err != nil {
			return
		}
		return
	},
}
