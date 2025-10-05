CREATE TABLE users
(
    id serial not null unique,
    first_name varchar(255) not null,
    last_name varchar(255) not null,
    username varchar(255) not null unique,
    email varchar(255) not null unique,
    password_hash varchar(255)
);

CREATE TABLE rooms
(
    id serial not null unique,
    name varchar(255) not null unique,
    description varchar(255),
    created_by int references users(id) on delete cascade not null,
    created_at timestamp default current_timestamp
);

CREATE TABLE clients
(
    id serial not null unique,
    room_id int references rooms(id) on delete cascade not null,
    user_id int references users(id) on delete cascade not null,
    connected_at timestamp default current_timestamp,
    disconnected_at timestamp
);

CREATE TABLE messages
(
    id serial not null unique,
    room_id int references rooms(id) on delete cascade not null,
    user_id int references users(id) on delete cascade not null,
    content text not null,
    created_at timestamp default current_timestamp
);
