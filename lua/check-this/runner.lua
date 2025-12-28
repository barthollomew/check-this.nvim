local diagnostics = require("check-this.diagnostics")

local uv = vim.uv or vim.loop

local M = {}
local timers = {}

local function buf_text(bufnr)
  if not vim.api.nvim_buf_is_valid(bufnr) then
    return nil
  end
  local lines = vim.api.nvim_buf_get_lines(bufnr, 0, -1, false)
  return table.concat(lines, "\n")
end

local function should_analyze(bufnr, opts)
  local ft = vim.bo[bufnr].filetype
  if not ft or ft == "" then
    return false
  end
  if not opts.filetypes or #opts.filetypes == 0 then
    return true
  end
  for _, allowed in ipairs(opts.filetypes) do
    if allowed == ft then
      return true
    end
  end
  return false
end

local function decode_json(payload)
  local ok, decoded = pcall(vim.json.decode, payload)
  if not ok then
    return nil, decoded
  end
  return decoded, nil
end

local function run_system(cmd, stdin, cb)
  vim.system(cmd, { stdin = stdin }, function(res)
    if res.code ~= 0 then
      cb(nil, string.format("analyzer exited with code %s: %s", res.code, res.stderr))
      return
    end
    local decoded, err = decode_json(res.stdout)
    if err then
      cb(nil, string.format("failed to decode analyzer output: %s", err))
      return
    end
    cb(decoded, nil)
  end)
end

local function run_jobstart(cmd, stdin, cb)
  local stdout, stderr = {}, {}
  local handle = vim.fn.jobstart(cmd, {
    stdin = "pipe",
    on_stdout = function(_, data, _)
      if data then
        table.insert(stdout, table.concat(data, "\n"))
      end
    end,
    on_stderr = function(_, data, _)
      if data then
        table.insert(stderr, table.concat(data, "\n"))
      end
    end,
    on_exit = function(_, code, _)
      if code ~= 0 then
        cb(nil, string.format("analyzer exited with code %s: %s", code, table.concat(stderr, "\n")))
        return
      end
      local decoded, err = decode_json(table.concat(stdout, "\n"))
      if err then
        cb(nil, string.format("failed to decode analyzer output: %s", err))
        return
      end
      cb(decoded, nil)
    end,
  })

  if handle > 0 then
    vim.fn.chansend(handle, stdin)
    vim.fn.chanclose(handle, "stdin")
  else
    cb(nil, "failed to start analyzer process")
  end
end

local function build_cmd(bufnr, opts)
  local path = vim.api.nvim_buf_get_name(bufnr)
  local lang = vim.bo[bufnr].filetype
  return {
    opts.analyzer_path or "check-this",
    "analyze",
    "--path",
    path,
    "--lang",
    lang,
    "--format",
    "json",
  }
end

local function run_once(bufnr, opts)
  if not should_analyze(bufnr, opts) then
    diagnostics.clear(bufnr)
    return
  end
  local text = buf_text(bufnr)
  if not text then
    return
  end
  local cmd = build_cmd(bufnr, opts)
  local runner = vim.system and run_system or run_jobstart
  runner(cmd, text, function(output, err)
    vim.schedule(function()
      if err then
        diagnostics.publish(bufnr, {
          diagnostics = {
            {
              rule_id = "internal.analyzer_error",
              severity = "error",
              message = "check-this.nvim: analyzer error (see :messages)",
              range = { start = { line = 0, col = 0 }, ["end"] = { line = 0, col = 1 } },
            },
          },
        }, opts)
        vim.notify(err, vim.log.levels.ERROR)
        return
      end
      diagnostics.publish(bufnr, output, opts)
    end)
  end)
end

function M.run_debounced(bufnr, opts)
  local existing = timers[bufnr]
  if existing then
    existing:stop()
    existing:close()
  end
  local timer = uv.new_timer()
  timers[bufnr] = timer
  timer:start(opts.debounce_ms or 500, 0, function()
    timer:stop()
    timer:close()
    timers[bufnr] = nil
    run_once(bufnr, opts)
  end)
end

function M.run_immediate(bufnr, opts)
  run_once(bufnr, opts)
end

return M
