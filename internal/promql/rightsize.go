package promql

import "fmt"

// p95 over time of mem usage/request, aggregated at service-level (container)
func MemP95Ratio(namespace, cluster, window, subStep string) string {
	// Use subquery so quantile_over_time works over the evaluated ratio over time
	return fmt.Sprintf(`
quantile_over_time(
  0.95,
  (
    avg by (namespace, container, uw_cluster) (
      container_memory_working_set_bytes{namespace="%s",uw_cluster="%s"}
    )
    /
    avg by (namespace, container, uw_cluster) (
      kube_pod_container_resource_requests{namespace="%s",resource="memory",uw_cluster="%s"}
    )
  )[%s:%s]
)
`, namespace, cluster, namespace, cluster, window, subStep)
}

func CpuP95Ratio(namespace, cluster, window, subStep string) string {
	// CPU usage in cores = rate(cpu_seconds_total[5m])
	return fmt.Sprintf(`
quantile_over_time(
  0.95,
  (
    avg by (namespace, container, uw_cluster) (
      rate(container_cpu_usage_seconds_total{namespace="%s",uw_cluster="%s",container!="POD",container!=""}[5m])
    )
    /
    avg by (namespace, container, uw_cluster) (
      kube_pod_container_resource_requests{namespace="%s",resource="cpu",uw_cluster="%s"}
    )
  )[%s:%s]
)
`, namespace, cluster, namespace, cluster, window, subStep)
}

func MemRequests(namespace, cluster string) string {
	return fmt.Sprintf(`
avg by (namespace, container, uw_cluster) (
  kube_pod_container_resource_requests{namespace="%s",resource="memory",uw_cluster="%s"}
)
`, namespace, cluster)
}

func CpuRequests(namespace, cluster string) string {
	return fmt.Sprintf(`
avg by (namespace, container, uw_cluster) (
  kube_pod_container_resource_requests{namespace="%s",resource="cpu",uw_cluster="%s"}
)
`, namespace, cluster)
}

// OOMKilled guardrail (best-effort; depends on kube-state-metrics availability)
func OOMKilled(namespace, cluster, window string) string {
	return fmt.Sprintf(`
max by (namespace, container, uw_cluster) (
  max_over_time(
    kube_pod_container_status_last_terminated_reason{
      namespace="%s",
      uw_cluster="%s",
      reason="OOMKilled"
    }[%s]
  )
)
`, namespace, cluster, window)
}

func CPUThrottling(namespace, cluster, window string) string {
	return fmt.Sprintf(`
max by (namespace, container, uw_cluster) (
  increase(
    container_cpu_cfs_throttled_seconds_total{
      namespace="%s",
      uw_cluster="%s",
      container!="POD",
      container!=""
    }[%s]
  )
)
`, namespace, cluster, window)
}

func JVMHeapAfterGC(namespace, cluster string) string {
	return fmt.Sprintf(`
avg by (namespace, container, uw_cluster) (
  jvm_memory_usage_after_gc{
    namespace="%s",
    uw_cluster="%s",
    area="heap"
  }
)
`, namespace, cluster)
}

func JVMNonHeapBytes(namespace, cluster string) string {
	return fmt.Sprintf(`
avg by (namespace, container, uw_cluster) (
  jvm_memory_used_bytes{
    namespace="%s",
    uw_cluster="%s",
    area!="heap"
  }
)
`, namespace, cluster)
}
