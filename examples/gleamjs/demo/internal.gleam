// A tail-recursive helper function to calculate the nth Fibonacci number.
// It takes the current index 'n', and the two previous Fibonacci numbers 'a' and 'b'.
pub fn fibonacci_helper(n: Int, a: Int, b: Int) -> Int {
  case n {
    0 -> a // Base case: if n is 0, return the first number (0)
    1 -> b // Base case: if n is 1, return the second number (1)
    _ -> fibonacci_helper(n - 1, b, a + b) // Recursive step: call with n-1, new 'a' is old 'b', new 'b' is sum of old 'a' and 'b'
  }
}

// No imports from the Gleam standard library are used.

/// A generic binary tree, which is either empty or a node with a value
/// and two child trees (left and right).
pub type Tree(a) {
  Empty
  Node(value: a, left: Tree(a), right: Tree(a))
}

// A recursive helper function to get the length of a list.
fn length(of list: List(a)) -> Int {
  case list {
    [] -> 0
    [_, ..tail] -> 1 + length(of: tail)
  }
}

fn split_at(list: List(a), index: Int) -> #(List(a), List(a)) {
  case list, index {
    // If index is 0 or less, the first part is empty.
    l, n if n <= 0 -> #([], l)
    // If the list is empty, both parts are empty.
    [], _ -> #([], [])
    // The recursive step: peel off the head and recurse on the tail.
    [head, ..tail], n -> {
      let #(first_part, second_part) = split_at(tail, n - 1)
      #([head, ..first_part], second_part)
    }
  }
}

/// ## Converts a Sorted List into a Balanced Binary Search Tree
pub fn sorted_list_to_balanced_tree(from list: List(a)) -> Tree(a) {
  let size = length(of: list)
  case size {
    0 -> Empty
    _ -> {
      let middle_index = size / 2
      // Use our custom helper to split the list at the middle.
      let #(left_items, remainder) = split_at(list, middle_index)

      // The remainder contains the root and all right-side items.
      // We pattern match to extract the root and the rest.
      case remainder {
        [root_value, ..right_items] ->
          Node(
            value: root_value,
            // Recursively build the left and right subtrees.
            left: sorted_list_to_balanced_tree(from: left_items),
            right: sorted_list_to_balanced_tree(from: right_items),
          )

        // This case is for an empty remainder, which shouldn't happen
        // if size > 0, but makes the pattern match exhaustive.
        [] -> Empty
      }
    }
  }
}
