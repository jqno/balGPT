CREATE TABLE teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL
);

CREATE INDEX idx_teams_id ON teams(id);
CREATE INDEX idx_teams_name ON teams(name);
