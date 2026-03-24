#!/bin/bash
# Test script for State Sync API endpoints
# This script tests the cloud synchronization features

BASE_URL="http://localhost:8080"
API_KEY="test-agent-key-12345"  # Replace with actual test agent key

echo "🧪 Testing State Sync API Endpoints..."
echo ""

# Test 1: Create agent state snapshot
echo "1️⃣ Creating agent state snapshot..."
RESPONSE=$(curl -s -X POST "${BASE_URL}/state/sync/agents/test-agent-001/snapshots" \
  -H "X-API-KEY: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "snapshot_type": "full",
    "metadata": {
      "test": true,
      "timestamp": "'$(date +%s)'"
    }
  }')

echo "$RESPONSE" | python3 -m json.tool
echo ""

# Test 2: Get current snapshot
echo "2️⃣ Getting current agent snapshot..."
SNAPSHOT_ID=$(curl -s -X GET "${BASE_URL}/state/sync/agents/test-agent-001/snapshot/current" \
  -H "X-API-KEY: ${API_KEY}" | python3 -c "import sys, json; print(json.load(sys.stdin)['id'])")

echo "$SNAPSHOT_ID"
echo ""

# Test 3: Get snapshot history
echo "3️⃣ Getting snapshot history..."
curl -s -X GET "${BASE_URL}/state/sync/agents/test-agent-001/snapshots?limit=5" \
  -H "X-API-KEY: ${API_KEY}" | python3 -m json.tool
echo ""

# Test 4: Store task execution
echo "4️⃣ Storing task execution..."
TASK_ID=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
curl -s -X POST "${BASE_URL}/state/sync/agents/test-agent-001/executions" \
  -H "X-API-KEY: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d "{
    \"task_id\": \"$TASK_ID\",
    \"title\": \"Test Task\",
    \"description\": \"Testing state sync functionality\",
    \"response_text\": \"Task completed successfully!\",
    \"status\": \"completed\",
    \"execution_time_ms\": 1250,
    \"metadata\": {
      \"test_data\": true
    }
  }" | python3 -m json.tool
echo ""

# Test 5: Get task completion metrics
echo "5️⃣ Getting task completion metrics..."
curl -s -X GET "${BASE_URL}/state/sync/agents/test-agent-001/metrics/task-completion?days=7" \
  -H "X-API-KEY: ${API_KEY}" | python3 -m json.tool
echo ""

# Test 6: Export task history (JSON format)
echo "6️⃣ Exporting task history (JSON)..."
curl -s -X GET "${BASE_URL}/state/sync/agents/test-agent-001/export?format=json&limit=10" \
  -H "X-API-KEY: ${API_KEY}" | python3 -m json.tool
echo ""

# Test 7: Export task history (CSV format)
echo "7️⃣ Exporting task history (CSV)..."
curl -s -X GET "${BASE_URL}/state/sync/agents/test-agent-001/export?format=csv&limit=10" \
  -H "X-API-KEY: ${API_KEY}" | head -5
echo ""

# Test 8: Add team member
echo "8️⃣ Adding team member..."
USER_ID=$(python3 -c "import uuid; print(str(uuid.uuid4()))")
curl -s -X POST "${BASE_URL}/state/sync/organisations/test-agent-001/members" \
  -H "X-API-KEY: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$USER_ID\",
    \"role\": \"member\",
    \"status\": \"invited\"
  }" | python3 -m json.tool
echo ""

# Test 9: Get team members
echo "9️⃣ Getting team members..."
curl -s -X GET "${BASE_URL}/state/sync/organisations/test-agent-001/members" \
  -H "X-API-KEY: ${API_KEY}" | python3 -m json.tool
echo ""

# Test 10: Update member status
echo "🔟 Updating team member status..."
curl -s -X PATCH "${BASE_URL}/state/sync/organisations/test-agent-001/members/${USER_ID}/status" \
  -H "X-API-KEY: ${API_KEY}" \
  -H "Content-Type: application/json" \
  -d '{
    "status": "active"
  }' | python3 -m json.tool
echo ""

echo "✅ All tests completed!"
