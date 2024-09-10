create table tenders(
    id serial primary key,
    name text not null ,
    description text not null ,
    service_type text not null ,
    status text not null ,
    organization_id int not null references organization(id),
    creator_username text not null references employee(username),
    created_at timestamp not null default now(),
    version int default 1
)