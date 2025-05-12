#!/bin/bash

# Config
LOGIN_URL="http://localhost:8080/api/auth/login"
TARGET_URL="http://localhost:8080/api/services"
USERNAME="test@test.com"
PASSWORD="test"
DURATION=10
STEP=50
MAX_REQUESTS=500
TMP_RESULTS="results.txt"
TMP_DATA="graph_data.dat"

rm -f "$TMP_RESULTS" "$TMP_DATA"

# Step 1: Get JWT token
echo "🔐 Authenticating..."
TOKEN=$(curl -s -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d '{"email": "'"$USERNAME"'", "password": "'"$PASSWORD"'"}' | jq -r '.access_token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "❌ Failed to obtain token"
  exit 1
fi
echo "✅ Token: ${TOKEN:0:30}..."

# Step 2: Incremental load testing
for ((REQ=$STEP; REQ<=$MAX_REQUESTS; REQ+=STEP)); do
  echo "🚀 Running load test with $REQ requests..."
  OUTPUT=$(ddosify -t "$TARGET_URL" -h "Authorization: Bearer $TOKEN" -m GET -n $REQ -d $DURATION 2>/dev/null)
  echo "$OUTPUT" >> "$TMP_RESULTS"

  # Extract average response time in ms
  AVG_TIME=$(echo "$OUTPUT" | awk '/Avg\. Duration:/ { for(i=1;i<=NF;i++) if ($i ~ /^[0-9]+\.[0-9]+s$/) { gsub("s", "", $i); print $i; exit } }')
  if [[ -n "$AVG_TIME" ]]; then
    echo "$REQ $AVG_TIME" >> "$TMP_DATA"
    echo "$OUTPUT"
  else
    echo "⚠️ No avg time for $REQ requests. Output was:"
    echo "$OUTPUT"
  fi
done

# Step 3: Plot results using gnuplot
if command -v gnuplot &>/dev/null; then
  echo "📈 Generating graph..."
  gnuplot -persist <<-EOFMarker
    set title "Response Time vs Number of Requests"
    set xlabel "Number of Requests"
    set ylabel "Avg Response Time (ms)"
    set grid
    plot "$TMP_DATA" using 1:2 with linespoints title "Avg Time"
EOFMarker
else
  echo "⚠️ Gnuplot not found. Printing raw data:"
  cat "$TMP_DATA"
fi