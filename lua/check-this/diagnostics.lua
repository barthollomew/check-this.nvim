local M = {}

local ns = vim.api.nvim_create_namespace("check-this")

local severity_map = {
  error = vim.diagnostic.severity.ERROR,
  warning = vim.diagnostic.severity.WARN,
  warn = vim.diagnostic.severity.WARN,
  info = vim.diagnostic.severity.INFO,
  hint = vim.diagnostic.severity.HINT,
}

local function resolve_severity(rule_id, severity, overrides)
  if overrides and overrides[rule_id] then
    return overrides[rule_id]
  end
  if not severity then
    return vim.diagnostic.severity.WARN
  end
  return severity_map[string.lower(severity)] or vim.diagnostic.severity.WARN
end

function M.publish(bufnr, output, opts)
  local diagnostics = {}
  for _, d in ipairs(output.diagnostics or {}) do
    local sev = resolve_severity(d.rule_id, d.severity, opts.severity)
    table.insert(diagnostics, {
      lnum = d.range.start.line,
      col = d.range.start.col,
      end_lnum = d.range["end"].line,
      end_col = d.range["end"].col,
      severity = sev,
      message = d.message,
      source = "check-this",
      user_data = {
        explanation = d.explanation,
        rule_id = d.rule_id,
        tags = d.tags,
      },
    })
  end
  vim.diagnostic.set(ns, bufnr, diagnostics, { virtual_text = true, underline = true })
end

function M.clear(bufnr)
  vim.diagnostic.reset(ns, bufnr)
end

function M.namespace()
  return ns
end

return M
