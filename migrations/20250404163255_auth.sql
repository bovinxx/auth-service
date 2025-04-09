-- +goose Up
-- +goose StatementBegin
create table users (
    id serial primary key,
    name text unique not null,
    email text not null,
    password text not null,
    is_admin boolean not null
);

CREATE UNIQUE INDEX idx_users_username ON users(name);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
