# Architecture

Backend is organized as a small set of services connected through Docker Compose.

The API Gateway exposes the HTTP API and delegates work to internal gRPC services. Domain services own their own storage and expose contracts through protobuf definitions from the contracts module.

## Runtime Components

- API Gateway receives HTTP requests.
- Authorization service manages users, login, tokens, and sessions.
- Feed service stores and returns car listings.
- Prediction service estimates prices and returns model-related data.
- PostgreSQL instances persist service-owned data.
- Redis stores authorization sessions.

## Communication

External clients talk to the gateway over HTTP. Internal service-to-service calls use gRPC over the compose network.
