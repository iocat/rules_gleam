import demo/demo.{type Tree, empty, fibonacci, create_node, sorted_list_to_balanced_tree}

@external(erlang, "example_ffi", "tail_recursive")
fn fib_ffi(a: Int) -> Int

@external(erlang, "example_ffi", "naive")
fn naive(a: Int) -> Int

pub fn main() {
  echo fibonacci(3)

  let tree: Tree(Int) = create_node(1, empty, empty)
  echo tree

  echo sorted_list_to_balanced_tree([1 , 3, 3, 4, 5, 6, 19, 12])

  echo fib_ffi(100)
  echo naive(20)
}
