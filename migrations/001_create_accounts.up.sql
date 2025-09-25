

create table if not exists accounts (
    id bigserial primary key,
    user_id bigint not null,
    name varchar(255) not null,
    display_name varchar(255) not null,
    timezone varchar(100) default 'Asia/Almaty',
    is_active boolean default true,
    created_at timestamp with time zone default now (),
    updated_at timestamp with time zone default now (),


        constraint unique_user_id unique (user_id)
    );


CREATE INDEX idx_accounts_user_id ON accounts(user_id);
CREATE INDEX idx_accounts_is_active ON accounts(is_active);