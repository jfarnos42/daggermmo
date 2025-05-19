#!/bin/bash
set -e

BASE_URL="http://localhost:8080"

echo "Testing GET /"
curl -s "$BASE_URL/"
echo -e "\n"

echo "Testing GET /httphealth"
curl -s "$BASE_URL/httphealth"
echo -e "\n"

echo "Testing GET /bdhealth"
curl -s "$BASE_URL/bdhealth"
echo -e "\n"

echo "Testing GET /listplayers (before adding player)"
curl -s "$BASE_URL/listplayers"
echo -e "\n"

echo "Testing POST /addplayer?username=testplayer&role=player"
curl -s -X POST "$BASE_URL/addplayer?username=testplayer&role=player"
echo -e "\n"

echo "Testing GET /listplayers (after adding player)"
curl -s "$BASE_URL/listplayers"
echo -e "\n"

echo "Testing POST /login?username=testplayer"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login?username=testplayer")
echo "Login response: $LOGIN_RESPONSE"

# Extraer token manualmente del JSON usando shell básico
TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*"' | sed 's/"token":"\([^"]*\)"/\1/')
ROLE=$(echo "$LOGIN_RESPONSE" | grep -o '"role":"[^"]*"' | sed 's/"role":"\([^"]*\)"/\1/')
echo "Extracted token: $TOKEN"
echo "Extracted role: $ROLE"
echo -e "\n"

if [ -n "$TOKEN" ]; then
  echo "Testing POST /logout?token=$TOKEN"
  curl -s -X POST "$BASE_URL/logout?token=$TOKEN"
  echo -e "\n"
else
  echo "No valid token to test logout"
fi

echo "Testing GET /getplayerrole?username=testplayer"
curl -s "$BASE_URL/getplayerrole?username=testplayer"
echo -e "\n"

echo "Testing GET /listplayers (after logout)"
curl -s "$BASE_URL/listplayers"
echo -e "\n"

echo "All endpoints tested."
