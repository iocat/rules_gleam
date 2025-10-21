import todo
import gleeunit
import gleeunit/should

type Todo {
  Todo(
    id: Int,
    title: String,
    description: Option(String),
    completed: Bool,
    priority: todo.Priority,
    due_date: Option(Int),
    tags: List(todo.Tag),
  )
}

pub fn main() {
  gleeunit.main()
}

pub fn todo_creation_test() {
  let todo = todo.create(1, "Learn Gleam")
  case todo {
    Todo(todo) -> {
      todo.title
        |> should.equal("Learn Gleam")
      todo.completed
        |> should.be_false()
    }
  }
}

pub fn todo_toggle_test() {
  let todo = todo.create(1, "Learn Gleam")
    |> todo.toggle()
  case todo {
    Todo(todo) -> {
      todo.completed
        |> should.be_true()
    }
  }
}
