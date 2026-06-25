package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/benchmark"
	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/cli/output"
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
		verbose    bool
	)

	cmd := &cobra.Command{
		Use:   "run [suite.yaml]",
		Short: "Run a benchmark case or YAML suite against an OpenAI-compatible endpoint",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 1 {
				return runSuite(args[0], verbose)
			}
			return runSingleCase(model, prompt, expectJSON, timeoutMS, retries, verbose)
		},
	}

	cmd.Flags().StringVar(&model, "model", "", "Model id (single-case mode)")
	cmd.Flags().StringVar(&prompt, "prompt", "", "User prompt (single-case mode)")
	cmd.Flags().BoolVar(&expectJSON, "expect-json", false, "Require valid JSON in assistant content")
	cmd.Flags().Int64Var(&timeoutMS, "timeout-ms", 30000, "Per-attempt timeout in ms")
	cmd.Flags().IntVar(&retries, "retries", 2, "Max retries on transient errors")
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show body and attempts on passed cases")

	return cmd
}

func runSingleCase(model, prompt string, expectJSON bool, timeoutMS int64, retries int, verbose bool) error {
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
		entities.ExpectationsFromCase("", expectJSON),
		req,
		timeout,
		retry,
		client,
	)
	if err != nil {
		return err
	}

	view := output.RunView{
		Version: Version,
		Target:  client.Target(),
		Cases: []output.CaseView{{
			ID:         result.CaseID,
			Model:      model,
			ExpectJSON: expectJSON,
			Result:     result,
		}},
	}
	output.PrintRun(os.Stdout, view, verbose)
	return exitIfFailed(result.Passed)
}

func runSuite(path string, verbose bool) error {
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

	cases := make([]output.CaseView, len(spec.Cases))
	for i, c := range spec.Cases {
		cases[i] = output.CaseView{
			ID:             c.ID,
			Model:          c.Model,
			ExpectContains: c.ExpectContains,
			ExpectJSON:     c.ExpectJSON,
			Result:         suite.Results[i],
		}
	}

	view := output.RunView{
		Version: Version,
		Target:  client.Target(),
		SuiteID: suite.SuiteID,
		Cases:   cases,
	}
	output.PrintRun(os.Stdout, view, verbose)
	return exitIfFailed(suite.AllPassed())
}

func exitIfFailed(ok bool) error {
	if ok {
		return nil
	}
	os.Exit(1)
	return nil
}

func isYAMLPath(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".yaml") || strings.HasSuffix(lower, ".yml")
}
