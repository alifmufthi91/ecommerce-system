# E-commerce System

A simple e-commerce system built as a microservices monorepo, featuring multiple services including user management, product catalog, shop management, order processing, and warehouse operations.

## 🏗️ Architecture

This system consists of the following microservices:
- **User Service** - User authentication
- **Product Service** - Product catalog and inventory
- **Shop Service** - Shop management
- **Order Service** - Order processing
- **Warehouse Service** - Inventory and stock management

## 🚀 Getting Started

### Option 1: Kubernetes (Recommended)

#### Prerequisites
- [Minikube](https://minikube.sigs.k8s.io/docs/start)
- [Skaffold](https://skaffold.dev/docs/install/)
- [Docker](https://docs.docker.com/get-docker/)

#### Steps
1. Start Minikube:
   ```bash
   minikube start
   ```

2. Deploy with Skaffold:
   ```bash
   skaffold dev
   ```

### Option 2: Docker Compose

#### Prerequisites
- [Docker](https://docs.docker.com/get-docker/)

#### Steps
1. Navigate to the build directory:
   ```bash
   cd build
   ```

2. Start all services:
   ```bash
   docker-compose up
   ```

### Option 3: Local Development

#### Prerequisites
- [Air](https://github.com/air-verse/air) - Live reload for Go applications
- [Swaggo](https://github.com/swaggo/swag) - Swagger documentation generator

#### Steps
1. Navigate to a specific service:
   ```bash
   cd services/{service_name}
   ```

2. Start the service:
   ```bash
   make start
   ```

## 🔌 Port Forwarding

For Kubernetes deployments, use the provided script to forward ports:
```bash
./port-forwarding.sh
```

## 📚 API Documentation

Swagger documentation is available for each service:

| Service | Swagger URL |
|---------|-------------|
| User Service | http://localhost:8001/api/docs/swagger/index.html |
| Product Service | http://localhost:8002/api/docs/swagger/index.html |
| Shop Service | http://localhost:8003/api/docs/swagger/index.html |
| Order Service | http://localhost:8004/api/docs/swagger/index.html |
| Warehouse Service | http://localhost:8005/api/docs/swagger/index.html |

## 📮 Postman Collection

Import the Postman collection to test all API endpoints:

1. Download the collection: [`ecommerce.postman_collection.json`](./ecommerce.postman_collection.json)
2. Open Postman
3. Click "Import" → "Upload Files"
4. Select the downloaded JSON file

### Environment Variables
Set these variables in Postman:
- `base_url`: `http://localhost` (for local development)
- `user_port`: `8001`
- `product_port`: `8002`
- `shop_port`: `8003`
- `order_port`: `8004`
- `warehouse_port`: `8005`

## 🧪 Testing

### Unit Tests
Run unit tests for any service:
```bash
cd services/{service_name}
make test
```

### Integration Tests
The warehouse service includes integration tests for concurrent stock reservation handling:
```bash
cd services/warehouse/tests
# Run integration tests here
```

## 📊 Monitoring

### Health Checks
Monitor service health at:
```
http://localhost:{service_port}/health
```

### Metrics
Prometheus metrics are available at:
```
http://localhost:{service_port}/metrics
```

## 🛠️ Development

### Service Ports
- User Service: `8001`
- Product Service: `8002`
- Shop Service: `8003`
- Order Service: `8004`
- Warehouse Service: `8005`

---

## 📝 Notes

This system demonstrates microservices architecture with proper service separation, health monitoring, and concurrent handling capabilities.

*README enhanced by GitHub Copilot*