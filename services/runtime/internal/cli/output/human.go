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
	theme := NewTheme(w)
	printHeader(w, view, theme)
	for _, cv := range view.Cases {
		printCaseLine(w, cv, theme)
		printCaseDetail(w, cv, verbose, theme)
	}
	printFooter(w, view, theme)
}

func printHeader(w io.Writer, view RunView, theme Theme) {
	fmt.Fprintf(w, "%s  %s\n", theme.render(theme.bold, "Connor"), view.Version)
	fmt.Fprintf(w, "%s  %s\n", theme.render(theme.dim, "Target"), view.Target)
	if view.SuiteID != "" {
		fmt.Fprintf(w, "%s   %s (%d cases)\n\n",
			theme.render(theme.dim, "Suite"),
			view.SuiteID,
			len(view.Cases),
		)
		return
	}
	fmt.Fprintln(w)
}

func printCaseLine(w io.Writer, cv CaseView, theme Theme) {
	r := cv.Result
	icon := theme.render(theme.fail, "✗")
	if r.Passed {
		icon = theme.render(theme.pass, "✓")
	}

	gates := gatesSuffix(cv, r.Passed, r.Reason, theme)
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

func printCaseDetail(w io.Writer, cv CaseView, verbose bool, theme Theme) {
	r := cv.Result
	if r.Passed && !verbose {
		return
	}

	if !r.Passed {
		if cv.ExpectContains != "" {
			fmt.Fprintf(w, "   %s     expect_contains\n", theme.render(theme.label, "gate:"))
			fmt.Fprintf(w, "   %s %s\n", theme.render(theme.dim, "expected:"), cv.ExpectContains)
		}
		if cv.ExpectJSONSchema {
			fmt.Fprintf(w, "   %s     expect_json_schema\n", theme.render(theme.label, "gate:"))
		} else if cv.ExpectJSON {
			fmt.Fprintf(w, "   %s     expect_json\n", theme.render(theme.label, "gate:"))
		}
		if r.Reason != "" {
			fmt.Fprintf(w, "   %s   %s\n", theme.render(theme.label, "reason:"), r.Reason)
		}
		if body := r.Response.BodyPreview(200); body != "" {
			fmt.Fprintf(w, "   %s     %s\n", theme.render(theme.dim, "body:"), body)
		}
		if hint := failHint(r.Reason, r.Response, cv.ExpectContainsIgnoreCase); hint != "" {
			fmt.Fprintf(w, "   %s     %s\n", theme.render(theme.dim, "hint:"), hint)
		}
		return
	}

	fmt.Fprintf(w, "   %s %d\n", theme.render(theme.dim, "attempts:"), r.Response.Attempts)
	if body := r.Response.BodyPreview(200); body != "" {
		fmt.Fprintf(w, "   %s     %s\n", theme.render(theme.dim, "body:"), body)
	}
}

func printFooter(w io.Writer, view RunView, theme Theme) {
	passed, total, totalMs, slowest := summarize(view)

	fmt.Fprintln(w, theme.render(theme.dim, "────────────────────────────────────"))

	slowestLabel := ""
	if slowest.ID != "" && total > 0 {
		slowestLabel = fmt.Sprintf(" · slowest %dms (%s)", slowest.Result.Response.LatencyMs, slowest.ID)
	}

	summary := fmt.Sprintf("%d/%d passed%s · total %.1fs",
		passed, total, slowestLabel, float64(totalMs)/1000)
	fmt.Fprintln(w, theme.render(theme.dim, summary))

	if passed == total && total > 0 {
		fmt.Fprintln(w, theme.render(theme.pass, "GATE PASSED — safe to merge"))
		fmt.Fprintln(w, theme.render(theme.dim, "exit 0"))
		return
	}
	if total == 0 {
		fmt.Fprintln(w, theme.render(theme.fail, "GATE FAILED — no cases ran"))
		fmt.Fprintln(w, theme.render(theme.dim, "exit 1"))
		return
	}
	fmt.Fprintln(w, theme.render(theme.fail, "GATE FAILED — do not merge"))
	fmt.Fprintln(w, theme.render(theme.dim, "exit 1"))
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

func gatesSuffix(cv CaseView, passed bool, reason entities.FailReason, theme Theme) string {
	var parts []string
	if cv.ExpectContains != "" {
		parts = append(parts, themedGateLabel(theme, "contains", passed, reason, entities.FailReasonContentMismatch))
	}
	if cv.ExpectJSONSchema {
		parts = append(parts, themedSchemaGateLabel(theme, passed, reason))
	} else if cv.ExpectJSON {
		parts = append(parts, themedGateLabel(theme, "json", passed, reason, entities.FailReasonInvalidJSON))
	}
	if len(parts) == 0 {
		return ""
	}
	return "   " + strings.Join(parts, " ")
}

func themedGateLabel(theme Theme, name string, passed bool, reason, failReason entities.FailReason) string {
	label := gateLabel(name, passed, reason, failReason)
	return styleGateBadge(theme, label, passed, reason == failReason)
}

func themedSchemaGateLabel(theme Theme, passed bool, reason entities.FailReason) string {
	label := schemaGateLabel(passed, reason)
	failed := reason == entities.FailReasonSchemaMismatch || reason == entities.FailReasonInvalidJSON
	return styleGateBadge(theme, label, passed, failed)
}

func styleGateBadge(theme Theme, label string, passed, failed bool) string {
	if passed {
		return theme.render(theme.pass, label)
	}
	if failed {
		return theme.render(theme.fail, label)
	}
	return theme.render(theme.dim, label)
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

func schemaGateLabel(passed bool, reason entities.FailReason) string {
	if passed {
		return "schema ✓"
	}
	if reason == entities.FailReasonSchemaMismatch || reason == entities.FailReasonInvalidJSON {
		return "schema ✗"
	}
	return "schema"
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
