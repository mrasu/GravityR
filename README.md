GravityR is a Gravity-Radar.  
This exists to reduce time to find bottleneck in your application.
And also this is to solve the problems faster and easier.  

# Features

### Analyze application bottleneck

You can find bottleneck of application using Jaeger or AWS's PerformanceInsights

![Jaeger](/docs/images/jaeger.png)

### Analyze database bottleneck

You can find bottleneck of PostgreSQL, Hasura and MySQL

![PostgreSQL](/docs/images/postgres.png)

# Example

### Dig

Below shows performance information by digging run history like traces.  

* [app dig jaeger](https://mrasu.github.io/GravityR/jaeger.html)
* [db dig performance-insights](https://mrasu.github.io/GravityR/performance-insights.html)

### Suggest

Below shows how you can accelerate the query by adding index.

* [db suggest mysql](https://mrasu.github.io/GravityR/mysql.html)
* [db suggest mysql --with-examine](https://mrasu.github.io/GravityR/mysql_examine.html)
* [db suggest postgres](https://mrasu.github.io/GravityR/postgres.html)
* [db suggest postgres --with-examine](https://mrasu.github.io/GravityR/postgres_examine.html)
* [db suggest hasura postgres](https://mrasu.github.io/GravityR/hasura_postgres.html)
* [db suggest hasura postgres --with-examine](https://mrasu.github.io/GravityR/hasura_postgres_examine.html)

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
    export JAEGER_UI_URL="http://localhost:16686" JAEGER_GRPC_ADDRESS="localhost:16685"
    
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
   
    # Get application information with Jaeger
    gr app dig jaeger -o "example/jaeger.html" --service-name bff
    ```

# Environment Variables

| Command                    | Key                 | Description                                      |
|----------------------------|---------------------|--------------------------------------------------|
| db suggest hasura postgres | HASURA_URL          | URL to hasura. (e.g. http://localhost:8081)      |
|                            | HASURA_ADMIN_SECRET | Token to connect as admin                        |
| db suggest mysql           | DB_USERNAME         | User of mysql                                    |
|                            | DB_PASSWORD         | Password of mysql                                |
|                            | DB_PROTOCOL         | Protocol to connect mysql. `tcp` by default      |
|                            | DB_ADDRESS          | Network address of mysql                         |
|                            | DB_DATABASE         | Database name of mysql                           |
|                            | DB_PARAM_TEXT       | Arbitrary text used after `?` in DSN             |
| db suggest postgres        | DB_USERNAME         | User of PostgreSQL                               |
|                            | DB_PASSWORD         | Password of PostgreSQL                           |
|                            | DB_HOST             | Host of PostgreSQL                               |
|                            | DB_PORT             | Port of PostgreSQL                               |
|                            | DB_DATABASE         | Database name of PostgreSQL                      |
|                            | DB_SEARCH_PATH      | Search path of PostgreSQL. `public` by default   |
|                            | DB_PARAM_TEXT       | Arbitrary text for connection string             |
| app dig jaeger             | JAEGER_UI_URL       | URL to Jaeger's UI (e.g. http://localhost:16686) |
|                            | JAEGER_GRPC_ADDRESS | Address to GRPC's query server                   |

# Build
```sh
$ cd client && yarn
$ make build
```
