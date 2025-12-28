local M = {}

M.defaults = {
  analyzer_path = "check-this",
  debounce_ms = 500,
  run_on_save = true,
  severity = {},
  rules = {},
  filetypes = { "python", "javascript", "typescript" },
}

function M.merge(user_opts)
  return vim.tbl_deep_extend("force", M.defaults, user_opts or {})
end

return M
