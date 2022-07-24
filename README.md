GravityR is a gravity radar. Find bottleneck in your application.

# Sample
GravityR's results for slow SQL because no index exists

* [suggest](https://htmlpreview.github.io/?https://github.com/mrasu/GravityR/blob/main/sample/output.html)
* [suggest --with-examine](https://htmlpreview.github.io/?https://github.com/mrasu/GravityR/blob/main/sample/output_examine.html)

# Usage
```sh
# Set evn to connect DB
export DB_USERNAME=root DB_DATABASE=gravityr

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
