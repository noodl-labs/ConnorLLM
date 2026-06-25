package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/noodl-labs/ConnorLLM/services/runtime/internal/runtime/domain/entities"
)

const modelColWidth = 34

// PrintRun renders a human-readable run report to w.
func PrintRun(w io.Writer, view RunView, verbose bool) {
	printHeader(w, view)
	for _, cv := range view.Cases {
		printCaseLine(w, cv)
		printCaseDetail(w, cv, verbose)
	}
	printFooter(w, view)
}

func printHeader(w io.Writer, view RunView) {
	fmt.Fprintf(w, "Connor  %s\n", view.Version)
	fmt.Fprintf(w, "Target  %s\n", view.Target)
	if view.SuiteID != "" {
		fmt.Fprintf(w, "Suite   %s (%d cases)\n\n", view.SuiteID, len(view.Cases))
		return
	}
	fmt.Fprintln(w)
}

func printCaseLine(w io.Writer, cv CaseView) {
	r := cv.Result
	icon := "✗"
	if r.Passed {
		icon = "✓"
	}

	gates := gatesSuffix(cv, r.Passed, r.Reason)
	fmt.Fprintf(w, "%s  %-12s  %-*s  %4dms  HTTP %d%s\n",
		icon,
		cv.ID,
		modelColWidth,
		truncate(cv.Model, modelColWidth),
		r.Response.LatencyMs,
		r.Response.HTTPStatus,
		gates,
	)
}

func printCaseDetail(w io.Writer, cv CaseView, verbose bool) {
	r := cv.Result
	if r.Passed && !verbose {
		return
	}

	if !r.Passed {
		if cv.ExpectContains != "" {
			fmt.Fprintln(w, "   gate:     expect_contains")
			fmt.Fprintf(w, "   expected: %s\n", cv.ExpectContains)
		}
		if cv.ExpectJSON {
			fmt.Fprintln(w, "   gate:     expect_json")
		}
		if r.Reason != "" {
			fmt.Fprintf(w, "   reason:   %s\n", r.Reason)
		}
		if body := r.Response.BodyPreview(200); body != "" {
			fmt.Fprintf(w, "   body:     %s\n", body)
		}
		if hint := failHint(r.Reason, r.Response); hint != "" {
			fmt.Fprintf(w, "   hint:     %s\n", hint)
		}
		return
	}

	fmt.Fprintf(w, "   attempts: %d\n", r.Response.Attempts)
	if body := r.Response.BodyPreview(200); body != "" {
		fmt.Fprintf(w, "   body:     %s\n", body)
	}
}

func printFooter(w io.Writer, view RunView) {
	passed, total, totalMs, slowest := summarize(view)

	fmt.Fprintln(w, "────────────────────────────────────")

	slowestLabel := ""
	if slowest.ID != "" && total > 0 {
		slowestLabel = fmt.Sprintf(" · slowest %dms (%s)", slowest.Result.Response.LatencyMs, slowest.ID)
	}

	fmt.Fprintf(w, "%d/%d passed%s · total %.1fs\n",
		passed, total, slowestLabel, float64(totalMs)/1000)

	if passed == total && total > 0 {
		fmt.Fprintln(w, "GATE PASSED — safe to merge")
		fmt.Fprintln(w, "exit 0")
		return
	}
	if total == 0 {
		fmt.Fprintln(w, "GATE FAILED — no cases ran")
		fmt.Fprintln(w, "exit 1")
		return
	}
	fmt.Fprintln(w, "GATE FAILED — do not merge")
	fmt.Fprintln(w, "exit 1")
}

func summarize(view RunView) (passed, total int, totalMs int64, slowest CaseView) {
	total = len(view.Cases)
	for _, cv := range view.Cases {
		if cv.Result.Passed {
			passed++
		}
		totalMs += cv.Result.Response.LatencyMs
		if slowest.ID == "" || cv.Result.Response.LatencyMs > slowest.Result.Response.LatencyMs {
			slowest = cv
		}
	}
	return passed, total, totalMs, slowest
}

func gatesSuffix(cv CaseView, passed bool, reason entities.FailReason) string {
	var parts []string
	if cv.ExpectContains != "" {
		parts = append(parts, gateLabel("contains", passed, reason, entities.FailReasonContentMismatch))
	}
	if cv.ExpectJSON {
		parts = append(parts, gateLabel("json", passed, reason, entities.FailReasonInvalidJSON))
	}
	if len(parts) == 0 {
		return ""
	}
	return "   " + strings.Join(parts, " ")
}

func gateLabel(name string, passed bool, reason, failReason entities.FailReason) string {
	if passed {
		return name + " ✓"
	}
	if reason == failReason {
		return name + " ✗"
	}
	return name
}

func truncate(s string, max int) string {
	if max <= 0 || len(s) <= max {
		return s
	}
	if max <= 1 {
		return s[:max]
	}
	return s[:max-1] + "…"
}
