package output

import (
	"math"

	"github.com/jedib0t/go-pretty/v6/text"
)

func formatCPUChange(req, rec float64, decision string) string {
	if decision == "SKIP_THROTTLING" {
		return text.FgHiRed.Sprintf("⏭ SKIP")
	}

	delta := rec - req
	if math.Abs(delta) < 0.001 {
		return text.FgYellow.Sprintf("→ %.2f", rec)
	}

	if delta < 0 {
		return text.FgGreen.Sprintf("↓ %.2f (%.2f)", rec, delta)
	}

	return text.FgRed.Sprintf("↑ %.2f (+%.2f)", rec, delta)
}

func formatMemChange(reqBytes, recBytes int64, decision string) string {
	const MiB = 1024 * 1024

	if decision == "SKIP_OOM" {
		return text.FgHiRed.Sprintf("⏭ SKIP")
	}

	deltaMiB := float64(recBytes-reqBytes) / MiB
	recMiB := float64(recBytes) / MiB

	if math.Abs(deltaMiB) < 1 {
		return text.FgYellow.Sprintf("→ %.0fMi", recMiB)
	}

	if deltaMiB < 0 {
		return text.FgGreen.Sprintf("↓ %.0fMi (%.0f)", recMiB, deltaMiB)
	}

	return text.FgRed.Sprintf("↑ %.0fMi (+%.0f)", recMiB, deltaMiB)
}
