## check-this.nvim

Most production outages come from boring mistakes that linters ignore: an HTTP call with no timeout, a `while True` retry loop with no backoff, or an `except Exception: pass` block that hides a crash.

check-this.nvim is a fast Go analyzer and Neovim plugin. It catches these risky patterns instantly as you code.

### Sample Output

When you run `check-this analyze --lang python` on risky code, it outputs clean structured JSON:

```json
{
  "version": "1.0",
  "path": "stdin",
  "language": "python",
  "diagnostics": [
    {
      "rule_id": "errors.swallowed",
      "severity": "warning",
      "message": "Exception handled but nothing done",
      "explanation": "Swallowing exceptions makes outages harder to detect; log or re-raise instead.",
      "range": { "start": { "line": 2, "col": 0 }, "end": { "line": 3, "col": 8 } },
      "tags": ["reliability", "errors"]
    }
  ],
  "stats": { "parse_ms": 0, "analyze_ms": 0, "rules_run": 4 }
}
```

### Architecture

The plugin uses two main parts working together:

1. Asynchronous Neovim Lua Runner
   - Written in Lua using `vim.system` to run the analyzer process asynchronously.
   - Pipes buffer contents via stdin so unsaved changes are checked instantly without disk I/O.
   - Parses the JSON output from the analyzer and updates native Neovim diagnostics in a dedicated `check-this` namespace using `vim.diagnostic`.
   - Debounces checks on save to keep the editor responsive.

2. CGO Tree-sitter Go AST Analyzer
   - Built in Go using CGO and Tree-sitter for fast parsing of Python and JS/TS code.
   - Evaluates specific heuristic rules:
     - `retry.unbounded`: detects loops that retry without limits or backoff.
     - `net.no_timeout`: detects network calls missing timeout options.
     - `errors.swallowed`: detects empty catch or except blocks.
     - `state.global_mutable`: detects module level mutable shared state.

### Testing

Testing is split into two layers:

1. Analyzer Unit Tests
   - Go unit tests covering AST rule matching, suppression comments, config parsing, and language detection. Run with `go test -v ./...`.

2. Headless Neovim Plugin Tests
   - Integration tests running inside Neovim in headless mode via `nvim --headless -u tests/minimal_init.lua`. These tests verify end-to-end execution, stdin piping, JSON decoding, and diagnostic population.

### Quick install

- Plugin (lazy.nvim)
  ```lua
  {
    "barthollomew/check-this.nvim",
    config = function()
      require("check-this").setup()
    end,
  }
  ```

- Analyzer binary
  ```sh
  cd analyzer
  # Tree-sitter requires CGO.
  # macOS/Linux: gcc or clang
  # Windows: MSYS2 / MinGW
  CGO_ENABLED=1 go build -o check-this ./cmd/check-this
  ```

  Put the binary on your `PATH` (Windows: `check-this.exe`), or point to it with `analyzer_path`.

### Quick usage

- In Neovim: open a Python or JS/TS file and run `:CheckThisAnalyze`. Diagnostics show inline with a short explanation.
- CLI smoke test:
  ```sh
  printf "try:\n    risky()\nexcept Exception:\n    pass\n" | ./analyzer/check-this analyze --lang python
  ```

### What it does

- Runs asynchronously using Tree-sitter.
- Works on unsaved buffers via stdin.
- Current rules:
  - `retry.unbounded`
  - `net.no_timeout`
  - `errors.swallowed`
  - `state.global_mutable`
- Debounced on save, with a manual command when you want it.
- Stable JSON output for scripting and tests.

### Configuration (minimal)

```lua
require("check-this").setup({
  analyzer_path = "check-this",
  debounce_ms = 500,
  run_on_save = true,
  severity = {
    ["retry.unbounded"] = vim.diagnostic.severity.WARN,
    ["net.no_timeout"] = vim.diagnostic.severity.INFO,
  },
  rules = {
    ["state.global_mutable"] = { enabled = false },
  },
})
```

### More details

See `doc/check-this.txt` for CLI flags, JSON schemas, rule heuristics, suppression comments, and troubleshooting notes.
