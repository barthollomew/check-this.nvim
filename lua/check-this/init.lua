local config = require("check-this.config")
local diagnostics = require("check-this.diagnostics")
local runner = require("check-this.runner")

local M = {
  _opts = config.defaults,
}

local function current_buf()
  return vim.api.nvim_get_current_buf()
end

local function explain_under_cursor(bufnr)
  bufnr = bufnr or current_buf()
  local cursor = vim.api.nvim_win_get_cursor(0)
  local lnum = cursor[1] - 1
  local diags = vim.diagnostic.get(bufnr, { lnum = lnum, namespace = diagnostics.namespace() })
  if #diags == 0 then
    vim.notify("check-this: no diagnostic under cursor", vim.log.levels.INFO)
    return
  end
  local target = diags[1]
  local explanation = target.user_data and target.user_data.explanation or ""
  local rule_id = target.user_data and target.user_data.rule_id or "unknown"
  local lines = {
    string.format("[%s] %s", rule_id, target.message),
  }
  if explanation ~= "" then
    table.insert(lines, explanation)
  end
  vim.notify(table.concat(lines, "\n"), vim.log.levels.INFO)
end

local function setup_autocmds()
  local group = vim.api.nvim_create_augroup("CheckThisAuto", { clear = true })
  if not M._opts.run_on_save then
    return
  end
  vim.api.nvim_create_autocmd("BufWritePost", {
    group = group,
    callback = function(args)
      runner.run_debounced(args.buf, M._opts)
    end,
  })
end

function M.setup(user_opts)
  M._opts = config.merge(user_opts or {})
  setup_autocmds()
end

function M.analyze(bufnr)
  runner.run_immediate(bufnr or current_buf(), M._opts)
end

function M.analyze_all()
  for _, buf in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_loaded(buf) then
      runner.run_immediate(buf, M._opts)
    end
  end
end

function M.explain(bufnr)
  explain_under_cursor(bufnr)
end

return M
