set statement_timeout = 0;
set client_encoding = 'UTF8';
-- set standart_conforming_strings = on;
set check_function_bodies = FALSE;
set client_min_messages = WARNING;
set search_path = public, extension;
set default_tablespace = '';
set default_with_oids = FALSE;

-- EXTENSIONS --
create extension if not exists pgcrypto;


-- TABLES --
create table if not exists users
(
    id uuid primary key default gen_random_uuid(),
    username varchar(255) not null unique,
    email text not null unique,
    pass_hash bytea not null,
    created_at timestamp default now()
);
create index if not exists index_email on users(email);

create table if not exists apps
(
    id uuid primary key default gen_random_uuid(),
    name text not null unique,
    description text,
    secret text not null unique,
    created_at timestamp default now()
);

create table if not exists groups
(
    id serial not null unique,
    app_id uuid not null references apps(id),
    name text not null unique,
    description text
);

create table if not exists roles
(
    id serial not null unique,
    name text not null unique,
    description text
);

create table if not exists admins
(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    app_id uuid not null references apps(id),
    is_admin bool not null default false,
    created_at timestamp default now()
);

create table if not exists groups_roles
(
    id serial not null unique,
    group_id integer not null references groups(id),
    role_id integer not null references roles(id),
    created_at timestamp default now()
);


create table if not exists users_permissions
(
    id uuid primary key default gen_random_uuid(),
    user_id uuid not null references users(id),
    group_id integer not null references groups(id),
    add_flag bool not null default false,
    created_at timestamp default now()
);


commit;

