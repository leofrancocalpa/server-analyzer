# server-analyzer

last update in [mods](https://github.com/leofrancocalpa/server-analyzer/tree/mods) branch

## Steps to run the server side

1. Start CockroachDB and open SQL Shell and run queries below as root user

>cockroach start --insecure

>cockroach sql --insecure

``` SQL
CREATE USER <username>;
``` 
``` SQL
CREATE DATABASE challenge;
``` 
``` SQL
SET DATABASE = challenge;
``` 
``` SQL
GRANT ALL ON DATABASE challenge TO <username>;
```

``` SQL
CREATE TABLE serverinfo ( 
    id UUID DEFAULT uuid_v4()::UUID PRIMARY KEY, 
    dns STRING, 
    last_updated TIMESTAMP DEFAULT now(), 
    data JSONB 
);
```
2. Now run the main.go file
``` shell
go run cmd/main.go
```

To run the frontend go to [server-analyzer-gui](https://github.com/leofrancocalpa/server-analyzer-gui)
