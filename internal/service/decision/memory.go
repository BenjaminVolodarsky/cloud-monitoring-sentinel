package decision

import (
	"fmt"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
)

func DecideMemory(p95Ratio float64, oomKilled bool) (model.MemoryDecision, string) {
	if oomKilled {
		return model.MemSkipOOM, "OOMKilled detected in lookback window"
	}
	if p95Ratio <= 0 {
		return model.MemKeep, "no memory ratio data (keeping)"
	}
	if p95Ratio > 0.90 {
		return model.MemIncrease, fmt.Sprintf("mem p95 ratio %.2f > 0.90 (risk)", p95Ratio)
	}
	if p95Ratio < 0.60 {
		return model.MemReduce, fmt.Sprintf("mem p95 ratio %.2f < 0.60 (overprovisioned)", p95Ratio)
	}
	return model.MemKeep, fmt.Sprintf("mem p95 ratio %.2f within healthy band", p95Ratio)
}
