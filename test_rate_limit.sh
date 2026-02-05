#!/bin/bash

URL="http://localhost:8080/api/resource"

echo "🚦 Testing HoldUp Rate Limiter"
echo "--------------------------------"

for i in {1..15}
do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL")
  echo "Request $i → Status: $STATUS"
  sleep 0.2
done

echo ""
echo "⏳ Waiting for token refill..."
sleep 3

echo ""
echo "🔄 Retesting after refill"
for i in {1..5}
do
  STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$URL")
  echo "Retry $i → Status: $STATUS"
done
