# MeetUp Clone

A backend API for a meetup-clone

## Getting Started

- Install dependencies

```sh
go mod tidy
```

- Create a `.env` directory and add `api.env` and `postgres.env` files

```sh
mkdir .env
touch .env/api.env
touch .env/postgres.env # Only needed if you are using `docker-compose.dev.yml`
```

- Add environment variables to the `.env/*.env` files

  - .env/api.env

    ```sh
    # App environment variables
    export DEBUG=true
    export FRONTEND_URL=http://localhost:3000

    # Server environment variables
    export SERVER_ADDR=:8000
    export SERVER_ENVIRONMENT=dev # dev, test, prod

    # Database environment variables
    export DB_ADDR=postgres://<user>:<password>@<host>:<port>/<dbName>?sslmode=disable
    export DB_MAX_OPEN_CONNS=30
    export DB_MAX_IDLE_CONNS=30
    export DB_MAX_IDLE_TIME=15m

    # JWT environment variables
    export JWT_ISSUER=meetup_clone
    export JWT_AUDIENCE=meetup_clone
    export JWT_SECRET_KEY=<secret-key>
    export JWT_ACCESS_EXP=3

    # Test environment variables
    export TEST_DB_ADDR=postgres://<user>:<password>@<host>:<port>/<dbName>_test?sslmode=disable
    ```

  - .env/postgres.env - Only needed if you are using `docker-compose.dev.yml`

    ```sh
    POSTGRES_PORT=5432
    POSTGRES_USER=postgres
    POSTGRES_PASSWORD=postgres
    POSTGRES_MULTIPLE_DATABASES=meetup, meetup_test
    ```

- Start the api services i.e postgres

    **NOTE**: This step is for only those using the `docker` option. If you have local instances, skip this step.

```sh
make services-up
# OR if you need sudo
sudo make services-up

# Stop the services
make services-down

# Stop services and destroy container
make services-kill
```

- Run migrations on the test and main db

    **NOTE**: You must have [migrate](https://github.com/golang-migrate/migrate) installed.

```sh
make migrate-up
make test-migrate-up
```

- Run tests

```sh
make test
```

- Run server

```sh
make runserver
```
