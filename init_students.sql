create table if not exists students (
    id serial primary key,
    name varchar(256) not null,
    grade integer not null
);

insert into students (name, grade) values('Alice', 90);
insert into students (name, grade) values('Bob', 95);
