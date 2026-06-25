package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/benchmark"
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
		Use:   "run [suite.yaml]",
		Short: "Run a benchmark case or YAML suite against an OpenAI-compatible endpoint",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return runSuite(args[0])
			}
			return runSingleCase(model, prompt, expectJSON, timeoutMS, retries)
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Model id (single-case mode)")
	cmd.Flags().StringVar(&prompt, "prompt", "", "User prompt (single-case mode)")
	cmd.Flags().BoolVar(&expectJSON, "expect-json", false, "Require valid JSON in assistant content")
	cmd.Flags().Int64Var(&timeoutMS, "timeout-ms", 30000, "Per-attempt timeout in ms")
	cmd.Flags().IntVar(&retries, "retries", 2, "Max retries on transient errors")

	return cmd
}

func runSingleCase(model, prompt string, expectJSON bool, timeoutMS int64, retries int) error {
	if model == "" || prompt == "" {
		return fmt.Errorf("connor: --model and --prompt are required (or pass a .yaml suite file)")
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
	retry, err := reliability.NewRetryPolicy(retries, reliability.DefaultBackoffBase)
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

	printCaseResult(result)
	return exitIfFailed(result.Passed, 1, 1)
}

func runSuite(path string) error {
	if !isYAMLPath(path) {
		return fmt.Errorf("connor: suite file must end with .yaml or .yml")
	}

	spec, err := benchmark.ParseFile(path)
	if err != nil {
		return err
	}

	client, err := openai_compatible.NewClientFromEnv(nil)
	if err != nil {
		return err
	}

	suite, err := application.ExecuteSuite(context.Background(), spec, client)
	if err != nil {
		return err
	}

	fmt.Printf("suite:     %s\n", suite.SuiteID)
	fmt.Printf("cases:     %d\n\n", len(suite.Results))

	for _, result := range suite.Results {
		printCaseResult(result)
		fmt.Println()
	}

	total := len(suite.Results)
	passed := suite.PassedCount()
	fmt.Printf("summary:   %d/%d passed\n", passed, total)

	return exitIfFailed(suite.AllPassed(), passed, total)
}

func printCaseResult(result entities.CaseResult) {
	fmt.Printf("case_id:   %s\n", result.CaseID)
	fmt.Printf("passed:    %v\n", result.Passed)
	fmt.Printf("reason:    %s\n", result.Reason)
	fmt.Printf("status:    %d\n", result.Response.HTTPStatus)
	fmt.Printf("latency:   %dms\n", result.Response.LatencyMs)
	fmt.Printf("attempts:  %d\n", result.Response.Attempts)
	fmt.Printf("body:      %s\n", result.Response.BodyPreview(200))
}

func exitIfFailed(ok bool, _, _ int) error {
	if ok {
		fmt.Printf("exit:      0\n")
		return nil
	}
	fmt.Printf("exit:      1\n")
	os.Exit(1)
	return nil
}

func isYAMLPath(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}
