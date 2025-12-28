package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/barthollomew/check-this.nvim/analyzer/internal/config"
	"github.com/barthollomew/check-this.nvim/analyzer/internal/engine"
	"github.com/barthollomew/check-this.nvim/analyzer/internal/ts"
)

// Run executes the CLI and writes JSON output to stdout. It returns the exit
// code that should be used by main.
func Run(args []string, stdin []byte) (int, error) {
	if len(args) == 0 {
		return usageError("missing command")
	}
	cmd := args[0]
	switch cmd {
	case "analyze":
		opts, err := parseAnalyzeFlags(args[1:])
		if err != nil {
			return usageError(err.Error())
		}
		lang := opts.Lang
		if lang == "" {
			lang = ts.DetectLanguage("", opts.Path)
		}
		request := engine.AnalyzeInput{
			Path:    opts.Path,
			Lang:    lang,
			Source:  stdin,
			Config:  config.Config{},
			Version: "1.0",
		}
		if err := engine.ValidateInput(request); err != nil {
			return 2, err
		}
		out, err := engine.NewEngine().Analyze(request)
		if err != nil {
			return 4, err
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			return 4, err
		}
		return 0, nil
	default:
		return usageError(fmt.Sprintf("unknown command %q", cmd))
	}
}

type analyzeFlags struct {
	Path       string
	Lang       string
	Format     string
	ConfigPath string
}

func parseAnalyzeFlags(args []string) (analyzeFlags, error) {
	fs := flag.NewFlagSet("analyze", flag.ContinueOnError)
	var opts analyzeFlags
	fs.StringVar(&opts.Path, "path", "", "optional path for diagnostics")
	fs.StringVar(&opts.Lang, "lang", "", "language override (python, javascript, typescript)")
	fs.StringVar(&opts.Format, "format", "json", "output format (json)")
	fs.StringVar(&opts.ConfigPath, "config", "", "path to config file")
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		return opts, err
	}
	if opts.Format != "json" {
		return opts, fmt.Errorf("unsupported format %s", opts.Format)
	}
	return opts, nil
}

func usageError(msg string) (int, error) {
	return 2, fmt.Errorf("usage error: %s", strings.TrimSpace(msg))
}
