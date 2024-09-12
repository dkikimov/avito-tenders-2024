create table bids_approvals(
    bid_id uuid not null references bids(id),
    user_id uuid not null references employee(id),
    primary key (bid_id, user_id)
)