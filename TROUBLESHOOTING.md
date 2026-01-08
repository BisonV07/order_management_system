# Troubleshooting Frontend-Backend Connection

## Common Issues and Solutions

### 1. CORS Errors

**Symptom:** Browser console shows "CORS policy" errors

**Solution:** The backend now includes CORS middleware. Make sure you've restarted the backend after the latest changes.

### 2. Connection Refused

**Symptom:** "Failed to connect" or "Connection refused"

**Check:**
```bash
# Verify backend is running
curl http://localhost:8080/api/v1/health

# Check if port 8080 is in use
lsof -i :8080
```

**Solution:** Restart the backend:
```bash
cd backend
go run cmd/main.go --api --port=8080
```

### 3. 404 Not Found

**Symptom:** "404" errors when calling API

**Check:**
- Is the endpoint correct? Should be `/api/v1/products`
- Is the backend router properly set up?

**Solution:** Verify the backend router has the products endpoint registered.

### 4. 401 Unauthorized

**Symptom:** "Unauthorized" errors

**Solution:** The products endpoint should be public (no auth required). Check that the auth middleware is skipping products GET requests.

### 5. Frontend Can't Reach Backend

**Check Browser Console:**
1. Open browser DevTools (F12)
2. Go to Network tab
3. Try loading the page
4. Check what request is being made and what error you get

**Common Issues:**
- Wrong URL: Check `VITE_API_BASE_URL` in `frontend/.env`
- Backend not running: Verify backend is on port 8080
- CORS: Should be fixed with CORS middleware

### 6. Testing the Connection

**Test Backend Directly:**
```bash
# Test health endpoint
curl http://localhost:8080/api/v1/health

# Test products endpoint
curl http://localhost:8080/api/v1/products
```

**Expected Response:**
- Health: `OK`
- Products: `[]` (empty array)

### 7. Frontend Environment Variables

**Check:** `frontend/.env` should contain:
```
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

**After changing .env:** Restart the frontend dev server.

### 8. Network Tab Debugging

1. Open browser DevTools → Network tab
2. Refresh the page
3. Look for the `/products` request
4. Check:
   - Request URL (should be `http://localhost:8080/api/v1/products`)
   - Status code (should be 200)
   - Response (should be `[]`)
   - Headers (check for CORS headers)

### 9. Restart Everything

If nothing works, restart both:

**Backend:**
```bash
cd backend
# Stop with Ctrl+C, then:
go run cmd/main.go --api --port=8080
```

**Frontend:**
```bash
cd frontend
# Stop with Ctrl+C, then:
npm run dev
```

### 10. Check Backend Logs

Look at the terminal where the backend is running. You should see:
- Server starting message
- Request logs (if logging middleware is working)
- Any error messages

## Quick Diagnostic Commands

```bash
# 1. Check if backend is running
curl http://localhost:8080/api/v1/health

# 2. Test products endpoint
curl http://localhost:8080/api/v1/products

# 3. Check what's on port 8080
lsof -i :8080

# 4. Check frontend is on port 3000
lsof -i :3000
```

