import multiplegleamonemain/character.{Character}
import multiplegleamonemain/faction.{Faction}

pub type Guild {
  Guild(
    name: String,
    faction: Faction,
    members: List(Character),
  )
}

pub fn new(name: String, faction: Faction) -> Guild {
  Guild(name: name, faction: faction, members: [])
}

pub fn add_member(guild: Guild, character: Character) -> Guild {
  case faction.is_same(guild.faction, character.faction) {
    True -> Guild(..guild, members: [character, ..guild.members])
    False -> guild
  }
}
