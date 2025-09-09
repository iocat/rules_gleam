import multiplegleamonemain/faction
import multiplegleamonemain/character
import multiplegleamonemain/guild

pub fn main() {
  let knights_faction = faction.Knights
  let wizards_faction = faction.Wizards

  let arthur = character.new("Arthur", knights_faction)
  let lancelot = character.new("Lancelot", knights_faction)
  let merlin = character.new("Merlin", wizards_faction)

  let round_table = guild.new("The Round Table", knights_faction)

  let guild_with_arthur = guild.add_member(round_table, arthur)
  let guild_with_knights = guild.add_member(guild_with_arthur, lancelot)
  let final_guild = guild.add_member(guild_with_knights, merlin)

  let _ = final_guild

  echo final_guild
}