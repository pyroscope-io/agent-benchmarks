#!/bin/bash
if [[ -z "${PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING}" ]]; then
    /opt/agent/fib
else
    export PYROSCOPE_APPLICATION_NAME=fibonacci-dotnet-cpu-exec
    export PYROSCOPE_SERVER_ADDRESS=http://ingester:4040/
    export PYROSCOPE_SPY_NAME=dotnetspy

    pyroscope exec /opt/agent/fib
fi
