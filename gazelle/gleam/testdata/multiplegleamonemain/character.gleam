
import multiplegleamonemain/faction.{Faction}

pub type Character {
  Character(name: String, faction: Faction)
}

pub fn new(name: String, faction: Faction) -> Character {
  Character(name: name, faction: faction)
}