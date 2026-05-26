-- +goose up
Create table feed_follows(
    id uuid primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    user_id uuid not null,
    Foreign key (user_id)
        References users(id)
        on delete cascade,
    feed_id uuid not null,
    Foreign key (feed_id)
        References feeds(id)
        on delete cascade,
    constraint following unique (user_id, feed_id)
);

-- +goose down
drop table feed_follows;