# Go healthcheck
Go healthcheck is an opensource healthcheck system that ensures your HTTP applications are up and running.

## Run
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
  mysql:
    image: mysql
    environment:
      - MYSQL_ROOT_PASSWORD=password
      - MYSQL_DATABASE=healthcheck
```

Then run :
```
    docker-compose up
```

Go to http://localhost:8080/app to see your dashboard