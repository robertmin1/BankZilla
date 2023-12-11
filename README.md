# BankZilla (Simple Bank Backend Service)

## Project Overview (On-going project)

### Features:
- Create and manage bank accounts with owner's name, balance, and currency.
- Record all balance changes, generating entry records for each transaction.
- Perform secure money transfers between two accounts, ensuring transactional consistency.

## Key Highlights
## Database Design and Management
### Database Design
Designed a robust database schema.
Implemented SQL code generation for seamless database interaction.

### Database Migration and Testing:
Executed database migration using Golang, ensuring consistency and reliability.
Conducted extensive unit tests for database CRUD operations.

### RESTful API Development
Constructed a set of RESTful HTTP APIs using Gin.
Covered various aspects, including config loading, unit testing, error handling, and user authentication with JWT and PASETO tokens.

## Production Deployment
### Docker and Kubernetes Deployment:
Built a minimal Golang Docker image and deployed the service to a production Kubernetes cluster on AWS.
Set up HTTPS, automatic TLS certificate renewal with Let's Encrypt, and domain registration with Route53.

## Advanced Backend Topics
Implemented advanced features such as managing user sessions, gRPC APIs, and Swagger documentation (Not finished yet)
Explored partial record updates, structured logging, and role-based access control (RBAC).

## Asynchronous Processing with Background Workers
### Asynchronous Processing:
Integrated background workers with Redis for asynchronous processing.
Covered email integration and unit testing for asynchronous tasks.

## Improving Stability and Security
### Continuous Improvement:
Covered topics like updating dependency packages, improving refresh token security, and implementing graceful server shutdown.


## Setup local development

### Install tools

- [Docker desktop](https://www.docker.com/products/docker-desktop)
- [TablePlus](https://tableplus.com/)
- [Golang](https://golang.org/)
- [Homebrew](https://brew.sh/)
- [Migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

    ```bash
    brew install golang-migrate
    ```

- [DB Docs](https://dbdocs.io/docs)

    ```bash
    npm install -g dbdocs
    dbdocs login
    ```

- [DBML CLI](https://www.dbml.org/cli/#installation)

    ```bash
    npm install -g @dbml/cli
    dbml2sql --version
    ```

- [Sqlc](https://github.com/kyleconroy/sqlc#installation)

    ```bash
    brew install sqlc
    ```

- [Gomock](https://github.com/golang/mock)

    ``` bash
    go install github.com/golang/mock/mockgen@v1.6.0
    ```

### Setup infrastructure

- Create the bank-network

    ``` bash
    make network
    ```

- Start postgres container:

    ```bash
    make postgres
    ```

- Create simple_bank database:

    ```bash
    make createdb
    ```

- Run db migration up all versions:

    ```bash
    make migrateup
    ```

- Run db migration up 1 version:

    ```bash
    make migrateup1
    ```

- Run db migration down all versions:

    ```bash
    make migratedown
    ```

- Run db migration down 1 version:

    ```bash
    make migratedown1
    ```

### Documentation

- Generate DB documentation:

    ```bash
    make db_docs
    ```

- Access the DB documentation at [this address](https://dbdocs.io/techschool.guru/simple_bank). Password: `secret`

### How to generate code

- Generate schema SQL file with DBML:

    ```bash
    make db_schema
    ```

- Generate SQL CRUD with sqlc:

    ```bash
    make sqlc
    ```

- Generate DB mock with gomock:

    ```bash
    make mock
    ```

- Create a new db migration:

    ```bash
    make new_migration name=<migration_name>
    ```

### How to run

- Run server:

    ```bash
    make server
    ```

- Run test:

    ```bash
    make test
    ```

## Deploy to kubernetes cluster

- [Install nginx ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/#aws):

    ```bash
    kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v0.48.1/deploy/static/provider/aws/deploy.yaml
    ```

- [Install cert-manager](https://cert-manager.io/docs/installation/kubernetes/):

    ```bash
    kubectl apply -f https://github.com/jetstack/cert-manager/releases/download/v1.4.0/cert-manager.yaml
    ```
