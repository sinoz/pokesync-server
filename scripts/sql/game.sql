DROP INDEX IF EXISTS account_email;
DROP INDEX IF EXISTS character_name;
DROP INDEX IF EXISTS character_owner;
DROP INDEX IF EXISTS pc_owner;
DROP INDEX IF EXISTS party_owner;

DROP TABLE IF EXISTS pc_entry;
DROP TABLE IF EXISTS party_entry;
DROP TABLE IF EXISTS monster;
DROP TABLE IF EXISTS character;
DROP TABLE IF EXISTS account;

DROP TYPE IF EXISTS user_group;
DROP TYPE IF EXISTS gender;
DROP TYPE IF EXISTS bicycle_type;

CREATE TYPE user_group AS ENUM (
    'regular',
    'patron',
    'mod',
    'admin',
    'game_design',
    'web_dev',
    'game_dev'
);

CREATE TYPE gender AS ENUM (
    'man',
    'woman',
    'genderless'
);

CREATE TYPE bicycle_type AS ENUM (
    'acro',
    'mach'
)

CREATE TABLE account (
    id serial PRIMARY KEY,
    email varchar(128) NOT NULL UNIQUE,
    password varchar(1024) NOT NULL
);

CREATE UNIQUE INDEX account_email ON account(email);

CREATE TABLE character (
    id serial PRIMARY KEY,
    display_name varchar(32) UNIQUE NOT NULL,
    user_group user_group NOT NULL DEFAULT 'regular',
    gender gender NOT NULL CHECK (gender = 'man' AND gender = 'woman'),
    bicycle_type bicycle_type NOT NULL,
    pokedollars integer DEFAULT 0 CHECK (pokedollars > 0),
    donator_points integer DEFAULT 0 CHECK (donator_points > 0),
    map_x smallint CHECK (map_x > 0),
    map_z smallint CHECK (map_z > 0),
    local_x smallint CHECK (local_x > 0),
    local_z smallint CHECK (local_z > 0),
    muted_until timestamp,
    banned_until timestamp,
    last_logged_in timestamp,
    account_id integer REFERENCES account (id)
);

CREATE UNIQUE INDEX character_name ON character(display_name);

CREATE INDEX character_owner ON character(account_id);

CREATE TABLE monster (
    id serial PRIMARY KEY
    model_id smallint CHECK (model_id >= 0)
    nickname varchar (32),
    original_trainer varchar (32) NOT NULL
);

CREATE TABLE pc_entry (
    id serial PRIMARY KEY,
    box_id smallint CHECK (box_id >= 0 AND box_id < 32),
    slot smallint CHECK (slot >= 0 AND slot < 100),
    monster_id integer REFERENCES monster (id),
    character_id integer REFERENCES character (id)
);

CREATE INDEX pc_owner ON pc_entry(character_id);

CREATE TABLE party_entry (
    id serial PRIMARY KEY,
    slot smallint CHECK (slot >= 0 AND slot < 6),
    monster_id integer REFERENCES monster (id),
    character_id integer REFERENCES character (id)
);

CREATE INDEX party_owner ON pc_entry(character_id);