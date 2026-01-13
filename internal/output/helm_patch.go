package output

import (
	"fmt"
	"os"
	"sort"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
)

func WriteHelmValuesPatch(path string, results []model.RightsizeResult) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Keep deterministic ordering
	sort.Slice(results, func(i, j int) bool {
		return results[i].Container < results[j].Container
	})

	// A generic values.yaml snippet keyed by "services.<name>.resources.requests"
	// Your platform team can align charts to consume this structure.
	_, _ = fmt.Fprintln(f, "# upctl-generated Helm values snippet")
	_, _ = fmt.Fprintln(f, "# Merge this into your chart values (or adapt to your chart schema).")
	_, _ = fmt.Fprintln(f, "services:")

	for _, r := range results {
		// Skip if both are KEEP and no mem/cpu change desired (optional). For now include all.
		cpu := cpuString(r.CpuRecommendedCores)
		mem := memString(r.MemRecommendedBytes)

		_, _ = fmt.Fprintf(f, "  %s:\n", r.Container)
		_, _ = fmt.Fprintln(f, "    resources:")
		_, _ = fmt.Fprintln(f, "      requests:")
		_, _ = fmt.Fprintf(f, "        cpu: %q\n", cpu)
		_, _ = fmt.Fprintf(f, "        memory: %q\n", mem)
	}

	return nil
}

func cpuString(cores float64) string {
	// Convert cores to millicores for readability
	m := int64(cores * 1000.0)
	return fmt.Sprintf("%dm", m)
}

func memString(bytes int64) string {
	// Convert bytes to Mi
	const MiB = 1024 * 1024
	mi := (bytes + MiB - 1) / MiB
	return fmt.Sprintf("%dMi", mi)
}
