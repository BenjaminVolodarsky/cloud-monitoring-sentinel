package output

import (
	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
	"github.com/jedib0t/go-pretty/v6/text"
)

func colorDecision(d string) string {
	switch d {
	case "REDUCE":
		return text.FgGreen.Sprintf(d)
	case "KEEP":
		return text.FgYellow.Sprintf(d)
	case "INCREASE":
		return text.FgRed.Sprintf(d)
	case "SKIP_OOM", "SKIP_THROTTLING":
		return text.FgHiRed.Sprintf(d)
	default:
		return d
	}
}

func colorCPU(d model.CPUDecision) string {
	return colorDecision(string(d))
}

func colorMemory(d model.MemoryDecision) string {
	return colorDecision(string(d))
}

func colorJVM(d model.JVMDecision) string {
	return colorDecision(string(d))
}
