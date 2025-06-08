
-- Version: 1
-- Description: Create teams table
CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    budget BIGINT NOT NULL
);

-- Version: 2
-- Description: Create players table
CREATE TABLE players (
    id SERIAL PRIMARY KEY,
    team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    age INTEGER NOT NULL,
    country VARCHAR(100) NOT NULL,
    value BIGINT NOT NULL,
    position VARCHAR(50) NOT NULL
);

-- Version: 3
-- Description: Create users table
CREATE TABLE users (
    username VARCHAR(100) PRIMARY KEY,
    password_hash TEXT NOT NULL,
    team_id INTEGER REFERENCES teams(id) ON DELETE SET NULL,
    date_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    date_updated TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Version: 4
-- Description: Create transfers table
CREATE TABLE transfers (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    seller_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    asking_price BIGINT NOT NULL,
    status VARCHAR(50) NOT NULL
);
