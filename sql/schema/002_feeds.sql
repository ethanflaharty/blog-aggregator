-- +goose up
Create table feeds(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text not null,
    url text unique not null,
    user_id uuid not null,
    Foreign key (user_id)
        References users(id)
        on delete cascade
);

-- +goose down
drop table feeds;
