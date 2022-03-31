var DataTypes = require("sequelize").DataTypes;
var _SequelizeMeta = require("./SequelizeMeta");
var _decks = require("./decks");
var _users = require("./users");

function initModels(sequelize) {
  var SequelizeMeta = _SequelizeMeta(sequelize, DataTypes);
  var decks = _decks(sequelize, DataTypes);
  var users = _users(sequelize, DataTypes);

  decks.belongsTo(users, { as: "user", foreignKey: "user_id"});
  users.hasMany(decks, { as: "decks", foreignKey: "user_id"});

  return {
    SequelizeMeta,
    decks,
    users,
  };
}
module.exports = initModels;
module.exports.initModels = initModels;
module.exports.default = initModels;
