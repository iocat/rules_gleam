import demo/internal

// The main Fibonacci function that initiates the helper function.
// It calculates the nth Fibonacci number.
pub fn fibonacci(n: Int) -> Int {
  internal.fibonacci_helper(n, 0, 1) // Start with F(0) = 0 and F(1) = 1
}
