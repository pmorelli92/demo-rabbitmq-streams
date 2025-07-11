# RabbitMQ Streams Demo with Golang

## 🎯 Purpose

This repository demonstrates **event-driven microservices architecture** using **RabbitMQ Streams** for reliable event streaming and **eventual consistency** between services. It showcases how services can maintain their own data while staying synchronized through events, providing a production-ready pattern for distributed systems.

## 🏗️ Architecture

![Architecture Overview](./diagrams/architecture-overview.svg)

![Event Flow](./diagrams/event-flow.svg)

![Single Active Consumer Pattern](./diagrams/single-active-consumer.svg)

## 🚀 Features

- Source of truth in streams.
- Automatic synchronization through events.
- Avoid duplicate processing.
- Automatic failover for consumers.
- Reliable producer and consumers for network intermittency or downtime with RabbitMQ.

## 🧪 Testing Suite

You need to have both `docker` and `go` installed to run the demo.

```bash
make dependencies
make run-e2e
make run-load
```
