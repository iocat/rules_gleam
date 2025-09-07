// A tail-recursive helper function to calculate the nth Fibonacci number.
// It takes the current index 'n', and the two previous Fibonacci numbers 'a' and 'b'.
pub fn fibonacci_helper(n: Int, a: Int, b: Int) -> Int {
  case n {
    0 -> a // Base case: if n is 0, return the first number (0)
    1 -> b // Base case: if n is 1, return the second number (1)
    _ -> fibonacci_helper(n - 1, b, a + b) // Recursive step: call with n-1, new 'a' is old 'b', new 'b' is sum of old 'a' and 'b'
  }
}