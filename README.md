# HealthtrackerProjectGo
Welcome to HealthtrackerProjectGo, a powerful health tracking application built using the Go programming language. This project aims to provide users with a comprehensive platform to monitor and manage various aspects of their health and well-being seamlessly.
## What is HealthtrackerProjectGo?
HealthtrackerProjectGo is designed to simplify the process of tracking health-related data, including vital signs, exercise routines, dietary habits, and more. Leveraging the efficiency and concurrency features of Go, this application offers robust performance and scalability, ensuring smooth operation even with large volumes of data.
## Getting Started
### Using app golang directly on Terminal
Provide all needed correct values by using flags. The default command line:
```
go run ./cmd/api
-dsn="postgres://healthtracker:12345@localhost/healthtracker?sslmode=disable"
-migrations=file://migrations
-fill=false
-env=development
-port=4001
```
#### List of flags
```dsn``` — postgres connection string with username, password, address, port, database name, and SSL mode. Default: Value is not correct by security reasons.

```migrations``` — Path to the folder with migration files. If not provided, migrations do not apply.

```fill``` — Fill the database with dummy data. Default: false.

```env``` - App running mode. Default: development

```port``` - App port. Default: 4001


### Run with docker-compose
```
env POSTGRES_PASSWORD="STRONG PASSWORD" APP_DSN="postgres://postgres:postgres@db:5432/healthtracker?sslmode=disable" docker-compose --env-file .env.example up --build
```

```env POSTGRES_PASSWORD="postgres"``` This command adds the environment variable then available in docker-compose.

```APP_DSN``` contains the connection string to the dockerized Postgres.

Overall, your DSN for docker should be like this: postgres://postgres:postgres@db:5432/healthtracker?sslmode=disable.

```--build``` flag forces docker-compose to rebuild the app. For example, if you have changed the source code, you need this flag.

## Healthtracker REST API
```
*list of all trackers*
GET /health/daily  
POST /health/daily  
GET /health/view/:id  
PUT /health/view/:id  
DELETE /health/view/:id  
```
## API Reference
**Healthcheck**

``` 
GET /
 ```

**Registration**

```
POST /users
```

```json
{
    "email": "yourEmail",
    "name": "yourName",
    "password": "yourPassword"
}
```

**Activation**

```
PUT /users/activated
```

```json
{
    "token": "receivedToken",
}
```

**Login**

```
POST /tokens/authentication
```

```json
{
    "email": "{{email}}",
    "password": "{{password}}"
}
```

**List your trackers**

```
GET /health/daily
```

Parameters

```
**Sorting**

sort=$1, where $1 can be a positive or negative value, depends on sorting in ascending or descending order

**Filtering**

calories=$1&walking=$2&hydrate=$3&sleep=$4, where $1, $2, $3, $4 are the tracker values

**Pagination**

page=$1&age_size=$2, where $1 - the first page in dataset and #$2 = the last page in dataset
```

**Create a new tracker** *token is required!*

```
POST /health/daily
```

```json
{  
     "calories":"3000 calories",
     "walking":"10000 steps",
     "hydrate":"2 liters",
     "sleep":"8 hours"
}
```

**Create a new goal** *token is required!*

```
POST /goals/daily
```

```json
{  
     "walking":"10000 steps",
}
```

## DB Structure
```
Table healthtracker {
    id bigserial [primary key]
    created_at timestamp
    walking text
    hydrate text
    sleep text
    user_id integer
    version integer
}


Table goals {
    id bigserial PRIMARY KEY,
    created_at timestamp,
    walking integer,
    achieved bool,
    user_id integer, 
    version integer
}

// many-to-many
Table healthtracker_and_goals {
  id bigserial [primary key]
  created_at timestamp
  healthtracker bigserial
  goal bigserial
}

Ref: healthtracker_and_goals.healthtracker < healthtracker.id
Ref: healthtracker_and_goals.goal < goals.id
```

## Contributing 
Contributions to HealthtrackerProjectGo are welcome! Whether you're interested in adding new features, fixing bugs, or improving documentation, your contributions help make this project better for everyone.

## Authors
Dussekenov Elnar 21B030333 @depayka
