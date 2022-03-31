const Sequelize = require('sequelize');
module.exports = function(sequelize, DataTypes) {
  return sequelize.define('league_decks', {
    id: {
      autoIncrement: true,
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true
    },
    league_id: {
      type: DataTypes.INTEGER,
      allowNull: false,
      references: {
        model: 'leagues',
        key: 'id'
      }
    },
    deck_id: {
      type: DataTypes.INTEGER,
      allowNull: false,
      references: {
        model: 'decks',
        key: 'id'
      }
    }
  }, {
    sequelize,
    tableName: 'league_decks',
    schema: 'public',
    timestamps: true,
    indexes: [
      {
        name: "league_decks_pkey",
        unique: true,
        fields: [
          { name: "id" },
        ]
      },
    ]
  });
};
