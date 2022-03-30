const pyroscope = require('@pyroscope/nodejs');

function fibonacci(n) {
    if ( n == 1 ) { return 1; }
    if ( n == 2 ) { return 1; }

    return fibonacci(n - 1) + fibonacci(n - 2);
}

const num = 48;

if ( process.env['PYROSCOPE_AGENT_BENCHMARK_ENABLE_PROFILING'] ) {
    pyroscope.init({server: 'http://ingester:4000', autoStart: false, name: 'fibonacci-nodejs-cpu-push'});
    pyroscope.startHeapProfiling();
    fibonacci(num);
    setTimeout(() => { console.log("Disabling"); pyroscope.stopHeapProfiling()} );
    process.exit(0);
} else {    
    pyroscope.init({autoStart: false});
    fibonacci(num);
    process.exit(0);

}
