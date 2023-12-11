# BankZilla (Simple Bank Backend Service)

## Project Overview

### Features:
Create and manage bank accounts with owner's name, balance, and currency.
Record all balance changes, generating entry records for each transaction.
Perform secure money transfers between two accounts, ensuring transactional consistency.

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
