package output

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
)

func WriteCSV(path string, results []model.RightsizeResult, meta model.RightsizeMeta) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// ------------------------------------------------------------------
	// Metadata (key/value rows)
	// ------------------------------------------------------------------

	_ = w.Write([]string{"namespace", meta.Namespace})
	_ = w.Write([]string{"cluster", meta.Cluster})
	_ = w.Write([]string{"window", meta.Window})
	_ = w.Write([]string{"oom_window", meta.OOMWindow})
	_ = w.Write([]string{"target_util", fmt.Sprintf("%f", meta.TargetUtil)})
	_ = w.Write([]string{"safety_factor", fmt.Sprintf("%f", meta.SafetyFactor)})
	_ = w.Write([]string{}) // blank line

	// ------------------------------------------------------------------
	// Header
	// ------------------------------------------------------------------

	_ = w.Write([]string{
		"container",
		"mem_p95_ratio",
		"cpu_p95_ratio",
		"mem_request_bytes",
		"mem_recommended_bytes",
		"cpu_request_cores",
		"cpu_recommended_cores",
		"oom_killed",
		"cpu_decision",
		"memory_decision",
		"jvm_heap_decision",
		"jvm_non_heap_decision",
	})

	// ------------------------------------------------------------------
	// Rows
	// ------------------------------------------------------------------

	for _, r := range results {
		_ = w.Write([]string{
			r.Container,
			fmt.Sprintf("%f", r.MemP95Ratio),
			fmt.Sprintf("%f", r.CpuP95Ratio),
			fmt.Sprintf("%d", r.MemRequestBytes),
			fmt.Sprintf("%d", r.MemRecommendedBytes),
			fmt.Sprintf("%f", r.CpuRequestCores),
			fmt.Sprintf("%f", r.CpuRecommendedCores),
			fmt.Sprintf("%t", r.OOMKilled),
			string(r.CPUDecision),
			string(r.MemoryDecision),
			string(r.JVMHeapDecision),
			string(r.JVMNonHeapDecision),
		})
	}

	return w.Error()
}
