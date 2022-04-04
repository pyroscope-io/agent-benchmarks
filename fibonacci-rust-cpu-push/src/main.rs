use pyroscope::PyroscopeAgent;
use std::env;

fn fib(n: i64) -> i64 {
    if n < 2 {
        n
    } else {
        fib(n - 1) + fib(n - 2)
    }
}

fn run() {
    fib(51);
}

fn main() -> Result<(), Box<dyn std::error::Error>> {
    if env::var("PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING").is_ok() {
        let mut agent =
            PyroscopeAgent::builder("http://ingester:4040", "fibonacci-rust-cpu-push").build()?;
        agent.start();
        run();
        agent.stop();
    } else {
        run();
    }
    Ok(())
}
