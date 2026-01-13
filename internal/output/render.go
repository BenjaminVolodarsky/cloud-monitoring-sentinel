package output

import (
	"fmt"
	"os"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
	"github.com/jedib0t/go-pretty/v6/table"
)

func RenderTable(results []model.RightsizeResult) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Header
	t.AppendHeader(table.Row{
		"CONTAINER",
		"MEM P95",
		"CPU P95",
		"MEM REQ",
		"MEM REC",
		"CPU REQ",
		"CPU REC",
		"CPU DECISION",
		"MEM DECISION",
		"JVM HEAP",
		"JVM NON-HEAP",
	})

	// Style
	t.SetStyle(table.Style{
		Name:    "upctl",
		Box:     table.StyleBoxRounded,
		Options: table.Options{DrawBorder: true, SeparateRows: true},
	})

	for _, r := range results {
		t.AppendRow(table.Row{
			r.Container,
			fmt.Sprintf("%.2f", r.MemP95Ratio),
			fmt.Sprintf("%.2f", r.CpuP95Ratio),
			bytes(r.MemRequestBytes),
			formatMemChange(
				r.MemRequestBytes,
				r.MemRecommendedBytes,
				string(r.MemoryDecision),
			),
			fmt.Sprintf("%.2f", r.CpuRequestCores),
			formatCPUChange(
				r.CpuRequestCores,
				r.CpuRecommendedCores,
				string(r.CPUDecision),
			),
			colorCPU(r.CPUDecision),
			colorMemory(r.MemoryDecision),
			colorJVM(r.JVMHeapDecision),
			colorJVM(r.JVMNonHeapDecision),
		})

	}

	t.Render()
}

func bytes(b int64) string {
	const (
		KiB = 1024
		MiB = KiB * 1024
		GiB = MiB * 1024
	)

	switch {
	case b >= GiB:
		return fmt.Sprintf("%.1fGi", float64(b)/GiB)
	case b >= MiB:
		return fmt.Sprintf("%.0fMi", float64(b)/MiB)
	default:
		return fmt.Sprintf("%dB", b)
	}
}
