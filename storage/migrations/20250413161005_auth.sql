-- +goose Up
-- +goose StatementBegin
create table users (
    id serial primary key,
    name text unique not null,
    email text not null,
    password text not null,
    role int not null
);

CREATE UNIQUE INDEX idx_users_username ON users(name);

create table sessions (
    id serial primary key,
    user_id serial not null references "users"(id) on delete cascade,
    refresh_token text not null unique,
    created_at timestamp not null default now(),
    expires_at timestamp not null,
    revoked_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
drop table if exists sessions;
-- +goose StatementEnd
