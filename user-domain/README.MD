internal/
├── application/
├── domain/
└── infrastructure/
    ├── http/
    │   ├── router.go
    │   ├── middleware/
    │   │   ├── request_id.go
    │   │   └── auth.go
    │   └── handler/
    │       └── user_handler.go
    │
    ├── logger/
    │   └── zap_logger.go
    │
    ├── persistence/
    │   ├── postgres/
    │   │   ├── user_repository.go
    │   ├── redis/
    │   │   └── cache_adapter.go
    │   └── dao/
    │       └── dao.go
    │
    ├── messaging/
    │   ├── kafka_producer.go
    │   └── kafka_consumer.go
    │
    ├── config/
    │   └── config_loader.go
    │
    ├── tracing/
    │   └── opentelemetry.go
    │
    └── metrics/
        └── prometheus.go
