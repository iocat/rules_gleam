import gleam/io
import gleam/json

pub fn main() {
  let deps = io.args()
    |> io.map(fn(arg) {
      let content = io.read_file(arg)
        |> io.unwrap_or_exit("Failed to read file: \(arg)")
      json.decode(content)
        |> io.unwrap_or_exit("Failed to decode JSON in file: \(arg)")
    })
    |> io.collect()

  io.println("Dependencies: \(deps)")
}