# Pyroscope Agent Benchmark Suite

The benchmark suite is designed around three components:
- A ingestor that represents a pyroscope server.
  It can support built-in behavior in order to test the agents in different situations (e.g. fast vs slow ingestion).
- A runner that contains the benchmarking logic. 
  For a benchmark, it'll run the baseline, non-profiled program along with different profiled versions in different situations (no ingestor, a fast ingestor, etc.).
- The set of benchmarks, which are generic programs that may behave differently.

The benchmarked programs are dockerized, which should make it easier to reproduce the benchmarks, and benchmarks are based on running time.

The runner uses the docker SDK instead of building the suite on top of something like docker-compose to have more flexibility.

Pull-mode is not supported for now, but support should be added at some point.
