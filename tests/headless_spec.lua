local M = {}

local function analyzer_path()
  local root = vim.fn.getcwd()
  local base = root .. "/analyzer/check-this"
  if vim.loop.os_uname().sysname == "Windows_NT" then
    local exe = base .. ".exe"
    if vim.loop.fs_stat(exe) then
      return exe
    end
  end
  return base
end

function M.run()
  local analyzer = analyzer_path()
  assert(vim.loop.fs_stat(analyzer), "analyzer binary not found at " .. analyzer)
  local check_this = require("check-this")
  check_this.setup({
    analyzer_path = analyzer,
    run_on_save = false,
    debounce_ms = 10,
  })

  local buf = vim.api.nvim_create_buf(false, true)
  vim.api.nvim_set_current_buf(buf)
  vim.api.nvim_buf_set_lines(buf, 0, -1, false, {
    "try:",
    "    risky()",
    "except Exception:",
    "    pass",
  })
  vim.bo[buf].filetype = "python"

  check_this.analyze(buf)

  local ns = require("check-this.diagnostics").namespace()
  local ok = vim.wait(5000, function()
    local diags = vim.diagnostic.get(buf, { namespace = ns })
    return #diags > 0
  end, 100)
  assert(ok, "expected diagnostics for swallowed exception")
end

return M
