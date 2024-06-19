CREATE TABLE articles (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE, 
    content TEXT,
    created_at DATETIME NOT NULL
);