exports.up = (pgm) => {
  pgm.createTable(
    'user', 
    {
      id: {
        primaryKey: true,
        type: 'varchar(256)',
        notNull: true,
      },
      password: {
        type: 'varchar(256)',
        notNull: true,
      },
      name: {
        type: 'varchar(256)',
        notNull: true,
      },
      created_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      updated_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      deleted_at: {
        type: 'timestamp',
      },
    }
  )
  pgm.createTable(
    'deck', 
    {
      id: 'id',
      name: {
        type: 'varchar(256)',
        notNull: true,
      },
      description: {
        type: 'varchar(256)',
        notNull: true,
      },
      user_id: {
        type: 'varchar(256)',
        notNull: true,
        references: '"user"',
        onDelete: 'cascade',
      },
      created_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      updated_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      deleted_at: {
        type: 'timestamp'
      },
    }
  )
  pgm.createTable(
    'league', 
    {
      id: 'id',
      name: {
        type: 'varchar(256)',
        notNull: true,
      },
      color: {
        type: 'varchar(256)',
        notNull: true,
      },
      created_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      updated_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      deleted_at: {
        type: 'timestamp'
      },
    }
  )
  pgm.createTable(
    'league_deck', 
    {
      id: 'id',
      league_id: {
        type: 'integer',
        notNull: true,
        references: '"league"',
        onDelete: 'cascade',
      },
      deck_id: {
        type: 'integer',
        notNull: true,
        references: '"deck"',
        onDelete: 'cascade',
      },
      created_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
    }
  )
  pgm.createTable(
    'duel', 
    {
      id: 'id',
      league_id: {
        type: 'integer',
        notNull: true,
        references: '"league"',
        onDelete: 'cascade',
      },
      user_1_id: {
        type: 'varchar(256)',
        notNull: true,
        references: '"user"',
        onDelete: 'cascade',
      },
      user_2_id: {
        type: 'varchar(256)',
        notNull: true,
        references: '"user"',
        onDelete: 'cascade',
      },
      deck_1_id: {
        type: 'integer',
        notNull: true,
        references: '"deck"',
        onDelete: 'cascade',
      },
      deck_2_id: {
        type: 'integer',
        notNull: true,
        references: '"deck"',
        onDelete: 'cascade',
      },
      result: {
        type: 'integer',
        notNull: true,
      },
      created_at: {
        type: 'timestamp',
        notNull: true,
        default: pgm.func('current_timestamp'),
      },
      confirmed_at: {
        type: 'timestamp',
      },
      deleted_at: {
        type: 'timestamp',
      },
    }
  )
}