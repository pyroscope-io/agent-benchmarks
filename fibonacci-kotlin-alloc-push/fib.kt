fun fib(n: Long): Long = if (n < 2) n else fib(n-1) + fib(n-2)

fun main() {
    fib(50);
}
