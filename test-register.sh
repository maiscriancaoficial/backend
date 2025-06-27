#!/bin/bash

# Registrar usuário ADMIN
echo "Registrando usuário ADMIN..."
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Admin Teste",
    "email": "admin@exemplo.com",
    "password": "senha123",
    "role": "ADMIN"
  }'
echo -e "\n\n"

# Registrar usuário CLIENT
echo "Registrando usuário CLIENT..."
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Cliente Teste",
    "email": "cliente@exemplo.com",
    "password": "senha123",
    "role": "CLIENT"
  }'
echo -e "\n\n"

# Registrar usuário EMPLOYEE
echo "Registrando usuário EMPLOYEE..."
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Funcionario Teste",
    "email": "funcionario@exemplo.com",
    "password": "senha123",
    "role": "EMPLOYEE"
  }'
echo -e "\n\n"

# Registrar usuário AFFILIATE
echo "Registrando usuário AFFILIATE..."
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Afiliado Teste",
    "email": "afiliado@exemplo.com",
    "password": "senha123",
    "role": "AFFILIATE"
  }'
echo -e "\n\n"
