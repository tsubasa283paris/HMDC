'use strict';
module.exports = {
  async up(queryInterface, Sequelize) {
    await queryInterface.createTable('duels', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: Sequelize.INTEGER
      },
      league_id: {
        allowNull: true,
        type: Sequelize.INTEGER,
        references: { model: 'leagues', key: 'id' }
      },
      user_1_id: {
        allowNull: false,
        type: Sequelize.STRING,
        references: { model: 'users', key: 'id' }
      },
      user_2_id: {
        allowNull: false,
        type: Sequelize.STRING,
        references: { model: 'users', key: 'id' }
      },
      deck_1_id: {
        allowNull: false,
        type: Sequelize.INTEGER,
        references: { model: 'decks', key: 'id' }
      },
      deck_2_id: {
        allowNull: false,
        type: Sequelize.INTEGER,
        references: { model: 'decks', key: 'id' }
      },
      result: {
        allowNull: false,
        type: Sequelize.INTEGER
      },
      created_at: {
        allowNull: false,
        type: Sequelize.DATE,
        defaultValue: Sequelize.fn('NOW')
      },
      confirmed_at: {
        allowNull: true,
        type: Sequelize.DATE
      },
      deleted_at: {
        allowNull: true,
        type: Sequelize.DATE
      }
    });
  },
  async down(queryInterface, Sequelize) {
    await queryInterface.dropTable('duels');
  }
};