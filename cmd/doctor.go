package cmd

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/buildinfo"
	"github.com/spf13/cobra"
)

const vmBaseURL = "http://vmselect.management.prod.internal:8481"

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check connectivity and environment prerequisites",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("upctl %s (commit=%s, built=%s)\n", buildinfo.Version, buildinfo.Commit, buildinfo.Date)

		ok := true

		// 1) DNS
		host := "vmselect.management.prod.internal"
		if _, err := net.LookupHost(host); err != nil {
			ok = false
			fmt.Fprintf(os.Stderr, "✗ DNS lookup failed for %s: %v\n", host, err)
		} else {
			fmt.Printf("✓ DNS: %s resolves\n", host)
		}

		// 2) HTTP health
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, vmBaseURL+"/-/healthy", nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			ok = false
			fmt.Fprintf(os.Stderr, "✗ VictoriaMetrics health check failed: %v\n", err)
		} else {
			defer resp.Body.Close()
			if resp.StatusCode/100 != 2 {
				ok = false
				fmt.Fprintf(os.Stderr, "✗ VictoriaMetrics unhealthy: %s\n", resp.Status)
			} else {
				fmt.Printf("✓ VictoriaMetrics: reachable (%s)\n", resp.Status)
			}
		}

		// 3) Optional: quick query endpoint sanity (super useful)
		// We hit the Prometheus API endpoint via vmselect /select/.../query
		// (If you prefer, you can reuse your vm client here)
		ctx2, cancel2 := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel2()

		url := vmBaseURL + "/select/0/prometheus/api/v1/query?query=1"
		req2, _ := http.NewRequestWithContext(ctx2, http.MethodGet, url, nil)
		resp2, err := http.DefaultClient.Do(req2)
		if err != nil {
			ok = false
			fmt.Fprintf(os.Stderr, "✗ Query endpoint failed: %v\n", err)
		} else {
			defer resp2.Body.Close()
			if resp2.StatusCode/100 != 2 {
				ok = false
				fmt.Fprintf(os.Stderr, "✗ Query endpoint error: %s\n", resp2.Status)
			} else {
				fmt.Printf("✓ Prometheus API query: ok\n")
			}
		}

		if !ok {
			return fmt.Errorf("doctor found issues")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}
