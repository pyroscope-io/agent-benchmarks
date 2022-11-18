#!/bin/bash
# set -ex
if [[ -z "${PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING}" ]]; then
    /opt/agent/fib
else
    export PYROSCOPE_APPLICATION_NAME=fibonacci-dotnet-cpu-exec
    export PYROSCOPE_SERVER_ADDRESS=http://ingester:4040/
    export PYROSCOPE_SPY_NAME=dotnetspy

    export CORECLR_ENABLE_PROFILING=1
    export CORECLR_PROFILER={BD1A650D-AC5D-4896-B64F-D6FA25D6B26A}
    export CORECLR_PROFILER_PATH=/opt/agent/Datadog.Profiler.Native.so
    export LD_PRELOAD=/opt/agent/Datadog.Linux.ApiWrapper.x64.so
    export PROFILING_ENABLED=1
    export PROFILING_CPU_ENABLED=true
    export PROFILING_WALLTIME_ENABLED=false
    
    



    exec /opt/agent/fib
fi
