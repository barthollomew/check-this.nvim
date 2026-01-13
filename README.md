## check-this.nvim

Small, opinionated risk checks for Neovim.

check-this.nvim runs a tiny Go analyzer over your buffer and flags things that might bite later. Missing timeouts, infinite retries, swallowed errors, sketchy globals. It is fast, readable, and deliberately heuristic. This is not formal verification

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

See `doc/check-this.txt` for:
- CLI flags and JSON schema
- Rule heuristics and suppression comments
- Platform notes and failure modes
- Design goals and known limitations
