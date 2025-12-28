## check-this.nvim

Lightweight risk checks for Neovim. check-this.nvim runs a small Go analyzer that spots “bad defaults” (missing timeouts, unbounded retries, swallowed errors, risky globals) and shows inline diagnostics. Fast, explainable, and intentionally heuristic.

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
  # Tree-sitter needs CGO; use gcc/clang on macOS/Linux, MSYS2/MinGW on Windows.
  CGO_ENABLED=1 go build -o check-this ./cmd/check-this
  ```
  Put the binary on `PATH` (Windows: `check-this.exe`) or set `analyzer_path` in setup.

### Quick usage

- In Neovim: open a Python or JS/TS file and run `:CheckThisAnalyze`. Diagnostics appear inline with explanations.
- CLI smoke test:
  ```sh
  printf "try:\n    risky()\nexcept Exception:\n    pass\n" | ./analyzer/check-this analyze --lang python
  ```

### What you get

- Async analysis driven by Tree-sitter; works on unsaved buffers (stdin).
- Rules: `retry.unbounded`, `net.no_timeout`, `errors.swallowed`, `state.global_mutable`.
- Debounced run on save; manual command available.
- JSON output contract for scripting and tests.

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

### Need more detail?

See `doc/check-this.txt` for deeper documentation:
- Full CLI flags and JSON schema
- Rule heuristics and suppression directives
- Failure handling, platform notes, and troubleshooting
- Design philosophy and limitations
