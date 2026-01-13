package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/output"
	"github.com/BenjaminVolodarsky/cloud-monitoring-sentinel/internal/service"
	"github.com/spf13/cobra"
)

var (
	rsNamespace string
	rsCluster   string
	rsWindow    string
	rsFormat    string
	rsCSVOut    string
	rsHelmPatch string

	rsTargetUtil   float64
	rsSafetyFactor float64
	rsMemRoundMiB  int64
	rsCPURoundm    int64

	rsOOMWindow string
	rsSubStep   string
	rsTopK      int
	rsBottom    bool
)

var benchRightsizeCmd = &cobra.Command{
	Use:   "rightsize",
	Short: "Compute memory+CPU over/under-provisioning and recommend new requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		svc := service.NewRightsizeService()

		results, meta, err := svc.Run(ctx, service.RightsizeParams{
			Namespace:    rsNamespace,
			Cluster:      rsCluster,
			Window:       rsWindow,
			SubqueryStep: rsSubStep,

			OOMWindow: rsOOMWindow,

			TargetUtil:   rsTargetUtil,
			SafetyFactor: rsSafetyFactor,
			MemRoundMiB:  rsMemRoundMiB,
			CPURoundm:    rsCPURoundm,

			TopK:   rsTopK,
			Bottom: rsBottom,
		})
		if err != nil {
			return err
		}

		// ---------- STDOUT ----------
		switch rsFormat {
		case "table":
			output.RenderTable(results)

		case "json":
			return fmt.Errorf("json format not implemented yet")

		default:
			return fmt.Errorf("unknown format: %s", rsFormat)
		}

		// ---------- CSV ----------
		if rsCSVOut != "" {
			if err := output.WriteCSV(rsCSVOut, results, meta); err != nil {
				return fmt.Errorf("write csv: %w", err)
			}
			fmt.Fprintf(os.Stderr, "✓ wrote CSV to %s\n", rsCSVOut)
		}

		// ---------- HELM PATCH ----------
		if rsHelmPatch != "" {
			if err := output.WriteHelmValuesPatch(rsHelmPatch, results); err != nil {
				return fmt.Errorf("write helm patch: %w", err)
			}
			fmt.Fprintf(os.Stderr, "✓ wrote Helm patch to %s\n", rsHelmPatch)
		}

		return nil
	},
}

func init() {
	benchCmd.AddCommand(benchRightsizeCmd)

	benchRightsizeCmd.Flags().StringVar(&rsNamespace, "namespace", "microservices", "Kubernetes namespace")
	benchRightsizeCmd.Flags().StringVar(&rsCluster, "cluster", "", "Cluster label (uw_cluster)")
	benchRightsizeCmd.Flags().StringVar(&rsWindow, "window", "24h", "Time window (e.g. 24h, 7d)")
	benchRightsizeCmd.Flags().StringVar(&rsSubStep, "sub-step", "5m", "Subquery step (e.g. 1m, 5m, 15m)")

	benchRightsizeCmd.Flags().StringVar(&rsOOMWindow, "oom-window", "14d", "Lookback window to detect OOMKilled")

	benchRightsizeCmd.Flags().StringVar(&rsFormat, "format", "table", "Output format: table|json")
	benchRightsizeCmd.Flags().StringVar(&rsCSVOut, "csv", "", "Write CSV to path (optional)")
	benchRightsizeCmd.Flags().StringVar(&rsHelmPatch, "helm-patch", "", "Write Helm values patch snippet (optional)")

	benchRightsizeCmd.Flags().Float64Var(&rsTargetUtil, "target-util", 0.70, "Target p95 usage/request ratio (e.g. 0.7)")
	benchRightsizeCmd.Flags().Float64Var(&rsSafetyFactor, "safety", 1.15, "Safety multiplier for recommendation (e.g. 1.15)")

	benchRightsizeCmd.Flags().Int64Var(&rsMemRoundMiB, "mem-round-mib", 64, "Round memory recommendation up to this MiB multiple")
	benchRightsizeCmd.Flags().Int64Var(&rsCPURoundm, "cpu-round-m", 10, "Round CPU recommendation up to this millicore multiple")

	benchRightsizeCmd.Flags().IntVar(&rsTopK, "topk", 50, "Limit results to top K (after ranking)")
	benchRightsizeCmd.Flags().BoolVar(&rsBottom, "bottom", true, "Rank by most overprovisioned (lowest ratios). Use --bottom=false for most underprovisioned.")

	_ = benchRightsizeCmd.MarkFlagRequired("cluster")
}
