GravityR is a gravity radar. Find bottleneck in your application.

# Example
GravityR's results for slow SQL because no index exists

* [dig performance-insights](https://htmlpreview.github.io/?https://github.com/mrasu/GravityR/blob/main/example/performance-insights.html)
* [suggest](https://htmlpreview.github.io/?https://github.com/mrasu/GravityR/blob/main/example/output.html)
* [suggest --with-examine](https://htmlpreview.github.io/?https://github.com/mrasu/GravityR/blob/main/example/output_examine.html)

# Usage
```sh
# Set evn to connect DB
export DB_USERNAME=root DB_DATABASE=gravityr

# Get info with AWS' Performance Insights
gr db suggest -o "example/output.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your SQL with EXPLAIN
gr db suggest -o "example/output.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your SQL by adding indexes temporarily
gr db suggest --with-examine -o "example/output_examine.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"
```

# Build
```sh
$ cd client && yarn
$ make build
```
