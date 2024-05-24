create table if not exists accounts (
    id bigserial primary key,
    email text unique not null,
    username text unique not null,
    first_name text not null,
    last_name text not null,
    password text not null,
    profile_type text not null,
);

create table if not exists sessions (
    id bigserial primary key,
    account_id bigserial references accounts(id),
    token text not null,
    expiry_time timestamptz not null
);

create table if not exists developers (
    id bigserial primary key,
    name text,
    created_by bigserial references accounts(id),
    members jsonb
);

create table if not exists publishers (
    id bigserial primary key,
    name text,
    created_by bigserial references accounts(id),
    members jsonb
);

create table if not exists games (
    id bigserial primary key,
	developer_id bigserial references developers(id),
	name text not null,
	description text not null,
	genre text not null,
	platform text not null
)