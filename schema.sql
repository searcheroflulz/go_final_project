-- schema.sql
CREATE TABLE IF NOT EXISTS scheduler (
   id INTEGER PRIMARY KEY AUTOINCREMENT,
   date TEXT,
   title TEXT,
   comment TEXT,
   repeat TEXT
);