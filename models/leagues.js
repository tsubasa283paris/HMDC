const Sequelize = require('sequelize');
module.exports = function(sequelize, DataTypes) {
  return sequelize.define('leagues', {
    id: {
      autoIncrement: true,
      type: DataTypes.INTEGER,
      allowNull: false,
      primaryKey: true
    },
    name: {
      type: DataTypes.STRING(255),
      allowNull: false
    }
  }, {
    sequelize,
    tableName: 'leagues',
    schema: 'public',
    timestamps: true,
    paranoid: true,
    indexes: [
      {
        name: "leagues_pkey",
        unique: true,
        fields: [
          { name: "id" },
        ]
      },
    ]
  });
};
