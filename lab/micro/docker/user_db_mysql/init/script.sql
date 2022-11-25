CREATE TABLE users (
  id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL
);

CREATE TABLE payment_methods (
  id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  last_numbers varchar(10) NOT NULL
);

CREATE TABLE payment_histories (
  id int PRIMARY KEY NOT NULL AUTO_INCREMENT,
  user_id int NOT NULL,
  payment_method_id int NOT NULL,
  amount int NOT NULL,
  payed_at datetime NOT NULL
);

/* Insert 16 rows */
INSERT INTO users(id, name, email, password)
WITH l0 AS (SELECT 0 AS num UNION ALL SELECT 1),
     l1 AS (SELECT a.num * 2 + b.num AS num FROM l0 AS a CROSS JOIN l0 AS b),
     l2 AS (SELECT a.num * 4 + b.num AS num FROM l1 AS a CROSS JOIN l1 AS b)
SELECT num+1,
       CONCAT('test', CAST(num AS char)),
       CONCAT('test', CAST(num AS char), '@example.com'),
       'p@ssw0rd'
FROM l2;

INSERT INTO payment_methods(id, user_id, last_numbers)
SELECT
  id,
  id,
  id%10*1111
FROM users;

/* Insert 16*16 rows */
INSERT INTO payment_histories(id, user_id, payment_method_id, amount, payed_at)
WITH l0 AS (SELECT 0 AS num UNION ALL SELECT 1),
     l1 AS (SELECT a.num * 2 + b.num AS num FROM l0 AS a CROSS JOIN l0 AS b),
     l2 AS (SELECT a.num * 4 + b.num AS num FROM l1 AS a CROSS JOIN l1 AS b)
SELECT us.id*16+num,
       us.id,
       us.id,
       (us.id*16+num) * 100,
       NOW() - INTERVAL (num%16) MONTH
FROM l2 CROSS JOIN (SELECT id FROM users) AS us;
