# Backend Connection Troubleshooting

## Network Error on Login

If you're seeing "Network Error" on the login page, the frontend cannot connect to the backend.

## Quick Checks

### 1. Is Backend Running?

Check if the backend is running:
```bash
lsof -i :8080
```

If nothing shows up, start the backend:
```bash
cd backend
go run cmd/main.go --api --port=8080
```

### 2. Test Backend Directly

Test if the backend is responding:
```bash
# Test health endpoint
curl http://localhost:8080/api/v1/health

# Should return: OK

# Test login endpoint
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test"}'

# Should return: {"token":"...","user_id":1}
```

### 3. Check Browser Console

1. Open browser DevTools (F12)
2. Go to **Console** tab - look for errors
3. Go to **Network** tab - refresh and try login
4. Look for the `/auth/login` request:
   - Status code (should be 200)
   - Response (should show token)
   - Any CORS errors

### 4. Common Issues

**Backend not started:**
- Solution: Start backend with `go run cmd/main.go --api --port=8080`

**Port conflict:**
- Another service is using port 8080
- Solution: Change port or stop conflicting service
```bash
# Find what's using port 8080
lsof -i :8080

# Or change backend port
go run cmd/main.go --api --port=8081
# Then update frontend .env: VITE_API_BASE_URL=http://localhost:8081/api/v1
```

**CORS errors:**
- Check browser console for CORS messages
- Solution: Ensure CORS middleware is first in router

**Firewall/Network:**
- Check if localhost:8080 is accessible
- Try: `curl http://localhost:8080/api/v1/health`

## Expected Backend Output

When backend starts successfully, you should see:
```
Starting API server on port 8080...
Server listening on http://localhost:8080
Health check: http://localhost:8080/api/v1/health
Products: http://localhost:8080/api/v1/products
```

## Verify Everything is Working

1. **Backend health check:**
   ```bash
   curl http://localhost:8080/api/v1/health
   # Expected: OK
   ```

2. **Backend login test:**
   ```bash
   curl -X POST http://localhost:8080/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"username":"test","password":"test"}'
   # Expected: {"token":"...","user_id":1}
   ```

3. **Frontend:**
   - Open http://localhost:3000
   - Should see login page
   - Enter any username/password
   - Should redirect to dashboard on success

## Still Having Issues?

1. Check backend terminal for error messages
2. Check browser console (F12) for detailed errors
3. Verify both services are running:
   - Backend: `lsof -i :8080`
   - Frontend: `lsof -i :3000`

