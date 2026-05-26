-- name: CreateFeed :one
insert into feeds (id, created_at, updated_at, name, url, user_id)
values (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
returning *;

-- name: GetFeeds :many
select feeds.*, users.name
from feeds
inner join users on users.id = feeds.user_id;

-- name: GetFeedByURL :one
select * from feeds where url = $1;