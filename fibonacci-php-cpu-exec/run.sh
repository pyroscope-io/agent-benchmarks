#!/bin/bash
if [[ -z "${PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING}" ]]; then
    php fib.php
else
    export PYROSCOPE_APPLICATION_NAME=fibonacci-php-cpu-exec
    export PYROSCOPE_SERVER_ADDRESS=http://ingester:4040/
    pyroscope exec php fib.php
fi
