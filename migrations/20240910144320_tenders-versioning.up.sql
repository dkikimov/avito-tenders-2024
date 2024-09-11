CREATE TABLE tenders_history
(
    tender_id        uuid references tenders (id),
    name             text      not null,
    description      text      not null,
    service_type     text      not null,
    status           text      not null,
    organization_id  uuid      not null references organization (id),
    creator_username text      not null references employee (username),
    created_at       timestamp not null,
    version          int       not null,
    modified_at      timestamp not null default now(),
    primary key (tender_id, version)
);

CREATE OR REPLACE FUNCTION log_tender_update() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO tenders_history(tender_id, name, description, service_type, status, organization_id, creator_username,
                               created_at, version)
    VALUES (OLD.id, OLD.name, OLD.description, OLD.service_type, OLD.status, OLD.organization_id, OLD.creator_username,
            OLD.created_at, OLD.version);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tender_update_trigger
    BEFORE UPDATE
    ON tenders
    FOR EACH ROW
EXECUTE FUNCTION log_tender_update();