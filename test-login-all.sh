#!/bin/bash

# Teste do login de ADMIN
echo "Testando login de ADMIN..."
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@exemplo.com",
    "password": "senha123"
  }'
echo -e "\n\n"

# Teste do login de CLIENT
echo "Testando login de CLIENT..."
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "cliente@exemplo.com",
    "password": "senha123"
  }'
echo -e "\n\n"

# Teste do login de EMPLOYEE
echo "Testando login de EMPLOYEE..."
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "funcionario@exemplo.com",
    "password": "senha123"
  }'
echo -e "\n\n"

# Teste do login de AFFILIATE
echo "Testando login de AFFILIATE..."
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "afiliado@exemplo.com",
    "password": "senha123"
  }'
echo -e "\n\n"
