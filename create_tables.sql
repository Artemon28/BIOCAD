create table devices
(
    mqtt       varchar(20),
    invid      varchar(50),
    unit_guid  varchar(50),
    msg_id     varchar(50),
    text       varchar(50),
    context    varchar(50),
    class      varchar(50),
    level      integer,
    area       varchar(50),
    addr       varchar(50),
    block      varchar(50),
    type       varchar(50),
    bit        integer,
    invert_bit integer,
    id         serial
        constraint devices_pk
            primary key
);

create unique index devices_id_uindex
    on devices using ??? (id);

create table files
(
    id   serial
        constraint files_pk
            primary key,
    name varchar(50)
);

create unique index files_id_uindex
    on files (id);

create unique index files_name_uindex
    on files (name);