# HealthtrackerProjectGo

## Healthtracker REST API 
```
*list of all trackers*
GET /health/daily  
POST /health/daily  
GET /health/view/:id  
PUT /health/view/:id  
DELETE /health/view/:id  
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
