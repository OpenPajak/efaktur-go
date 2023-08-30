package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/OpenPajak/efaktur-go/pkg/provider/web"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"golang.org/x/exp/slices"
)

const (
	KetSuksesPosting     = "SUKSES POSTING"      // original
	KetSuksesLapor       = "SUKSES LAPOR"        // original
	KetSptSedangDiproses = "SPT SEDANG DIPROSES" // original: "SPT sedang diproses"
)

type (
	cmdDumpT    = cmd[dataCmdDump]
	dataCmdDump struct {
		Masa          int `vd:"$>=1 && $<=12"`
		Tahun         int `vd:"$>=2000"`
		JenisLampiran jenisLampiranFlag

		OutputType string `vd:"in($,'wide','yaml','json','csv')"`
		Annually   bool

		WaitSptHeaderCheckDuration time.Duration
		WorkerCount                int
		DumpPerPageCount           int
	}
)

var cmdDump = &cmdDumpT{
	name: "dump",
	data: dataCmdDump{
		Masa:  1,
		Tahun: 2000,
		JenisLampiran: jenisLampiranFlag{
			values: []web.KodeFormSpt{
				web.KodeFormSpt_A1,
				web.KodeFormSpt_A2,
				web.KodeFormSpt_B1,
				web.KodeFormSpt_B2,
				web.KodeFormSpt_B3,
			},
			pos: 0,
		},
		Annually:   false,
		OutputType: "wide",

		WaitSptHeaderCheckDuration: time.Second * 2,
		WorkerCount:                5,
		DumpPerPageCount:           10,
	},
	setup: func(fs *flag.FlagSet, data *dataCmdDump) {

		fs.IntVar(&data.Masa, "m", data.Masa, "Masa pajak")
		fs.IntVar(&data.Tahun, "y", data.Tahun, "Tahun pajak")
		fs.Var(&data.JenisLampiran, "t", "Lampiran type")
		fs.BoolVar(&data.Annually, "annually", data.Annually, "Annually mode")
		fs.StringVar(&data.OutputType, "o", data.OutputType, "Output type")

		fs.IntVar(&data.WorkerCount, "wrk", data.WorkerCount, "Worker SptHeader count")
		fs.DurationVar(&data.WaitSptHeaderCheckDuration, "waitchk", data.WaitSptHeaderCheckDuration, "Wait SptHeader check duration")
		fs.IntVar(&data.DumpPerPageCount, "ppc", data.DumpPerPageCount, "Dump per page count")
	},
	callback: func(ctx context.Context, c *cmdDumpT) (err error) {
		if err = vd.Validate(c.data); err != nil {
			return
		}

		defer cmdPrerun(ctx)()
		client := GetClientFromContext(ctx)

		const (
			OpCode_Delete = 1 << 0
			OpCode_Create = 1 << 1
		)
		type compactSptHeader struct {
			Masa       int
			Tahun      int
			RevisionNo int

			OpCode int

			Header *web.SptHeader
		}

		// grab existing SptHeader
		var respList web.SptHeaderListResponse
		if respList, err = client.SptHeader.List(ctx, web.SptHeaderListRequest{
			TahunPajak: c.data.Tahun,
		}); err != nil {
			return
		}

		var taskCreate = map[uint]*compactSptHeader{}
		// Check NON EXIST masa pajak
		// pack (year: 16-bit, month: 8-bit)
		var tabYMPb = map[uint]int{}
		for _, mp := range respList.Data {
			mState := strings.ToUpper(strings.Trim(mp.Keterangan, " "))

			tahun, masa, pb := mp.Tahun, mp.Masa1, mp.RevisionNo
			key := packYM(uint16(tahun), uint8(masa))

			task := &compactSptHeader{
				Masa:  masa,
				Tahun: tahun,

				Header: mp,
			}

			switch mState {
			case KetSuksesPosting:
				task.RevisionNo = mp.RevisionNo
				task.OpCode = OpCode_Delete | OpCode_Create
			case KetSuksesLapor:
				// temp pb to retrieve previous reported data.
				task.RevisionNo = mp.RevisionNo + 1
				task.OpCode = OpCode_Create
			case KetSptSedangDiproses:
				log.Printf("warning: SPT sedang diproses %+#v", mp)

			default:
				return fmt.Errorf("unhandled keterangan: %q", mState)
			}

			// previous pb is higher
			if prevPb, exist := tabYMPb[key]; exist && prevPb > pb {
				continue
			}

			// if it's not annually
			if !c.data.Annually && (c.data.Tahun != tahun || c.data.Masa != masa) {
				continue
			}

			taskCreate[key] = task
			tabYMPb[key] = pb
		}

		if c.data.Annually {
			currentMonth := int(time.Now().Month())
			// fill missing month
			for month := 1; month <= currentMonth; month++ {
				year := c.data.Tahun
				key := packYM(uint16(year), uint8(month))
				if _, exist := tabYMPb[key]; exist {
					continue
				}
				// create new one
				task := &compactSptHeader{
					Masa:       month,
					Tahun:      year,
					RevisionNo: 0,
					OpCode:     OpCode_Create,
				}
				taskCreate[key] = task
				tabYMPb[key] = task.RevisionNo
			}
		}

		// --- parallel ---
		type taskCreateCh struct {
			pYM uint
			cr  *compactSptHeader
		}
		chTaskCreate := make(chan taskCreateCh, c.data.WorkerCount)
		var wgTaskCreate sync.WaitGroup
		for i := 0; i < c.data.WorkerCount; i++ {
			go func() {
				for task := range chTaskCreate {
					pYM, cr := task.pYM, task.cr
					if cr.OpCode == 0 {
						continue
					}
					if err := func() (err error) { // shadow err
						defer wgTaskCreate.Done()

						year, month := unpackYM(pYM)
						_, _ = year, month
						// Delete OP
						if cr.OpCode&OpCode_Delete != 0 && cr.Header != nil {
							hdr := cr.Header
							log.Printf("DELETE {%d %d %d}", year, month, hdr.RevisionNo)
							var resp web.SptHeaderDeleteResponse
							if resp, err = client.SptHeader.Delete(ctx, web.SptHeaderDeleteRequest{
								SptHeader: *hdr,
							}); err != nil {
								return
							}
							_ = resp
						}
						// Create OP
						if cr.OpCode&OpCode_Create != 0 {
							log.Printf("CREATE {%d %d %d}", year, month, cr.RevisionNo)
							var resp web.SptHeaderCreateResponse
							if resp, err = client.SptHeader.Create(ctx, web.SptHeaderCreateRequest{
								Masa1:      cr.Masa,
								RevisionNo: cr.RevisionNo,
								Tahun:      cr.Tahun,
							}); err != nil {
								return
							}
							_ = resp
						}

						return
					}(); err != nil {
						log.Printf("worker SptCreate caught an err: %s", err)
					}
				}
			}()
		}
		// enqueue task to worker
		wgTaskCreate.Add(len(taskCreate))
		for pYM, cr := range taskCreate {
			chTaskCreate <- taskCreateCh{pYM, cr}
		}
		wgTaskCreate.Wait()
		close(chTaskCreate)
		// --- ---

		// periodically check the SptHeader list to see what's complete
		var completionTabHdr map[uint]*web.SptHeader
		for {
			var checkTabHdr = map[uint]*web.SptHeader{}
			log.Printf("checking header list")
			var respList web.SptHeaderListResponse
			if respList, err = client.SptHeader.List(ctx, web.SptHeaderListRequest{
				TahunPajak: c.data.Tahun,
			}); err != nil {
				return
			}

			for _, hdr := range respList.Data {
				year, month := hdr.Tahun, hdr.Masa1
				key := packYM(uint16(year), uint8(month))
				// match only successful posting
				if strings.ToUpper(strings.Trim(hdr.Keterangan, " ")) != KetSuksesPosting {
					log.Printf("UNMET_CRITERIA_0 ({%d %d %d}) {ket: %q}",
						year, month, hdr.RevisionNo,
						hdr.Keterangan,
					)
					continue
				}
				// put to table
				if hdr2, exist := checkTabHdr[key]; !exist || hdr.RevisionNo > hdr2.RevisionNo {

					// if it's not annually
					if !c.data.Annually && (c.data.Tahun != year || c.data.Masa != month) {
						continue
					}

					checkTabHdr[key] = hdr
				}
			}

			// match with taskCreate
			for key, pb := range tabYMPb {
				year, month := unpackYM(key)
				if hdr, exist := checkTabHdr[key]; !exist || hdr.RevisionNo != pb {
					var hpb *int
					if hdr != nil {
						hpb = &hdr.RevisionNo
					}
					log.Printf("UNMET_CRITERIA_1 ({%d %d %d}) {%d %d %d}",
						year, month, pb,
						year, month, hpb)
					goto retryWait
				}

			}

			// all header created.
			completionTabHdr = checkTabHdr // move ready list to completion.
			break

		retryWait:
			dur := c.data.WaitSptHeaderCheckDuration
			log.Printf("cooldown for %v", dur)
			time.Sleep(dur)
		}

		// ==== CLEAN UP ====
		// // clean up resource
		// defer func() {
		// }()

		// ===== LOOKUP =====

		type taskLookupLampiran struct {
			hdr           *web.SptHeader
			JenisLampiran web.KodeFormSpt

			PageNum  int
			PageSize int

			results  []*web.SptDetail
			callback func(task taskLookupLampiran)
		}
		var taskLookup []taskLookupLampiran
		for _, hdr := range completionTabHdr {
			year, month, pb := hdr.Tahun, hdr.Masa1, hdr.RevisionNo

			lampTypes := c.data.JenisLampiran.GetValues()
			log.Printf("getting spt detail header {%d %d %d} %v", year, month, pb, lampTypes)
			for _, lampType := range lampTypes {
				var resp web.SptDetailHeaderFindResponse
				if resp, err = client.SptDetailHeader.Find(ctx, web.SptDetailHeaderFindRequest{
					MasaPajak:   month,
					TahunPajak:  year,
					RevisionNo:  pb,
					KodeFormSpt: lampType,
				}); err != nil {
					return
				}

				// --- scroll paging ---
				pageCount := resp.Data.JumlahRecord/c.data.DumpPerPageCount + 1
				for pageNum := 1; pageNum <= pageCount; pageNum++ {
					taskLookup = append(taskLookup, taskLookupLampiran{
						hdr:           hdr,
						JenisLampiran: lampType,
						PageNum:       pageNum,
						PageSize:      c.data.DumpPerPageCount,
					})
				}
			}

		}
		var chTaskLookup = make(chan taskLookupLampiran, c.data.WorkerCount)
		var wgTaskLookup sync.WaitGroup
		for i := 0; i < c.data.WorkerCount; i++ {
			go func() {
				for task := range chTaskLookup {
					if err := func() (err error) {
						defer wgTaskLookup.Done()

						var resp web.SptDetailFindResponse
						if resp, err = client.SptDetail.Find(ctx, web.SptDetailFindRequest{
							MasaPajak:   task.hdr.Masa1,
							TahunPajak:  task.hdr.Tahun,
							RevisionNo:  task.hdr.RevisionNo,
							KodeFormSpt: task.JenisLampiran,
							PageNum:     task.PageNum,
							PageSize:    task.PageSize,
						}); err != nil {
							return
						}

						task.results = resp.Data
						if f := task.callback; f != nil {
							f(task)
						}
						return
					}(); err != nil {
						log.Printf("worker SptLampiranDump caught an err: %s", err)
					}
				}
			}()
		}

		var (
			results   = []*taskLookupLampiran{}
			resultsMu sync.Mutex
		)
		wgTaskLookup.Add(len(taskLookup))
		for _, lk := range taskLookup {
			lk.callback = func(task taskLookupLampiran) {
				resultsMu.Lock()
				defer resultsMu.Unlock()

				results = append(results, &task)
			}
			chTaskLookup <- lk
		}
		wgTaskLookup.Wait()
		close(chTaskLookup)

		// ===== Formatter =====
		slices.SortFunc(results, func(a *taskLookupLampiran, b *taskLookupLampiran) bool {
			return a.hdr.Tahun <= b.hdr.Tahun &&
				a.hdr.Masa1 <= b.hdr.Masa1 &&
				a.JenisLampiran <= b.JenisLampiran
		})
		ff := GetFormatter(c.data.OutputType)
		ff.SetHeaderStrings([]string{
			"PAJAK_TAHUN", "PAJAK_MASA",
			"JENIS_LAMPIRAN",
			//
			"NAMA", "NPWP",
			"NOMOR", "TANGGAL",
			"JUMLAH_DPP", "JUMLAH_PPN", "JUMLAH_PPNBM",
			"NO_REF",
			"KETERANGAN",
		})
		for _, meta := range results {

			for _, ent := range meta.results {
				ff.Add([]any{
					meta.hdr.Tahun, meta.hdr.Masa1,
					meta.JenisLampiran,
					//
					ent.NamaLt, ent.NpwpPasporLt,
					ent.Nomor, ent.Tanggal,
					ent.JumlahDpp, ent.JumlahPpn, ent.JumlahPpnbm,
					ent.NoRef,
					ent.Keterangan,
				})
			}

		}
		if _, err = ff.WriteTo(os.Stdout); err != nil {
			return
		}
		return
	},
}

// Pack year (16-bit) and month (8-bit) to uint
func packYM(year uint16, month uint8) uint {
	var v uint
	v |= uint(year&0xffff) << 8
	v |= uint(month & 0xff)
	return v
}

// Unpack uint as (year: 16-bit, month: 8-bit)
func unpackYM(ym uint) (year uint16, month uint8) {
	year = uint16((ym >> 8) & 0xffff)
	month = uint8(ym & 0xff)
	return
}

type jenisLampiranFlag struct {
	values []web.KodeFormSpt
	pos    int
}

var _ flag.Value = (*jenisLampiranFlag)(nil)

func (j *jenisLampiranFlag) Set(v string) (err error) {
	if j.pos == 0 { // reset len
		j.values = j.values[:0]
	}
	parts := strings.Split(v, ",")
	var vs web.KodeFormSpt
	for _, part := range parts {
		part = strings.Trim(part, " ")
		if vs, err = web.KodeFormSptString(part); err != nil {
			return fmt.Errorf("arg %d, %q: %w", j.pos, v, err)
		}
		j.values = append(j.values, vs)
	}
	j.pos++
	return
}

func (j *jenisLampiranFlag) GetValues() []web.KodeFormSpt {
	return j.values
}

func (j jenisLampiranFlag) String() string {
	var sb strings.Builder
	for i, val := range j.values {
		sb.WriteString(val.String())
		if i+1 < len(j.values) {
			sb.WriteRune(',')
		}
	}
	return sb.String()
}
