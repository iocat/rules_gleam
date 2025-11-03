// Import the module that contains the Tree type and helpers.
import demo/internal

// ---------------------------------------------------------------------------
// Test: constructing a simple tree manually
pub fn manual_construction_test() {
  // Build a small tree:
  //        2
  //       / \
  //      1   3
  let _tree =
    internal.Node(
      value: 2,
      left: internal.Node(value: 1, left: internal.Empty, right: internal.Empty),
      right: internal.Node(
        value: 3,
        left: internal.Empty,
        right: internal.Empty,
      ),
    )
}

pub fn sorted_list_to_balanced_tree_test() {
  let tree = internal.sorted_list_to_balanced_tree([1, 10, 2, 3, 2, 12, 4])
  assert tree
    == internal.Node(
      value: 3,
      left: internal.Node(
        value: 10,
        left: internal.Node(
          value: 1,
          left: internal.Empty,
          right: internal.Empty,
        ),
        right: internal.Node(
          value: 2,
          left: internal.Empty,
          right: internal.Empty,
        ),
      ),
      right: internal.Node(
        value: 12,
        left: internal.Node(
          value: 2,
          left: internal.Empty,
          right: internal.Empty,
        ),
        right: internal.Node(
          value: 4,
          left: internal.Empty,
          right: internal.Empty,
        ),
      ),
    )
}
