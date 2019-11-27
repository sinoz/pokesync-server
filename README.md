# PokeSync - Game Service

## Building and Running the Application

From the root directory of the project:

```
sh start.sh
```

To make `start.sh` executable:

```
chmod u+x start.sh
```

### With Docker

```
docker build -t game-service .
docker run -t game-service .
```

Docker-Compose is also supported:

```
docker-compose up
docker-compose down
```

## How to Contribute

Please read [CONTRIBUTING.md](https://gitlab.com/pokesync/game-service/blob/master/CONTRIBUTING.md)