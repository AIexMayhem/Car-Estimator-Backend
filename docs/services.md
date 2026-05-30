# Services

## API Gateway

Entry point for external HTTP clients. It validates request shape, maps HTTP payloads to service requests, and forwards calls to gRPC services.

## Authorization

Handles registration, login, token generation, token validation, refresh flow, and session persistence.

## Feed

Owns listing and favorite-related data. It exposes listing search, listing details, creation, update, deletion, and favorite operations.

## Prediction

Runs price estimation logic and serves image data used by the prediction workflow.

## API Contracts

Stores protobuf definitions and generated code shared by services. Contract changes should be coordinated with consumers before service code is updated.
