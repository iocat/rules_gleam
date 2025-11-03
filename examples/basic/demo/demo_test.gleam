import demo/demo

pub fn fibonacci_test() {
    assert demo.fibonacci(0) == 0
    assert demo.fibonacci(1) == 1
    assert demo.fibonacci(2) == 1
    assert demo.fibonacci(3) == 2
    assert demo.fibonacci(4) == 3
    assert demo.fibonacci(5) == 5
    assert demo.fibonacci(6) == 8
    assert demo.fibonacci(7) == 13
    assert demo.fibonacci(8) == 21
    assert demo.fibonacci(9) == 34
    assert demo.fibonacci(10) == 55
}
