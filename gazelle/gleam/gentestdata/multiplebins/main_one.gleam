import gleam/io
import gleam/string as text

pub fn main() {
  io.println("Hello, Mike!")
  io.println(text.reverse("Hello, Joe!"))
}
