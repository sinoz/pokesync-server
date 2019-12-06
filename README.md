# PokeSync - Game Service

## Building and Running the Application

First, make sure you have a Redis instance running:

```
docker run --name some-redis -d -p 6379:6379 redis
```

Now that Redis is fired up, you can freely run the game service:

```
./start.sh
```

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