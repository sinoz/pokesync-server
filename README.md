# PokeSync - Game Service

## Building and Running the Application

First, make `start.sh` an executable:

```
chmod u+x start.sh
```

And now you can freely run it at the root directory of the project:

```
./start.sh
```

Alternatively, you can also just stick to calling `sh start.sh` but this may require `sudo`.

## Docker

### Running the service

```
docker-compose up
```

### Image storage

Images of this game service are stored on GitLab Container Registry:

```
docker login registry.gitlab.com
```

```
docker build -t registry.gitlab.com/pokesync/game-service .
```

```
docker push registry.gitlab.com/pokesync/game-service
```

## Running Tests

```
go test -v ../...
```

## How to Contribute

Please read [CONTRIBUTING.md](https://gitlab.com/pokesync/game-service/blob/master/CONTRIBUTING.md)