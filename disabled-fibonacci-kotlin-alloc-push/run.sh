#!/bin/bash
if [[ -z "${PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING}" ]]; then
    java -jar fib.jar
else
    export PYROSCOPE_APPLICATION_NAME=fibonacci-kotlin-alloc-push
    export PYROSCOPE_PROFILER_EVENT=alloc
    export PYROSCOPE_SERVER_ADDRESS=http://ingester:4040
    java -javaagent:pyroscope.jar -jar fib.jar
fi
