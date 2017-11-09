# Go healthcheck
Go healthcheck is an opensource healthcheck system that ensures your HTTP applications are up and running.

## Run
Create the config.yml file for app :
```
history:
  enabled: true
db:
  type: sqlite3
  path: data.db
```

Typical docker-compose.yml file :

```
version: '3.1'
services:
  go-healthcheck:
    build: .
    image: ghiltoniel/go-healthcheck
    ports:
      - 8080:8080
    depends_on:
      - mysql
    volumes:
      - ./config_docker.yml:/go/src/app/config.yml
```

Then run :
```
    docker-compose up
```

Go to http://localhost:8080/app to see your dashboard

## Configuration
- history:
  - enabled: true if you want the check history to be saved into db, false if you want to keep only latest check
- db:
  - type: mysql / sqlite3 / postgres
  - username: the database username
  - password: the database password
  - host: the database host
  - port: the database port
  - name: the database name
  - path: the file database path (for sqlite3)
