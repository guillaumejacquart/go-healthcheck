# Go healthcheck
Go healthcheck is an opensource healthcheck system that ensures your HTTP applications are up and running.

![Travis CI](https://travis-ci.org/guillaumejacquart/go-healthcheck.svg?branch=master) [![codecov](https://codecov.io/gh/guillaumejacquart/go-healthcheck/branch/master/graph/badge.svg)](https://codecov.io/gh/guillaumejacquart/go-healthcheck)

Live demo : https://check.apps.guillaumejacquart.com/app/

## Run
Create the docker-compose.yml file :

```
version: '3.1'
services:
  go-healthcheck:
    build: .
    image: ghiltoniel/go-healthcheck
    ports:
      - 8080:8080
    volumes:
      - ./config_docker.yml:/go/src/app/config.yml      
    environment:
      - DB.TYPE=sqlite3
      - DB.PATH=data.db
```

Then run :
```
    docker-compose up
```

Go to http://localhost:8080/app to see your dashboard

## Configuration
The configuration can be set in any of the following places :
- config.yml file at the root of the source
- config.yml file inside %HOME%/.go-healthcheck/
- config.yml file in /etc/go-healthcheck/
- in the environment variables (using capitalize letters, ex : DB.TYPE=sqlite3)

### Configuration variables
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

## TBD

- Notification email
- Notification pushbullet or else
- Test coverage
