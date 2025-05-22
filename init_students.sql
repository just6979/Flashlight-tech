DROP TABLE students;

CREATE TABLE students (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(256) NOT NULL,
    grade INTEGER NOT NULL
);

INSERT INTO students (name, grade) VALUES('Alice', 90);
INSERT INTO students (name, grade) VALUES('Bob', 95);
INSERT INTO students (name, grade) VALUES('Marty', 88);
INSERT INTO students (name, grade) VALUES('David', 88);
INSERT INTO students (name, grade) VALUES('Wayne', 99);
