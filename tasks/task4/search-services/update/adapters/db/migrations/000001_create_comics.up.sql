CREATE TABLE comics (
    comics_id INTEGER NOT NULL PRIMARY KEY CHECK(comics_id > 0),
    comics_url TEXT NOT NULL,
    words TEXT[]
);



