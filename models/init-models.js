var DataTypes = require("sequelize").DataTypes;
var _SequelizeMeta = require("./SequelizeMeta");
var _decks = require("./decks");
var _duels = require("./duels");
var _league_decks = require("./league_decks");
var _leagues = require("./leagues");
var _users = require("./users");

function initModels(sequelize) {
  var SequelizeMeta = _SequelizeMeta(sequelize, DataTypes);
  var decks = _decks(sequelize, DataTypes);
  var duels = _duels(sequelize, DataTypes);
  var league_decks = _league_decks(sequelize, DataTypes);
  var leagues = _leagues(sequelize, DataTypes);
  var users = _users(sequelize, DataTypes);

  duels.belongsTo(decks, { as: "deck_1", foreignKey: "deck_1_id"});
  decks.hasMany(duels, { as: "duels", foreignKey: "deck_1_id"});
  duels.belongsTo(decks, { as: "deck_2", foreignKey: "deck_2_id"});
  decks.hasMany(duels, { as: "deck_2_duels", foreignKey: "deck_2_id"});
  league_decks.belongsTo(decks, { as: "deck", foreignKey: "deck_id"});
  decks.hasMany(league_decks, { as: "league_decks", foreignKey: "deck_id"});
  duels.belongsTo(leagues, { as: "league", foreignKey: "league_id"});
  leagues.hasMany(duels, { as: "duels", foreignKey: "league_id"});
  league_decks.belongsTo(leagues, { as: "league", foreignKey: "league_id"});
  leagues.hasMany(league_decks, { as: "league_decks", foreignKey: "league_id"});
  decks.belongsTo(users, { as: "user", foreignKey: "user_id"});
  users.hasMany(decks, { as: "decks", foreignKey: "user_id"});
  duels.belongsTo(users, { as: "user_1", foreignKey: "user_1_id"});
  users.hasMany(duels, { as: "duels", foreignKey: "user_1_id"});
  duels.belongsTo(users, { as: "user_2", foreignKey: "user_2_id"});
  users.hasMany(duels, { as: "user_2_duels", foreignKey: "user_2_id"});

  return {
    SequelizeMeta,
    decks,
    duels,
    league_decks,
    leagues,
    users,
  };
}
module.exports = initModels;
module.exports.initModels = initModels;
module.exports.default = initModels;
