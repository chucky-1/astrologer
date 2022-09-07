CREATE TABLE picture
(
    title varchar(100),
    date  date,
    image bytea,

    PRIMARY KEY (title, date)
)
