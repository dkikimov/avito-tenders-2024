create table bids_reviews
(
    id uuid primary key default uuid_generate_v4(),
    description text,
    created_at timestamp not null default now(),
    bid_id uuid references bids(id)
)