# Hopper 🐇

*A transparent message broker built for simplicity and delivery guarantees*

[![Go Version](https://img.shields.io/badge/go-1.23.4-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

---

## 🚀 Why Hopper?

Modern microservices rely on message brokers to connect services reliably. While existing solutions are powerful, they often come with operational overhead:

- **RabbitMQ** → Reliable queuing, but requires separate management dashboards to understand queue health and message flow
- **Kafka** → Excellent for high-throughput event streaming, but complex to operate and optimized for stream replay rather than guaranteed delivery
- **Redis/SQS** → Simple but limited observability into message processing

**Hopper focuses on transparency and delivery guarantees.** It's designed so you can see exactly how messages flow through your system, with built-in observability that integrates seamlessly with monitoring tools like [Streamly](https://github.com/hoppermq/streamly).

---

## ✨ Key Features

- **🔍 Transparent delivery** → See exactly how messages move through queues and containers
- **✅ Delivery guarantees** → Reliable message processing without stream replay complexity  
- **🎯 Simple by design** → No complex clustering or configuration files to get started
- **📊 Built-in observability** → Native event bus and monitoring hooks
- **⚡ Lightweight & fast** → Written in Go with TCP transport layer
- **🎛️ Web UI included** → Built-in dashboard for queue management and monitoring

---

## 📦 Quick Start

### 1. Run Hopper Server

```bash
# Clone and run
git clone https://github.com/hoppermq/hopper
cd hopper
go run main.go
```

The server starts on default ports:
- **Message broker**: `localhost:5672` (TCP)
- **Web dashboard**: `localhost:8080` (HTTP)
- **Management API**: `localhost:9090` (HTTP)

### 2. Producer Service Example

```go
// order-service/main.go - Publishes order events
package main

import (
    "context"
    "encoding/json"
    "github.com/hoppermq/hopper/pkg/client"
)

type OrderEvent struct {
    OrderID   string `json:"order_id"`
    UserID    string `json:"user_id"`
    Status    string `json:"status"`
    Amount    float64 `json:"amount"`
}

func main() {
    producer := client.NewClient()
    
    ctx := context.Background()
    if err := producer.Run(ctx); err != nil {
        panic(err)
    }
    defer producer.Stop(ctx)
    
    // Publish order created event
    event := OrderEvent{
        OrderID: "order-123",
        UserID:  "user-456", 
        Status:  "created",
        Amount:  99.99,
    }
    
    data, _ := json.Marshal(event)
    // producer.Publish("orders.created", data) // Coming soon
}
```

### 3. Consumer Service Example

```go
// notification-service/main.go - Consumes order events
package main

import (
    "context"
    "encoding/json"
    "log"
    "github.com/hoppermq/hopper/pkg/client"
)

type OrderEvent struct {
    OrderID   string `json:"order_id"`
    UserID    string `json:"user_id"`
    Status    string `json:"status"`
    Amount    float64 `json:"amount"`
}

func main() {
    consumer := client.NewClient()
    
    ctx := context.Background()
    if err := consumer.Run(ctx); err != nil {
        panic(err)
    }
    defer consumer.Stop(ctx)
    
    // Subscribe to order events
    // consumer.Subscribe("orders.*", func(msg []byte) {
    //     var event OrderEvent
    //     json.Unmarshal(msg, &event)
    //     log.Printf("Processing order %s for user %s", event.OrderID, event.UserID)
    //     // Send notification logic here
    // }) // Coming soon
}
```

### 4. Dual Producer/Consumer Service

```go
// payment-service/main.go - Consumes orders, publishes payment events
package main

import (
    "context"
    "github.com/hoppermq/hopper/pkg/client"
)

func main() {
    service := client.NewClient()
    
    ctx := context.Background()
    if err := service.Run(ctx); err != nil {
        panic(err)
    }
    defer service.Stop(ctx)
    
    // service.Subscribe("orders.created", handleOrderPayment)
    // service.Subscribe("payments.retry", retryFailedPayment)
    // Both consume orders AND publish payment events
}
```

---

## 🔍 When to Choose Hopper

**🎯 Hopper excels at:**
- **Enterprise messaging with built-in transparency** → See message flow without external dashboards
- **Container/channel routing with observability** → Complex routing patterns with full visibility
- **Delivery guarantees + monitoring** → Reliable processing with real-time insights
- **Self-hosted control** → Own your message infrastructure and data
- **Microservices pub/sub** → Event-driven architectures with string-based topics

**📊 Operational advantages over alternatives:**
- **vs RabbitMQ**: Same routing power + built-in observability (no separate management UI)
- **vs Kafka**: Focused on delivery guarantees (not stream replay optimization)
- **vs SQS/Pub/Sub**: Self-hosted control + deeper operational insights
- **vs Redis**: Message durability + advanced routing (beyond simple key-value pub/sub)

---

## 🏗️ Architecture

Hopper is built around these core components:

- **Broker Core** → Message routing and delivery guarantees
- **TCP Transport** → High-performance binary protocol  
- **Event Bus** → Internal observability and monitoring hooks
- **Web UI** → Real-time dashboard for queue management
- **Client SDK** → Go client library for producers/consumers

The architecture emphasizes **transparency** - every message movement generates events that can be monitored, logged, or integrated with external observability tools.

---

## 🛣️ Roadmap

- [ ] **Complete Client SDK** → Full producer/consumer API
- [ ] **Persistent Storage** → Message durability across restarts  
- [ ] **Streamly Integration** → Native monitoring dashboard
- [ ] **Multi-language Clients** → Python, Node.js, Java SDKs
- [ ] **Performance Benchmarks** → Throughput and latency testing
- [ ] **Docker Images** → Official container distribution

---

## 🤝 Contributing

Hopper is open-source and community-driven. We welcome contributions!

- 🐛 **Bug reports** → [Open an issue](https://github.com/hoppermq/hopper/issues)
- 💡 **Feature requests** → [Start a discussion](https://github.com/hoppermq/hopper/discussions)  
- 🔧 **Pull requests** → See our [Contributing Guide](CONTRIBUTING.md)

---

## 📄 License

MIT © 2025 HopperMQ

---

> **"Message brokers shouldn't be black boxes. Hopper makes message flow transparent."**
