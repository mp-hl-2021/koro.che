create table accounts
(
    id       serial primary key,
    login    varchar(255) not null unique,
    password varchar(255) not null
);

create table links
(
    id           serial primary key,
    creator_id   int default null,
    realLink     varchar(255),
    key varchar(255) unique,
    use_counter  int default 0,

    constraint fk_creator
        foreign key (creator_id)
            references accounts (id)
);