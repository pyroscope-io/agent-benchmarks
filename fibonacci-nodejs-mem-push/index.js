const pyroscope = require('@pyroscope/nodejs');

function fibonacci(n) {
    if ( n < 2 ) { return n; }

    return fibonacci(n - 1) + fibonacci(n - 2);
}

const num = 49;

if ( process.env['PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING'] ) {
    pyroscope.init({server: 'http://ingester:4000', autoStart: false, name: 'fibonacci-nodejs-mem-push'});
    pyroscope.startHeapProfiling();
    fibonacci(num);
    setTimeout(() => { console.log("Disabling"); pyroscope.stopHeapProfiling()} );
    process.exit(0);
} else {
    fibonacci(num);
    process.exit(0);
}
