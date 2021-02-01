CREATE TABLE users (
  username   varchar,
  email      varchar,
  password   varchar
);

CREATE TABLE balances (
  username    varchar,
  balance     int
);

INSERT INTO users (username, email, password) VALUES
('test1', 'test1@test.com', 'test1pass'),
('test2', 'tes21@test.com', 'test2pass');

INSERT INTO balances (username, balance) VALUES
('test1', 100),
('test2', 200);

ALTER TABLE users
  ADD PRIMARY KEY (username);


ALTER TABLE balances
  ADD PRIMARY KEY (username);

COMMIT;
