package decision

import (
	"fmt"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
)

func DecideCPU(p95Ratio float64, throttled bool) (model.CPUDecision, string) {
	if throttled {
		return model.CPUSkipThrottling, "CPU throttling detected (skipping reductions)"
	}
	if p95Ratio <= 0 {
		return model.CPUKeep, "no cpu ratio data (keeping)"
	}
	if p95Ratio > 0.90 {
		return model.CPUIncrease, fmt.Sprintf("cpu p95 ratio %.2f > 0.90 (pressure)", p95Ratio)
	}
	if p95Ratio < 0.60 {
		return model.CPUReduce, fmt.Sprintf("cpu p95 ratio %.2f < 0.60 (overprovisioned)", p95Ratio)
	}
	return model.CPUKeep, fmt.Sprintf("cpu p95 ratio %.2f within healthy band", p95Ratio)
}
