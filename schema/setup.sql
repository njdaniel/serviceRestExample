-- docker run --name dnd_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres:latest

CREATE TABLE items (
                       name varchar UNIQUE NOT NULL ,
                       cost varchar,
                       description varchar
);