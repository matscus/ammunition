# Ammunition, a simple application, for test datapools and universal KV in-memory store 
![GitHub](https://img.shields.io/github/license/matscus/ammunition?color=31E311)

#### Attention - The previous version was moved to the v1 branch and marked as deprecated.

- Long term storage in postgresql.
- In-memory stored all data.
- Persist cache not supporteg get current value, only iterator or random.
- Returns unique for each pool one by one persisted pool.
- Prepared statement data, for minimum response time.
- Returns unique or current values ​​for a temporary pool.
- Convert upload csv file to json struct.
- Support for generic data structures in the request body for temporary cache - store all body from request , then [] byte.


### Quick start 

```sh
    #Clone repos
    git clone https://github.com/matscus/ammunition.git && cd ammunition

    #Build and running an application in a container 
    make run
```

- Swagger:  http://localhost:9443/swagger/index.html#/

### Stop

```sh
    #Stop application and postgres 
    make stop
```

### Stop and delete

```sh
    #Stop application with drop containers, volumes and images
    make kill
```

### Build binary
```sh
    #Build binary file
    make engine
```

## Docker image build
```
docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t <your registry>/ammunition --push .
```

### TODO
- add description api
- add all tests
