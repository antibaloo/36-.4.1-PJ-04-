DROP TABLE IF EXISTS news;
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL CHECK (title <> ''),
    content TEXT NOT NULL CHECK (content <> ''),
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE CHECK (link <> '')
);