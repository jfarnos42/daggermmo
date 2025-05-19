#!/bin/bash

BASE_URL="http://localhost:8080"

# Primero chequea la salud de la base de datos
DB_STATUS=$(curl -s "${BASE_URL}/bdhealth")
if [[ "$DB_STATUS" != "bd: ok" ]]; then
  echo "Database is not available. Aborting."
  exit 1
fi

# Asegura que el usuario exista (crea el jugador si es necesario)
echo "Adding player 'testuser' (if not exists)..."
curl -s "${BASE_URL}/addplayer?username=testuser" || true

# Login y obtener token
echo -e "\nLogging in user 'testuser'..."
RAW_RESPONSE=$(curl -s "${BASE_URL}/login?username=testuser")
echo "Raw login response: $RAW_RESPONSE"
TOKEN=$(echo "$RAW_RESPONSE" | jq -r '.token')

if [[ -z "$TOKEN" ]]; then
  echo "Failed to get login token."
  exit 1
fi

echo "Extracted token: $TOKEN"

# Logout usando el token
echo -e "\nLogging out with token..."
LOGOUT_RESULT=$(curl -s "${BASE_URL}/logout?token=$TOKEN")
echo "Logout response: $LOGOUT_RESULT"

# Intentar logout otra vez con el mismo token (debería fallar)
echo -e "\nTrying to logout again with same token (should fail)..."
LOGOUT_RESULT2=$(curl -s "${BASE_URL}/logout?token=$TOKEN")
echo "Logout response: $LOGOUT_RESULT2"
