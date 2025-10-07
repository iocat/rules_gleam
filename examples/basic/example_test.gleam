import example
import gleam/option.{None}

pub fn cat_to_json_test() {
  let cat = example.Cat("Dude", 9, None)
  let json = example.cat_to_json(cat)

  assert json == "{\"name\":\"Dude\",\"lives\":9,\"flaws\":null}"
}
