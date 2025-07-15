# Worker RabbitMQ Documentation

This module is part of the `cmd/worker` directory and is responsible for consuming and processing messages from specific message queue topics. It runs background tasks (workers/consumers) that handle asynchronous event processing.

## Available Topics

The table below outlines the supported topics and their purposes:

| Constant Name      | Topic Name         | Description                                |
|--------------------|--------------------|--------------------------------------------|
| `ProcessSyncLog`   | `log.insert`       | Handles log synchronization insert events. |
| `ProcessExample`   | `example.consumer` | Example consumer for demonstration/testing.|


## Consumer Process

All consumer processing logic is located in: `internal/queue/consumer/`.

IMPORTANT: For more complex business logic, it is recommended to use the **Usecase Layer** to separate concerns and maintain clean architecture. This helps in improving code readability, testability, and scalability.



## How to Run

Make sure the required dependencies and environment configurations are properly set, then run the worker:

```bash
go run cmd/worker/main.go <topic.name>
```

## Notes

Each topic should have a corresponding handler that processes the message payload.

Add new topic constants and consumer logic to extend functionality.

For more topic handlers and implementation logic, refer to the file in the `internal/queue/consumer/` directory.