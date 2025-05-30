#!/bin/bash
# port-forward-services.sh

echo "Setting up port forwarding for all services..."

kubectl port-forward -n ecommerce svc/postgres 5434:5432 &
kubectl port-forward -n ecommerce svc/user 8001:8080 &
kubectl port-forward -n ecommerce svc/product 8002:8080 &
kubectl port-forward -n ecommerce svc/shop 8003:8080 &
kubectl port-forward -n ecommerce svc/order 8004:8080 &
kubectl port-forward -n ecommerce svc/warehouse 8005:8080 &

echo "Port forwarding setup complete!"
echo "Services available at:"
echo "  User:      http://localhost:8001/api/users"
echo "  Product:   http://localhost:8002/api/products"
echo "  Shop:      http://localhost:8003/api/shops"
echo "  Order:     http://localhost:8004/api/orders"
echo "  Warehouse: http://localhost:8005/api/warehouses"

wait