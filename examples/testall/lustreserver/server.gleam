import gleam/int
import lustre/element.{text}
import lustre/element/html.{div, button, p}
import lustre/event.{on_click}

pub fn init(_flags) {
  0
}

pub type Msg {
  Incr
  Decr
}

pub fn update(model, msg: Msg) {
  case msg {
    Incr -> model + 1
    Decr -> model - 1
  }
}

pub fn view(model) {
  let count = int.to_string(model)

  div([], [
    button([on_click(Incr)], [text(" + ")]),
    p([], [text(count)]),
    button([on_click(Decr)], [text(" - ")])
  ])
}