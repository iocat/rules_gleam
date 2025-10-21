import gleeunit/should
import greeter

pub fn main() {
  gleeunit.main()
}

pub fn greeter_test() {
  greeter.greet("Phat")
    |> should.equal("Hello, Phat!")
}
