package decision

import "github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"

// After-GC heap ratio
func DecideJVMHeap(afterGCRatio float64) model.JVMDecision {
	if afterGCRatio > 0.80 {
		return model.JVMIncrease
	}
	return model.JVMKeep
}

// Non-heap pressure relative to container memory
func DecideJVMNonHeap(nonHeapBytes, memRequestBytes int64) model.JVMDecision {
	if memRequestBytes == 0 {
		return model.JVMKeep
	}

	if float64(nonHeapBytes)/float64(memRequestBytes) > 0.30 {
		return model.JVMIncrease
	}

	return model.JVMKeep
}
