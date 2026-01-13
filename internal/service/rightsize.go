package service

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/model"
	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/promql"
	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/service/decision"
	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/vm"
)

type RightsizeParams struct {
	Namespace    string
	Cluster      string
	Window       string
	SubqueryStep string

	OOMWindow string

	TargetUtil   float64
	SafetyFactor float64
	MemRoundMiB  int64
	CPURoundm    int64

	TopK   int
	Bottom bool
}

type RightsizeService struct {
	vm *vm.Client
}

func NewRightsizeService() *RightsizeService {
	const vmURL = "http://vmselect.management.prod.internal:8481/select/0/prometheus/api/v1"

	return &RightsizeService{
		vm: vm.NewClient(vmURL),
	}
}

func (s *RightsizeService) Run(
	ctx context.Context,
	p RightsizeParams,
) ([]model.RightsizeResult, model.RightsizeMeta, error) {

	meta := model.RightsizeMeta{
		Namespace:    p.Namespace,
		Cluster:      p.Cluster,
		Window:       p.Window,
		OOMWindow:    p.OOMWindow,
		TargetUtil:   p.TargetUtil,
		SafetyFactor: p.SafetyFactor,
		SubqueryStep: p.SubqueryStep,
	}

	// ---------------------------------------------------------------------
	// 1. Fetch signals (best-effort where appropriate)
	// ---------------------------------------------------------------------

	memP95, err := s.query(ctx, promql.MemP95Ratio(
		p.Namespace, p.Cluster, p.Window, p.SubqueryStep,
	))
	if err != nil {
		return nil, meta, fmt.Errorf("mem p95 ratio: %w", err)
	}

	cpuP95, err := s.query(ctx, promql.CpuP95Ratio(
		p.Namespace, p.Cluster, p.Window, p.SubqueryStep,
	))
	if err != nil {
		return nil, meta, fmt.Errorf("cpu p95 ratio: %w", err)
	}

	memReq, err := s.query(ctx, promql.MemRequests(p.Namespace, p.Cluster))
	if err != nil {
		return nil, meta, fmt.Errorf("mem requests: %w", err)
	}

	cpuReq, err := s.query(ctx, promql.CpuRequests(p.Namespace, p.Cluster))
	if err != nil {
		return nil, meta, fmt.Errorf("cpu requests: %w", err)
	}

	oom, _ := s.query(ctx, promql.OOMKilled(
		p.Namespace, p.Cluster, p.OOMWindow,
	))

	cpuThrottle, _ := s.query(ctx, promql.CPUThrottling(
		p.Namespace, p.Cluster, p.Window,
	))

	jvmHeapAfterGC, _ := s.query(ctx, promql.JVMHeapAfterGC(
		p.Namespace, p.Cluster,
	))

	jvmNonHeap, _ := s.query(ctx, promql.JVMNonHeapBytes(
		p.Namespace, p.Cluster,
	))

	// ---------------------------------------------------------------------
	// 2. Index all signals by (namespace|cluster|container)
	// ---------------------------------------------------------------------

	key := func(m map[string]string) string {
		return m["namespace"] + "|" + m["uw_cluster"] + "|" + m["container"]
	}

	memP95Map := map[string]float64{}
	for _, s := range memP95 {
		memP95Map[key(s.Metric)] = s.ValueFloat
	}

	cpuP95Map := map[string]float64{}
	for _, s := range cpuP95 {
		cpuP95Map[key(s.Metric)] = s.ValueFloat
	}

	memReqMap := map[string]float64{}
	for _, s := range memReq {
		memReqMap[key(s.Metric)] = s.ValueFloat
	}

	cpuReqMap := map[string]float64{}
	for _, s := range cpuReq {
		cpuReqMap[key(s.Metric)] = s.ValueFloat
	}

	oomMap := map[string]bool{}
	for _, s := range oom {
		oomMap[key(s.Metric)] = s.ValueFloat >= 1
	}

	cpuThrottleMap := map[string]bool{}
	for _, s := range cpuThrottle {
		cpuThrottleMap[key(s.Metric)] = s.ValueFloat > 0
	}

	jvmHeapAfterGCMap := map[string]float64{}
	for _, s := range jvmHeapAfterGC {
		jvmHeapAfterGCMap[key(s.Metric)] = s.ValueFloat
	}

	jvmNonHeapMap := map[string]int64{}
	for _, s := range jvmNonHeap {
		jvmNonHeapMap[key(s.Metric)] = int64(s.ValueFloat)
	}

	// ---------------------------------------------------------------------
	// 3. Build results (service-level)
	// ---------------------------------------------------------------------

	results := make([]model.RightsizeResult, 0, len(memReqMap))

	for k, memReqBytesF := range memReqMap {
		parts := split3(k)
		ns, cl, container := parts[0], parts[1], parts[2]

		memReqBytes := int64(memReqBytesF)
		cpuReqCores := cpuReqMap[k]

		memRatio := memP95Map[k]
		cpuRatio := cpuP95Map[k]

		r := model.RightsizeResult{
			Namespace: ns,
			Cluster:   cl,
			Container: container,

			MemP95Ratio: memRatio,
			CpuP95Ratio: cpuRatio,

			MemRequestBytes: memReqBytes,
			CpuRequestCores: cpuReqCores,

			OOMKilled:    oomMap[k],
			CPUThrottled: cpuThrottleMap[k],

			JVMHeapAfterGCRatio: jvmHeapAfterGCMap[k],
			JVMNonHeapBytes:     jvmNonHeapMap[k],
		}

		// -----------------------------------------------------------------
		// 4. Recommendation math (NO decisions here)
		// -----------------------------------------------------------------

		if r.OOMKilled {
			r.MemRecommendedBytes = memReqBytes
			r.CpuRecommendedCores = cpuReqCores
		} else {
			r.MemRecommendedBytes = recommendMem(
				memReqBytes,
				memRatio,
				p.TargetUtil,
				p.SafetyFactor,
				p.MemRoundMiB,
			)

			r.CpuRecommendedCores = recommendCPU(
				cpuReqCores,
				cpuRatio,
				p.TargetUtil,
				p.SafetyFactor,
				p.CPURoundm,
			)
		}

		// -----------------------------------------------------------------
		// 5. Decisions (pure policy layer)
		// -----------------------------------------------------------------

		r.MemoryDecision, r.MemoryWhy = decision.DecideMemory(memRatio, r.OOMKilled)
		r.CPUDecision, r.CPUWhy = decision.DecideCPU(cpuRatio, r.CPUThrottled)
		
		r.JVMHeapDecision = decision.DecideJVMHeap(
			r.JVMHeapAfterGCRatio,
		)

		r.JVMNonHeapDecision = decision.DecideJVMNonHeap(
			r.JVMNonHeapBytes,
			r.MemRequestBytes,
		)

		results = append(results, r)
	}

	// ---------------------------------------------------------------------
	// 6. Rank results
	// ---------------------------------------------------------------------

	sort.Slice(results, func(i, j int) bool {
		if p.Bottom {
			return results[i].MemP95Ratio < results[j].MemP95Ratio
		}
		return results[i].MemP95Ratio > results[j].MemP95Ratio
	})

	if p.TopK > 0 && len(results) > p.TopK {
		results = results[:p.TopK]
	}

	return results, meta, nil
}

// -------------------------------------------------------------------------
// Helpers
// -------------------------------------------------------------------------

func (s *RightsizeService) query(
	ctx context.Context,
	expr string,
) ([]instantSample, error) {
	raw, err := s.vm.Query(ctx, vm.QueryOptions{Expr: expr})
	if err != nil {
		return nil, err
	}
	return parseInstantVector(raw)
}

func recommendMem(
	currentBytes int64,
	ratio, target, safety float64,
	roundMiB int64,
) int64 {
	if currentBytes <= 0 || ratio <= 0 || target <= 0 {
		return currentBytes
	}
	factor := (ratio / target) * safety
	reco := float64(currentBytes) * factor
	step := float64(roundMiB) * 1024 * 1024
	return int64(math.Ceil(reco/step) * step)
}

func recommendCPU(
	currentCores float64,
	ratio, target, safety float64,
	roundm int64,
) float64 {
	if currentCores <= 0 || ratio <= 0 || target <= 0 {
		return currentCores
	}
	factor := (ratio / target) * safety
	reco := currentCores * factor
	step := float64(roundm) / 1000.0
	return math.Ceil(reco/step) * step
}

func split3(s string) [3]string {
	out := [3]string{"", "", ""}
	cur := 0
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '|' && cur < 2 {
			out[cur] = s[start:i]
			cur++
			start = i + 1
		}
	}
	out[cur] = s[start:]
	return out
}
