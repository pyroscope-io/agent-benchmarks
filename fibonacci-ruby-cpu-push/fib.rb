require "pyroscope"

def fib(n)
  if n < 2 then
    n
  else
    fib(n-1) + fib(n-2)
  end
end

if ENV["PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING"] then
  Pyroscope.configure do |config|
    config.app_name = "fibonacci-ruby-cpu-push"
    config.server_address = "http://ingester:4040"
  end
end

fib(43)
