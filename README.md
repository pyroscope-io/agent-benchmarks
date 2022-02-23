# Pyroscope Agent Benchmark Suite

The benchmark suite is designed around three components:
- A ingester that represents a pyroscope server.
  It can support built-in behavior in order to test the agents in different situations (e.g. fast vs slow ingestion).
- A runner that contains the benchmarking logic. 
  For a benchmark, it'll run the baseline, non-profiled program along with different profiled versions in different situations (no ingester, a fast ingester, etc.).
- The set of benchmarks, which are generic programs that may behave differently.

The benchmarked programs are dockerized, which should make it easier to reproduce the benchmarks, and benchmarks are based on running time.

The runner uses the docker SDK instead of building the suite on top of something like docker-compose to have more flexibility.

Pull-mode is not supported for now, but support should be added at some point.

## The benchmarks

Each benchmark is just a directory with a Dockerfile that speficies how to build the image of the program to benchmark.

The program to benchmark should handle the PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING environment variable:
- When the envvar is set, the benchmarked program should enable profiling.
- When the envvar is not set, the benchmarked program should not enable profiling.

### List of benchmarks

- fibonacci. It's a CPU intensive benchmark that requires no heap allocation and no concurrency. It should provide simple stack traces of some height.
- memory / GC intensive benchmark (TODO)
- concurrency intensive benchmark (TODO)
- ...

## The runner

The runner takes a sequence of paths to the actual benchmarks to run as command line arguments.

The runner currently executes three different versions of the benchmarked program:
- A program with profiling not enabled. This is the baseline.
- A program with profiling enabled with an ingester available.
- A program with profiling enabled with an ingester unavailable.
It will measure the time it takes to run all of them and compare the last two with the baseline.
The benchmarked programs are run several times to have more samples and generate more reliable results.

The runner takes care of the whole setup and teardown of each benchmark, including:
- Building the docker images and containres for the ingester and benchmarked program.
- Creating a network and connecting/disconnecting the ingester and benchmarked program.
- Removing the containers, images and network when no longer needed.

## Usage

To run all the benchmars just run `run-all.sh`:

```
$ ./run-all.sh
2022/02/22 15:24:12 Running fibonacci-go-cpu-push benchmark
[...]
```

Take a look to the script to get a better idea of how it works.
