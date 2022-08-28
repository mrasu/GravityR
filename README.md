GravityR is a Gravity-Radar.  
This exists to remove bottleneck in your application without help of experts.  
And also this is to help experts solving problems faster and easily.  

# Example

### Dig

Below shows slow queries by digging database's history from AWS's PerformanceInsights.  

* [dig performance-insights](https://mrasu.github.io/GravityR/performance-insights.html)

### Suggest

Below shows how you can accelerate the query by adding index.

* [suggest mysql](https://mrasu.github.io/GravityR/mysql.html)
* [suggest mysql --with-examine](https://mrasu.github.io/GravityR/mysql_examine.html)
* [suggest postgres](https://mrasu.github.io/GravityR/postgres.html)
* [suggest postgres --with-examine](https://mrasu.github.io/GravityR/postgres_examine.html)

# Usage
```sh
# Set envs to connect DB
export DB_USERNAME=root DB_DATABASE=gravityr

# Get database information with AWS' Performance Insights
gr db dig performance-insights -o "example/performance-insights.html"

# Search your MySQL's SQL with EXPLAIN
gr db suggest mysql -o "example/mysql.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your MySQL's SQL by adding indexes temporarily
gr db suggest mysql --with-examine -o "example/mysql_examine.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your PostgreSQL's SQL with EXPLAIN
gr db suggest postgres -o "example/postgres.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"

# Search your PostgreSQL's SQL by adding indexes temporarily
gr db suggest postgres --with-examine -o "example/postgres_examine.html" -q "SELECT name, t.description FROM users INNER JOIN todos AS t ON users.id = t.user_id WHERE users.name = 'foo'"
```

# Build
```sh
$ cd client && yarn
$ make build
```
