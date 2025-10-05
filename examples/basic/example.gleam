import demo/demo.{
  type Tree, create_node, empty, fibonacci, sorted_list_to_balanced_tree,
}
import gleam/int
import gleam/io
import gleam/json
import gleam/option.{None}

type Cat {
  Cat(name: String, lives: Int, flaws: option.Option(String))
}

fn cat_to_json(cat: Cat) -> String {
  json.object([
    #("name", json.string(cat.name)),
    #("lives", json.int(cat.lives)),
    #("flaws", json.null()),
  ])
  |> json.to_string
}

@external(erlang, "example_ffi", "tail_recursive")
fn fib_ffi(a: Int) -> Int

@external(erlang, "example_ffi", "naive")
fn naive(a: Int) -> Int

pub fn main() {
  io.println(fibonacci(3) |> int.to_string)

  let tree: Tree(Int) = create_node(1, empty, empty)
  echo tree
  echo sorted_list_to_balanced_tree([1, 3, 3, 4, 5, 6, 19, 12])

  io.println(fib_ffi(100) |> int.to_string)
  io.println(naive(20) |> int.to_string)
  io.println("Hello " <> "JS" <> "!")
  io.println(Cat("Dude", 9, None) |> cat_to_json)
}
