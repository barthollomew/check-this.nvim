-- Minimal init for headless testing.
local root = vim.fn.getcwd()
vim.opt.runtimepath:append(root)
vim.opt.runtimepath:append(root .. "/tests")
package.path = root .. "/tests/?.lua;" .. package.path
vim.g.loaded_python_provider = 0
vim.g.loaded_python3_provider = 0
