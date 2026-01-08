# Browser Debugging Guide

## Quick Test in Browser Console

Open your browser console (F12) and run this to test the connection:

```javascript
// Test 1: Direct fetch through proxy
fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'test' })
})
.then(r => {
  console.log('Status:', r.status);
  console.log('Headers:', [...r.headers.entries()]);
  return r.json();
})
.then(data => {
  console.log('✅ Success:', data);
})
.catch(err => {
  console.error('❌ Error:', err);
});

// Test 2: Check API base URL
console.log('API Base URL:', import.meta.env.VITE_API_BASE_URL || '/api/v1');

// Test 3: Test axios directly
import axios from 'axios';
axios.post('/api/v1/auth/login', { username: 'test', password: 'test' })
  .then(r => console.log('✅ Axios Success:', r.data))
  .catch(e => console.error('❌ Axios Error:', e));
```

## What to Look For

### In Network Tab:
1. **Request URL**: Should be `http://localhost:3000/api/v1/auth/login`
2. **Request Method**: Should be `POST`
3. **Status Code**: Should be `200` (not 404, 500, or CORS error)
4. **Response**: Should contain `{"token":"...","user_id":1}`

### Common Issues:

#### Issue 1: 404 Not Found
**Cause**: Request not reaching backend
**Check**: 
- Is backend running? `lsof -i :8080`
- Is Vite proxy configured? Check `vite.config.ts`

#### Issue 2: CORS Error
**Cause**: Browser blocking cross-origin request
**Check**: 
- Are you using proxy? (Should use `/api/v1` not `http://localhost:8080/api/v1`)
- Check Network tab for actual request URL

#### Issue 3: Network Error / Failed to Fetch
**Cause**: Connection refused or timeout
**Check**:
- Backend running? `curl http://localhost:8080/api/v1/health`
- Frontend running? `curl http://localhost:3000`
- Firewall blocking?

#### Issue 4: 401 Unauthorized
**Cause**: Auth middleware blocking
**Check**: 
- Is `/api/v1/auth/login` in public paths?
- Check backend logs for auth errors

## Step-by-Step Debugging

1. **Open Browser DevTools** (F12)
2. **Go to Network Tab**
3. **Clear Network Log** (trash icon)
4. **Try to Login**
5. **Find `/api/v1/auth/login` request**
6. **Click on it** to see details:
   - **Headers**: Request and Response headers
   - **Preview**: Response body
   - **Response**: Raw response
   - **Timing**: Request timing

## Expected Request Details

### Request Headers:
```
POST /api/v1/auth/login HTTP/1.1
Host: localhost:3000
Content-Type: application/json
Origin: http://localhost:3000
```

### Response Headers:
```
HTTP/1.1 200 OK
Content-Type: application/json
```

### Response Body:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1
}
```

## If Still Not Working

1. **Hard Refresh**: Ctrl+Shift+R (Windows) / Cmd+Shift+R (Mac)
2. **Clear Browser Cache**: DevTools → Application → Clear Storage
3. **Try Incognito Mode**: Rules out extensions/cache
4. **Check Backend Logs**: Look for incoming requests
5. **Check Frontend Logs**: Look for proxy errors

## Share These Details

If still having issues, share:
1. Browser console errors (screenshot)
2. Network tab request details (screenshot)
3. Backend terminal output
4. Frontend terminal output

