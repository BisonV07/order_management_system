# External Browser Debugging Guide

## Current Status

✅ **Backend**: Running on port 8080
✅ **Frontend**: Running on port 3000  
✅ **Curl Tests**: Both work perfectly
❌ **Browser**: Not connecting

## Not a Docker Issue

Docker is only used for PostgreSQL (port 5432). Backend and frontend run directly:
- Backend: `go run cmd/main.go --api --port=8080`
- Frontend: `npm run dev` (Vite on port 3000)

## Possible Issues

### Issue 1: IPv6 vs IPv4 Binding

Backend might be listening on IPv6 (`::1`) but browser tries IPv4 (`127.0.0.1`).

**Check:**
```bash
# Test IPv4 specifically
curl http://127.0.0.1:8080/api/v1/health

# Test IPv6 specifically  
curl http://[::1]:8080/api/v1/health
```

**Fix:** Ensure backend listens on both IPv4 and IPv6 (or just IPv4).

### Issue 2: Browser CORS Preflight

Browser sends OPTIONS request first. Check if it's handled.

**Test in Browser Console:**
```javascript
// Test direct fetch
fetch('http://localhost:8080/api/v1/health')
  .then(r => r.text())
  .then(console.log)
  .catch(console.error);

// Test through proxy
fetch('/api/v1/health')
  .then(r => r.text())
  .then(console.log)
  .catch(console.error);
```

### Issue 3: Proxy Not Working in Browser

Vite proxy might not be forwarding correctly.

**Check:**
1. Open browser DevTools (F12)
2. Go to Network tab
3. Try to login
4. Look for `/api/v1/auth/login` request:
   - **URL**: Should be `http://localhost:3000/api/v1/auth/login`
   - **Status**: Should be 200 (not 404 or ERR_FAILED)
   - **Response**: Should show token JSON

### Issue 4: Browser Cache

Old failed requests might be cached.

**Fix:**
1. Hard refresh: `Ctrl+Shift+R` (Windows) / `Cmd+Shift+R` (Mac)
2. Clear browser cache
3. Try incognito/private mode

## Step-by-Step Debugging

### Step 1: Check Browser Console
1. Open external browser (Chrome/Safari)
2. Go to `http://localhost:3000`
3. Open DevTools (F12)
4. Go to **Console** tab
5. Look for errors (red text)

### Step 2: Check Network Tab
1. In DevTools, go to **Network** tab
2. Clear network log
3. Try to login
4. Find `/api/v1/auth/login` request
5. Click on it and check:
   - **Status Code**: Should be 200
   - **Request URL**: Should be `http://localhost:3000/api/v1/auth/login`
   - **Response**: Should contain `{"token":"...","user_id":1}`

### Step 3: Test Direct Backend Access
In browser console:
```javascript
fetch('http://localhost:8080/api/v1/health')
  .then(r => r.text())
  .then(console.log)
  .catch(e => console.error('Error:', e));
```

If this fails, backend might not be accessible from browser.

### Step 4: Test Proxy
In browser console:
```javascript
fetch('/api/v1/health')
  .then(r => r.text())
  .then(console.log)
  .catch(e => console.error('Proxy Error:', e));
```

If this fails, Vite proxy isn't working.

## Common Error Messages

### "ERR_CONNECTION_REFUSED"
- Backend not running or wrong port
- Check: `lsof -i :8080`

### "ERR_FAILED" or "Network Error"
- Proxy not forwarding
- CORS issue
- Check Network tab for details

### "404 Not Found"
- Route not registered
- Wrong URL path
- Check backend routes

### CORS Error
- Backend CORS headers missing
- Check CORS middleware is first

## Quick Fixes

### Fix 1: Restart Both Services
```bash
# Terminal 1 - Backend
cd backend
go run cmd/main.go --api --port=8080

# Terminal 2 - Frontend  
cd frontend
npm run dev
```

### Fix 2: Check Backend Binding
Backend should listen on `0.0.0.0` or `127.0.0.1`, not just IPv6.

### Fix 3: Verify URLs
- Frontend: `http://localhost:3000`
- Backend: `http://localhost:8080`
- API calls should go through proxy: `/api/v1/*`

## Share These Details

If still not working, share:
1. Browser console errors (screenshot)
2. Network tab request details (screenshot)
3. Backend terminal output
4. Frontend terminal output
5. Result of browser console test commands above

