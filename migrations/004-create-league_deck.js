'use strict';
module.exports = {
  async up(queryInterface, Sequelize) {
    await queryInterface.createTable('league_decks', {
      id: {
        allowNull: false,
        autoIncrement: true,
        primaryKey: true,
        type: Sequelize.INTEGER
      },
      league_id: {
        allowNull: false,
        type: Sequelize.INTEGER,
        references: { model: 'leagues', key: 'id' }
      },
      deck_id: {
        allowNull: false,
        type: Sequelize.INTEGER,
        references: { model: 'decks', key: 'id' }
      },
      created_at: {
        allowNull: false,
        type: Sequelize.DATE,
        defaultValue: Sequelize.fn('NOW')
      }
    });
  },
  async down(queryInterface, Sequelize) {
    await queryInterface.dropTable('league_decks');
  }
};