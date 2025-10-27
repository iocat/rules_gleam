import demo/demo.{
  type Tree, create_node, empty, sorted_list_to_balanced_tree,
}


@external(javascript, "./example_ffi.mjs", "add")
pub fn add(a: Int, b: Int) -> Int

pub fn main() {

  let tree: Tree(Int) = create_node(1, empty, empty)
  echo tree
  echo sorted_list_to_balanced_tree([1, 3, 3, 4, 5, 6, 19, 12])
  echo add(1, 22)

  // io.println(fib_ffi(100) |> int.to_string)
  // io.println(naive(20) |> int.to_string)
  // io.println("Hello " <> "JS" <> "!")
  // io.println(Cat("Dude", 9, None) |> cat_to_json)
}
