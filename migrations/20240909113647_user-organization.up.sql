CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE employee
(
    id         uuid primary key default uuid_generate_v4(),
    username   VARCHAR(50) UNIQUE NOT NULL,
    first_name VARCHAR(50),
    last_name  VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE organization_type AS ENUM (
    'IE',
    'LLC',
    'JSC'
    );

CREATE TABLE organization
(
    id          uuid primary key default uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    description TEXT,
    type        organization_type,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organization_responsible
(
    id              uuid primary key default uuid_generate_v4(),
    organization_id uuid REFERENCES organization (id) ON DELETE CASCADE,
    user_id         uuid REFERENCES employee (id) ON DELETE CASCADE
);