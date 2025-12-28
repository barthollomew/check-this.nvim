if vim.g.loaded_check_this then
  return
end
vim.g.loaded_check_this = true

local check_this = require("check-this")

vim.api.nvim_create_user_command("CheckThisAnalyze", function(opts)
  local bufnr = opts.buf or vim.api.nvim_get_current_buf()
  check_this.analyze(bufnr)
end, { desc = "Run check-this analyzer for current buffer" })

vim.api.nvim_create_user_command("CheckThisAnalyzeAll", function()
  check_this.analyze_all()
end, { desc = "Run check-this analyzer for all loaded buffers" })

vim.api.nvim_create_user_command("CheckThisExplain", function()
  check_this.explain()
end, { desc = "Show explanation for the diagnostic under cursor" })
