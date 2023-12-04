CREATE TABLE ORDERS(
    order_uid  VARCHAR(60) PRIMARY KEY,
    track_number  VARCHAR(30) NOT NULL,
    entry  VARCHAR(30) NOT NULL,
    locale  VARCHAR(10) NOT NULL,
    internal_signature  VARCHAR(30) NOT NULL,
    customer_id  VARCHAR(30) NOT NULL,
    delivery_service  VARCHAR(30) NOT NULL,
    shardkey  VARCHAR(30) NOT NULL,
    sm_id BIGINT NOT NULL,
    date_created TIMESTAMP NOT NULL,
    oof_shard  VARCHAR(30) NOT NULL
);


CREATE TABLE DELIVERY(
    order_id VARCHAR(60) PRIMARY KEY REFERENCES ORDERS(order_uid),
    name VARCHAR(60) NOT NULL,
    phone VARCHAR(30) NOT NULL,
    zip VARCHAR(30) NOT NULL,
    city VARCHAR(30) NOT NULL,
    address VARCHAR(60) NOT NULL,
    region VARCHAR(60) NOT NULL,
    email VARCHAR(60) NOT NULL
);


CREATE TABLE PAYMENT(
    order_id VARCHAR(60) PRIMARY KEY REFERENCES ORDERS (order_uid),
    transaction VARCHAR(60) NOT NULL,
    request_id VARCHAR(30) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    provider VARCHAR(30) NOT NULL,
    amount DECIMAL(22, 3) NOT NULL,
    payment_dt BIGINT NOT NULL,
    bank VARCHAR(60) NOT NULL,
    delivery_cost DECIMAL(22, 3) NOT NULL,
    goods_total INTEGER NOT NULL,
    custom_fee DECIMAL(22, 3) NOT NULL
);


CREATE TABLE ITEM(
    id SERIAL PRIMARY KEY,
    order_id  VARCHAR(60) REFERENCES ORDERS (order_uid),
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(60) NOT NULL,
    price DECIMAL(22, 3) NOT NULL,
    rid VARCHAR(60) NOT NULL,
    name VARCHAR(60) NOT NULL,
    sale INTEGER NOT NULL,
    size VARCHAR(30) NOT NULL,
    total_price DECIMAL(22, 3) NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(30) NOT NULL,
    status INTEGER NOT NULL
);
