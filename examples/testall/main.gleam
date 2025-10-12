import lustre
import lustreserver/server

pub fn main() {
  let app = lustre.simple(server.init, server.update, server.view)
  let assert Ok(_) = lustre.start(app, "#app", Nil)

  Nil
}