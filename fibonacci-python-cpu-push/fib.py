"""Fibonacci CPU push benchmark for Python."""
import os
import pyroscope


def fib(n):
    """n-th Fibonacci number."""
    return n if n < 2 else fib(n-1) + fib(n-2)


if __name__ == '__main__':
    if os.getenv("PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING"):
        pyroscope.configure(app_name="fibonacci-python-cpu-push",
                            server_address="http://ingester:4040")
    fib(42)
