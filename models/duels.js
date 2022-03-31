const Sequelize = require('sequelize');
module.exports = function(sequelize, DataTypes) {
  return sequelize.define('duels', {
    id: {
      autoIncrement: true,
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true
    },
    league_id: {
      type: DataTypes.INTEGER,
      allowNull: true,
      references: {
        model: 'leagues',
        key: 'id'
      }
    },
    user_1_id: {
      type: DataTypes.STRING(255),
      allowNull: false,
      references: {
        model: 'users',
        key: 'id'
      }
    },
    user_2_id: {
      type: DataTypes.STRING(255),
      allowNull: false,
      references: {
        model: 'users',
        key: 'id'
      }
    },
    deck_1_id: {
      type: DataTypes.INTEGER,
      allowNull: false,
      references: {
        model: 'decks',
        key: 'id'
      }
    },
    deck_2_id: {
      type: DataTypes.INTEGER,
      allowNull: false,
      references: {
        model: 'decks',
        key: 'id'
      }
    },
    result: {
      type: DataTypes.INTEGER,
      allowNull: false
    },
    confirmed_at: {
      type: DataTypes.DATE,
      allowNull: true
    }
  }, {
    sequelize,
    tableName: 'duels',
    schema: 'public',
    timestamps: true,
    paranoid: true,
    indexes: [
      {
        name: "duels_pkey",
        unique: true,
        fields: [
          { name: "id" },
        ]
      },
    ]
  });
};
