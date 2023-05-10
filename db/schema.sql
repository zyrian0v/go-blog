CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY,
    slug TEXT NOT NULL UNIQUE, 
    title TEXT NOT NULL,
    content TEXT
);