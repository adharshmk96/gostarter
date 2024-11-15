# gostarter

## First Run

```bash
make maiden
```

## About

Golang starter project with simple webserver, migrations and clean architecture.

- [Cobra CLI](https://github.com/spf13/cobra)
- [Chi Router](https://github.com/go-chi/chi)
- [Swaggo](https://github.com/swaggo/swag)
- [PGX](https://github.com/jackc/pgx)

## Structure

- infra: Code related to infrastructure, like logging, database, o11y, config, etc.
- internal: Business logic, Storage logic, Http Handlers, etc.
- platform: Configs for services like nginx, db and platform tools.
- pkg: commonly used packages across the project.
- cmd: Entrypoint for the commands in the project, server, migrations, etc.

Services like message brokers are caches should be in infra folder. It should be initialized and passed down via cmd
package.

Note: These services should be global in nature

Services like External API Client should be initialized in pkg folder. It should be initialized from delivery layer and
passed down to the services that need it.

Note: These services are more specific, not necessarily global.

Internal

- domain: Business logic, entities, value objects, etc.
- service: Business logic, use cases, etc.
- storage: Database logic, migrations, etc.
- delivery: Http, Worker, etc.
    - http: Http handlers, middlewares, etc.
        - server: Server Framework, Router handlers, etc.
    - worker: Worker handlers, etc.
