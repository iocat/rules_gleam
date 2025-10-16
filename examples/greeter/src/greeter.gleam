import gleam/io

pub fn greet(name: String) -> String {
  "Hello, " <> name <> "!"
}

pub fn main() {
  let greeting = greet("World")
  io.println(greeting)
}
