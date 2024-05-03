# HealthtrackerProjectGo
Welcome to HealthtrackerProjectGo, a powerful health tracking application built using the Go programming language. This project aims to provide users with a comprehensive platform to monitor and manage various aspects of their health and well-being seamlessly.
## Getting Started
### Using app golang directly on Terminal
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
