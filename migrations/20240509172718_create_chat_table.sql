-- +goose Up
create table chats (
    id serial primary key,
    usernames text[]
);

-- +goose Down
drop table chats;
