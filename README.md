## check-this.nvim

check-this.nvim is a Neovim plugin backed by a Go analyzer that surfaces operational-risk “bad defaults” (missing timeouts, unbounded retries, swallowed errors, risky globals) as inline diagnostics. It is advisory and fast, not a formal verifier.

### What it does / does not

- Does: parse buffers with Tree-sitter, apply pragmatic heuristics, show diagnostics via `vim.diagnostic`.
- Does not: prove correctness, eliminate all false positives, or replace a full LSP.

### Install

- **Plugin (lazy.nvim)**
  ```lua
  {
    "check-this/check-this.nvim",
    config = function()
      require("check-this").setup()
    end,
  }
  ```

- **Build the analyzer**
  ```sh
  cd analyzer
  # CGO is required because Tree-sitter ships C code. Use gcc/clang on macOS/Linux,
  # or MSYS2/MinGW on Windows.
  CGO_ENABLED=1 go build -o check-this ./cmd/check-this
  ```
  Add the resulting binary to your `PATH` (on Windows the file is `check-this.exe`) or set `analyzer_path` in the setup call.

### Quickstart

- In Neovim: open a Python or JS/TS buffer and run `:CheckThisAnalyze`. Diagnostics appear inline with messages and explanations.
- CLI smoke test:
  ```sh
  echo 'try:\n    risky()\nexcept Exception:\n    pass' | ./analyzer/check-this analyze --lang python
  ```

### Configuration

```lua
require("check-this").setup({
  analyzer_path = "check-this", -- or absolute path to the built binary
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

### Analyzer CLI

```
check-this analyze [--path <file>] [--lang <lang>] [--format json] [--config <path>]
```

Input is read from stdin; `--path` is used for display and language inference. `--lang` overrides detection.

### Supported rules (v1)

- **retry.unbounded**
  - Why: infinite retries amplify partial outages and overload dependencies.
  - Bad: `while True: requests.get(url)` with no sleep/backoff.
  - Better: add a max-attempts counter and `time.sleep`/backoff.
  - Suppress: `# check-this: disable=retry.unbounded`

- **net.no_timeout**
  - Why: network calls without timeouts can hang and block resources.
  - Bad: `requests.get(url)` or `fetch("/api")` without AbortController/timeout.
  - Better: `requests.get(url, timeout=2)` or `fetch(url, { signal = controller.signal })`.
  - Suppress: `# check-this: disable=net.no_timeout`

- **errors.swallowed**
  - Why: empty handlers hide failures and delay detection.
  - Bad: `except Exception: pass` or `catch (e) {}`.
  - Better: log, rethrow, or handle explicitly.
  - Suppress: `# check-this: disable=errors.swallowed`

- **state.global_mutable**
  - Why: globals create hidden coupling and contention.
  - Bad: module-level `my_cache = {}` mutated later.
  - Better: encapsulate in functions or use immutables.
  - Suppress: `# check-this: disable=state.global_mutable`

### Troubleshooting

- No diagnostics:
  - Ensure the analyzer binary is on `PATH` or configure `analyzer_path`.
  - Confirm the buffer filetype is supported (python, javascript, typescript).
  - Run `:CheckThisAnalyze` and check `:messages` for errors.
- Analyzer too slow:
  - Increase `debounce_ms` or disable rules you do not need.
- False positives:
  - Lower severity per rule or disable/suppress it inline with `check-this: disable=<rule>`.
- Windows specifics:
  - Build `check-this.exe` with `CGO_ENABLED=1`; ensure MSYS2/MinGW compiler is available.

### Design philosophy and limits

- Heuristic and reliability-focused: aims to catch risky defaults quickly.
- Explainable: each diagnostic includes a short message and explanation.
- Not perfect: false positives are expected; no deep dataflow or whole-repo analysis.
