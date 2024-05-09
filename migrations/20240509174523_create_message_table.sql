-- +goose Up
create table messages (
    id serial primary key,
    "from" text not null,
    "text" text not null,
    timestamp timestamptz not null default now(),
    chat_id integer references chats(id) on delete set null
);

-- +goose Down
drop table messages;
