package model

type Decision string

const (
	DecisionReduce   Decision = "REDUCE"
	DecisionKeep     Decision = "KEEP"
	DecisionIncrease Decision = "INCREASE"
	DecisionSkipOOM  Decision = "SKIP_OOM"
)

type RightsizeMeta struct {
	Namespace    string  `json:"namespace"`
	Cluster      string  `json:"cluster"`
	Window       string  `json:"window"`
	OOMWindow    string  `json:"oom_window"`
	TargetUtil   float64 `json:"target_util"`
	SafetyFactor float64 `json:"safety_factor"`
	SubqueryStep string  `json:"subquery_step"`
}

type CPUDecision string
type MemoryDecision string
type JVMDecision string

const (
	CPUReduce         CPUDecision = "REDUCE"
	CPUKeep           CPUDecision = "KEEP"
	CPUIncrease       CPUDecision = "INCREASE"
	CPUSkipThrottling CPUDecision = "SKIP_THROTTLING"
)

const (
	MemReduce   MemoryDecision = "REDUCE"
	MemKeep     MemoryDecision = "KEEP"
	MemIncrease MemoryDecision = "INCREASE"
	MemSkipOOM  MemoryDecision = "SKIP_OOM"
)

const (
	JVMKeep     JVMDecision = "KEEP"
	JVMIncrease JVMDecision = "INCREASE"
)

type RightsizeResult struct {
	Namespace string
	Cluster   string
	Container string

	MemP95Ratio float64
	CpuP95Ratio float64

	MemRequestBytes int64
	CpuRequestCores float64

	MemRecommendedBytes int64
	CpuRecommendedCores float64

	OOMKilled    bool
	CPUThrottled bool

	JVMHeapAfterGCRatio float64
	JVMNonHeapBytes     int64

	MemoryDecision     MemoryDecision
	CPUDecision        CPUDecision
	JVMHeapDecision    JVMDecision
	JVMNonHeapDecision JVMDecision
	// Deltas (recommended - current)
	CpuDeltaCores float64
	MemDeltaBytes int64

	// Optional dollar estimate (if pricing provided)
	EstSavingsPerHourUSD float64

	// Explanations ("why")
	CPUWhy    string
	MemoryWhy string
	JVMWhy    string
}
