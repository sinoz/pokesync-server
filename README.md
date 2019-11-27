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

Otherwise, you can stick to `sh start.sh`

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

Rapid development with docker (compose) isn't supported yet.

## How to Contribute

Please read [CONTRIBUTING.md](https://gitlab.com/pokesync/game-service/blob/master/CONTRIBUTING.md)