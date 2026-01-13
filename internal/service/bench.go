package service

import "github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/types"

type BenchParams struct {
	Type       types.BenchType
	TimeWindow int
	Resource   types.ResourceType
}

func Benchmark(params BenchParams) {

}
