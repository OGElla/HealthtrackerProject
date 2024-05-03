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
```dsn``` — postgress connection string with username, password, address, port, database name, and SSL mode. Default: Value is not correct by security reasons.

```migrations``` — Path to folder with migration files. If not provided, migrations do not applied.

```fill``` — Fill database with dummy data. Default: false.

```env``` - App running mode. Default: development

```port``` - App port. Default: 8081
### Run with docker-compose

## Healthtracker REST API 
```
*list of all trackers*
GET /health/daily  
POST /health/daily  
GET /health/view/:id  
PUT /health/view/:id  
DELETE /health/view/:id  
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

Table goals{
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
