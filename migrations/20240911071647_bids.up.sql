create table bids
(
    id          uuid primary key default uuid_generate_v4(),
    name        text not null,
    description text not null,
    status      text not null,
    tender_id   uuid references tenders (id),
    author_type text not null,
    author_id   uuid not null,
    version     int  not null    default 1,
    created_at  timestamp        default now()
);

CREATE TABLE bids_history
(
    bid_id      uuid references bids (id),
    name        text      not null,
    description text      not null,
    status      text      not null,
    tender_id   uuid references tenders (id),
    author_type text      not null,
    author_id   text      not null,
    version     int       not null,
    created_at  timestamp,
    modified_at timestamp not null default now(),
    primary key (name, version)
);

CREATE OR REPLACE FUNCTION log_bid_update() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO bids_history(bid_id, name, description, status, tender_id, author_type, author_id, version, created_at)
    VALUES (old.id, old.name, old.description, old.status, old.tender_id, old.author_type, old.author_id, old.version, old.created_at);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER bid_update_trigger
    BEFORE UPDATE
    ON bids
    FOR EACH ROW
EXECUTE FUNCTION log_bid_update();