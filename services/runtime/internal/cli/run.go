package cli

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/application"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/reliability"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/infrastructure/providers/openai_compatible"
)

func newRunCmd() *cobra.Command {
	var (
		model      string
		prompt     string
		expectJSON bool
		timeoutMS  int64
		retries    int
	)

	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a single benchmark case against an OpenAI-compatible endpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			if model == "" || prompt == "" {
				return fmt.Errorf("connor: --model and --prompt are required")
			}

			client, err := openai_compatible.NewClientFromEnv(nil)
			if err != nil {
				return err
			}

			req, err := entities.NewRequest(model, []entities.Message{
				{Role: "user", Content: prompt},
			}, false)
			if err != nil {
				return err
			}

			timeout, err := reliability.NewTimeoutPolicyFromMS(timeoutMS)
			if err != nil {
				return err
			}
			retry, err := reliability.NewRetryPolicy(retries, 200*time.Millisecond)
			if err != nil {
				return err
			}

			result, err := application.EvaluateCase(
				context.Background(),
				"cli-run",
				expectJSON,
				req,
				timeout,
				retry,
				client,
			)
			if err != nil {
				return err
			}

			fmt.Printf("case_id:   %s\n", result.CaseID)
			fmt.Printf("passed:    %v\n", result.Passed)
			fmt.Printf("reason:    %s\n", result.Reason)
			fmt.Printf("status:    %d\n", result.Response.HTTPStatus)
			fmt.Printf("latency:   %dms\n", result.Response.LatencyMs)
			fmt.Printf("attempts:  %d\n", result.Response.Attempts)
			fmt.Printf("body:      %s\n", result.Response.BodyPreview(200))

			if !result.Passed {
				os.Exit(1)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Model id (required)")
	cmd.Flags().StringVar(&prompt, "prompt", "", "User prompt (required)")
	cmd.Flags().BoolVar(&expectJSON, "expect-json", false, "Require valid JSON in assistant content")
	cmd.Flags().Int64Var(&timeoutMS, "timeout-ms", 30000, "Per-attempt timeout in ms")
	cmd.Flags().IntVar(&retries, "retries", 2, "Max retries on transient errors")
	_ = cmd.MarkFlagRequired("model")
	_ = cmd.MarkFlagRequired("prompt")

	return cmd
}
