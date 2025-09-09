
/// A faction represents a group a character can belong to.
pub type Faction {
  Knights
  Wizards
  Rogues
}

pub fn is_same(a: Faction, b: Faction) -> Bool {
  case a, b {
    Knights, Knights -> True
    Wizards, Wizards -> True
    Rogues, Rogues -> True
    _, _ -> False
  }
}