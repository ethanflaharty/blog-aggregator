-- +goose up
Create table users(
    id integer primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null
);

-- +goose down
Drop table users;
