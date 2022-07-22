GravityR is a gravity radar. Find bottleneck in your application.

# Usage
```sh
# Search your SQL with EXPLAIN
gr db suggest -o "sample/output.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your SQL by adding indexes temporarily
gr db suggest --with-examine -o "sample/output_examine.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"
```

# Build
```sh
$ cd client && yarn
$ make build
```
