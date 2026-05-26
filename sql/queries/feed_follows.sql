-- name: CreateFeedFollow :one
with inserted_feed_follow as (
    insert into feed_follows (id, created_at, updated_at, user_id, feed_id)
    values (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    returning *
)

select 
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
from inserted_feed_follow
inner join users on users.id = inserted_feed_follow.user_id
inner join feeds on feeds.id = inserted_feed_follow.feed_id;

-- name: GetFeedFollowsForUser :many
select 
    feed_follows.*, 
    feeds.name as feed_name,
    users.name as user_name
from feed_follows
inner join users on users.id = feed_follows.user_id
inner join feeds on feeds.id = feed_follows.feed_id
where feed_follows.user_id = $1;