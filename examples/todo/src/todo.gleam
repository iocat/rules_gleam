import gleam/io

pub type Todo {
  Todo {
    id: Int,
    title: String,
    completed: Bool,
  }
}

pub fn create(id: Int, title: String) -> Todo {
  Todo(id: id, title: title, completed: False)
}

pub fn toggle(todo: Todo) -> Todo {
  Todo(..todo, completed: !todo.completed)
}

pub fn to_string(todo: Todo) -> String {
  let status = case todo.completed {
    True -> "✅"
    False -> "⬜"
  }
  status <> " " <> todo.title
}

pub fn main() {
  let todo = create(1, "Learn Gleam")
  let toggled = toggle(todo)

  io.println("Original: " <> to_string(todo))
  io.println("Toggled: " <> to_string(toggled))
}
