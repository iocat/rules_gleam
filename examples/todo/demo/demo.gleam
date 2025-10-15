import todo

type Todo {
  Todo(
    id: Int,
    title: String,
    completed: Bool,
  )
}

pub fn main() {
  // Create a new todo
  let my_todo = todo.create(1, "Learn Gleam")
  
  // Toggle the todo
  let toggled = todo.toggle(my_todo)
  
  // Print the results
  io.println("Original todo: " <> todo.to_string(my_todo))
  io.println("After toggle: " <> todo.to_string(toggled))
}
