var DataTypes = require("sequelize").DataTypes;
var _SequelizeMeta = require("./SequelizeMeta");
var _decks = require("./decks");
var _league_decks = require("./league_decks");
var _leagues = require("./leagues");
var _users = require("./users");

function initModels(sequelize) {
  var SequelizeMeta = _SequelizeMeta(sequelize, DataTypes);
  var decks = _decks(sequelize, DataTypes);
  var league_decks = _league_decks(sequelize, DataTypes);
  var leagues = _leagues(sequelize, DataTypes);
  var users = _users(sequelize, DataTypes);

  league_decks.belongsTo(decks, { as: "deck", foreignKey: "deck_id"});
  decks.hasMany(league_decks, { as: "league_decks", foreignKey: "deck_id"});
  league_decks.belongsTo(leagues, { as: "league", foreignKey: "league_id"});
  leagues.hasMany(league_decks, { as: "league_decks", foreignKey: "league_id"});
  decks.belongsTo(users, { as: "user", foreignKey: "user_id"});
  users.hasMany(decks, { as: "decks", foreignKey: "user_id"});

  return {
    SequelizeMeta,
    decks,
    league_decks,
    leagues,
    users,
  };
}
module.exports = initModels;
module.exports.initModels = initModels;
module.exports.default = initModels;
