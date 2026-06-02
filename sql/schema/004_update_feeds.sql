-- +goose up
alter table feeds
add column last_fetched_at timestamp;

-- +goose down
alter table feeds
drop column last_fetched_at;