-- +goose Up
-- +goose StatementBegin
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
drop table if exists sessions;
-- +goose StatementEnd
