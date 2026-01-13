package types

type BenchType string

const (
	BenchMicroservices BenchType = "microservices"
)

type ResourceType string

const (
	CPU    ResourceType = "cpu"
	Memory ResourceType = "memory"
)
