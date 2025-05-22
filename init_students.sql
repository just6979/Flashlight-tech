DROP TABLE students;

CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    grade INTEGER NOT NULL
);

INSERT INTO students (name, grade) VALUES('Alice', 90);
INSERT INTO students (name, grade) VALUES('Bob', 95);
