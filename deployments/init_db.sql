CREATE TABLE payment (
     payment_id SERIAL PRIMARY KEY,
     transaction VARCHAR,
     request_id VARCHAR,
     currency VARCHAR,
     provider VARCHAR,
     payment_dt INT,
     bank VARCHAR,
     delivery_cost INT,
     goods_total INT,
     custom_fee INT
);

CREATE TABLE delivery (
          delivery_id SERIAL PRIMARY KEY,
          name VARCHAR,
          phone VARCHAR,
          zip VARCHAR,
          city VARCHAR,
          address VARCHAR,
          region VARCHAR,
          email VARCHAR
);

CREATE TABLE orders (
     order_uid VARCHAR PRIMARY KEY,
     track_number VARCHAR,
     entry VARCHAR,
     locale VARCHAR,
     internal_signature VARCHAR,
     customer_id VARCHAR,
     delivery_service VARCHAR,
     shardkey VARCHAR,
     sm_id INT,
     date_created TIMESTAMP,
     oof_shard VARCHAR,
     delivery INT REFERENCES delivery(delivery_id),
     payment INT REFERENCES payment(payment_id));

CREATE TABLE item (
      item_id SERIAL PRIMARY KEY,
      order_uid VARCHAR REFERENCES orders(order_uid),
      chrt_id INT,
      track_number VARCHAR,
      price INT,
      rid VARCHAR,
      name VARCHAR,
      sale INT,
      size VARCHAR,
      total_price INT,
      nm_id INT,
      brand VARCHAR,
      status INT
);