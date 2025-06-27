#!/bin/bash

# Teste do endpoint de login
echo "Testando o endpoint de login..."
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "senha123"
  }'

echo -e "\n\n"
