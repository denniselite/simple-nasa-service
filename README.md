# simple-nasa-service

=======

Simple NASA service aggregates information about Near Earth Objects (NEOs) and helps to analyze data about them.

Installation
------------

```sh
make
make install
```

Configuration
-------------

Format:

```yml

listen: 127.0.0.1:3001

pgsql:
  host: 127.0.0.1
  port: 5432
  username: postgres
  password: postgres
  database: nasa

NASAServer:
  API-KEY: N7LkblDsc5aen05FJqBQ8wU4qSdmsftwJagVK7UD
  endPoint: https://api.nasa.gov/neo/rest
  APIVersion: /v1

```

Starting service
----------------

```sh
bin/simple-nasa-service --c cfg/config.yml
```


