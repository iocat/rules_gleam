import todo
import gleeunit
import gleeunit/should

pub fn main() {
  gleeunit.main()
}

pub fn todo_creation_test() {
  let todo = todo.create(1, "Learn Gleam")
  todo.title
    |> should.equal("Learn Gleam")
  todo.completed
    |> should.be_false()
}

pub fn todo_toggle_test() {
  let todo = todo.create(1, "Learn Gleam")
    |> todo.toggle()
  todo.completed
    |> should.be_true()
}
