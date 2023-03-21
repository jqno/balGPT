CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    home_team INTEGER REFERENCES teams(id),
    away_team INTEGER REFERENCES teams(id),
    home_goals INTEGER,
    away_goals INTEGER,
    date DATE NOT NULL
);

CREATE INDEX idx_matches_id ON matches(id);
CREATE INDEX idx_matches_home_team ON matches(home_team);
CREATE INDEX idx_matches_away_team ON matches(away_team);
CREATE INDEX idx_matches_date ON matches(date);
