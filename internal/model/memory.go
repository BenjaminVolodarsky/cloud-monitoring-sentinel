package model

type MemoryBenchResult struct {
	Namespace string
	Cluster   string
	Container string

	P95UsageRatio float64

	RequestBytes     int64
	RecommendedBytes int64
	Decision         string // REDUCE | KEEP | INCREASE
}
