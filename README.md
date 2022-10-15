GravityR is a Gravity-Radar.  
This exists to reduce time to find bottleneck in your application.
And also this is to solve the problems faster and easier.  

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
* [suggest hasura postgres](https://mrasu.github.io/GravityR/hasura_postgres.html)
* [suggest hasura postgres --with-examine](https://mrasu.github.io/GravityR/hasura_postgres_examine.html)

# Usage

1. Download the latest binary from [releases](https://github.com/mrasu/GravityR/releases) page
2. Give execution permission and rename to `gr`
    ```sh
    chmod u+x gr_xxx
    mv gr_xxx gr
    ```
3. Run commands to find bottlenecks!  
    examples:
    ```sh
    # Set envs to connect DB and Hasura
    export DB_USERNAME=root DB_DATABASE=gravityr
    export HASURA_URL="http://localhost:8081" HASURA_ADMIN_SECRET="myadminsecretkey" 
    
    # Get database information with AWS' Performance Insights
    gr db dig performance-insights -o "example/performance-insights.html"
    
    # Search your MySQL's SQL with EXPLAIN
    gr db suggest mysql -o "example/mysql.html" -q "SELECT name, t.description FROM users INNER JOIN tasks AS t ON users.id = t.user_id WHERE users.name = 'foo'"
    
    # Search your MySQL's SQL by adding indexes temporarily
    gr db suggest mysql --with-examine -o "example/mysql_examine.html" -q "SELECT name, t.description FROM users INNER JOIN tasks AS t ON users.id = t.user_id WHERE users.name = 'foo'"
    
    # Search your PostgreSQL's SQL with EXPLAIN
    gr db suggest postgres -o "example/postgres.html" -q "SELECT name, t.description FROM users INNER JOIN tasks AS t ON users.id = t.user_id WHERE users.name = 'foo'"
    
    # Search your PostgreSQL's SQL by adding indexes temporarily
    gr db suggest postgres --with-examine -o "example/postgres_examine.html" -q "SELECT name, t.description FROM users INNER JOIN tasks AS t ON users.id = t.user_id WHERE users.name = 'foo'"
    
    # Search your Hasura's GraphQL with EXPLAIN
    gr db suggest hasura postgres -o "example/hasura.html" -q "query MyQuery(\$email: String) {
      tasks(where: {user: {email: {_eq: \$email}}}) {
        user {
          name
        }
        description
      }
    }
    " --json-variables '{"email": "test1112@example.com"}'
    
    # Search your Hasura's GraphQL by adding indexes temporarily
    gr db suggest hasura postgres --with-examine -o "example/hasura_examine.html" -q "query MyQuery(\$email: String) {
      tasks(where: {user: {email: {_eq: \$email}}}) {
        user {
          name
        }
        description
      }
    }
    " --json-variables '{"email": "test1112@example.com"}'
    ```

# Build
```sh
$ cd client && yarn
$ make build
```
