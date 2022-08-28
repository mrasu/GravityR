CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(255) NOT NULL,
  created_at timestamp with time zone DEFAULT NULL,
  updated_at timestamp with time zone DEFAULT NULL
);

CREATE TABLE todos (
  id SERIAL PRIMARY KEY,
  user_id int NOT NULL,
  title varchar(50) NOT NULL,
  description varchar(100) DEFAULT NULL,
  status smallint NOT NULL,
  created_at timestamp with time zone DEFAULT NULL,
  updated_at timestamp with time zone DEFAULT NULL
);

/* Insert 256*256(=65,536) rows */
WITH l0 AS (SELECT 0 AS num UNION ALL SELECT 1),
     l1 AS (SELECT a.num * 2 + b.num AS num FROM l0 AS a CROSS JOIN l0 AS b),
     l2 AS (SELECT a.num * 4 + b.num AS num FROM l1 AS a CROSS JOIN l1 AS b),
     l3 AS (SELECT a.num * 16 + b.num AS num FROM l2 AS a CROSS JOIN l2 AS b),
     l4 AS (SELECT a.num * 256 + b.num AS num FROM l3 AS a CROSS JOIN l3 AS b)
INSERT INTO users(name, email, password, created_at, updated_at)
SELECT CONCAT('test', num),
       CONCAT('test', num, '@example.com'),
       'p@ssw0rd',
       CURRENT_TIMESTAMP - INTERVAL '1 DAY' * ((num%363) + 1),
       CURRENT_TIMESTAMP - INTERVAL '1 DAY' * (num%363)
FROM l4;

/* Insert 65,536*100(=6,553,600) rows */
INSERT INTO todos(user_id, title, description, status, created_at, updated_at)
SELECT users.id,
     'test title',
     'test description',
     tmp_100.id % 3,
     users.updated_at - INTERVAL '1 HOUR' * ((tmp_100.id%100) +5),
     users.updated_at - INTERVAL '1 HOUR' * (tmp_100.id%100)
FROM users
     CROSS JOIN (SELECT * FROM users LIMIT 100) tmp_100;
