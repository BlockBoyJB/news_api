create table if not exists news
(
    id      bigserial primary key,
    title   varchar not null,
    content text    not null
);

create table if not exists news_categories
(
    id          bigserial primary key,
    news_id     bigint references News (id) not null,
    category_id bigint                      not null
);
