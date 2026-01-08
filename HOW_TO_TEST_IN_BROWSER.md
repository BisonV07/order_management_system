# How to Test in Browser Console - Step by Step Guide

## Step 1: Open Your Application

1. Make sure both services are running:
   ```bash
   # Terminal 1 - Backend
   cd backend
   go run cmd/main.go --api --port=8080
   
   # Terminal 2 - Frontend
   cd frontend
   npm run dev
   ```

2. Open your external browser (Chrome, Safari, Firefox, or Edge)

3. Go to: `http://localhost:3000`

## Step 2: Open Developer Tools (DevTools)

### Method 1: Keyboard Shortcut (Easiest)
- **Mac**: Press `Cmd + Option + I` (or `Cmd + Shift + I`)
- **Windows/Linux**: Press `F12` or `Ctrl + Shift + I`

### Method 2: Right-Click Menu
1. Right-click anywhere on the page
2. Click "Inspect" or "Inspect Element"

### Method 3: Browser Menu
- **Chrome**: Menu (⋮) → More Tools → Developer Tools
- **Safari**: Safari → Settings → Advanced → Show Develop menu, then Develop → Show Web Inspector
- **Firefox**: Menu (☰) → More Tools → Web Developer Tools

## Step 3: Access the Console Tab

Once DevTools opens, you'll see tabs at the top:
- **Elements** (or Inspector)
- **Console** ← Click this one!
- **Sources**
- **Network**
- **Application**
- etc.

Click on the **Console** tab.

## Step 4: Run Test Commands

In the Console tab, you'll see a prompt that looks like `>` or `▶`. This is where you type commands.

### Test 1: Check if Console Works
Type this and press Enter:
```javascript
console.log('Hello! Console is working!')
```

You should see: `Hello! Console is working!` printed below.

### Test 2: Test Backend Health Endpoint (Direct)
Copy and paste this, then press Enter:
```javascript
fetch('http://localhost:8080/api/v1/health')
  .then(r => r.text())
  .then(result => console.log('✅ Backend Health:', result))
  .catch(error => console.error('❌ Backend Error:', error))
```

**Expected Result:**
- ✅ Success: `✅ Backend Health: OK`
- ❌ Error: Shows error message

### Test 3: Test Through Vite Proxy
Copy and paste this, then press Enter:
```javascript
fetch('/api/v1/health')
  .then(r => r.text())
  .then(result => console.log('✅ Proxy Health:', result))
  .catch(error => console.error('❌ Proxy Error:', error))
```

**Expected Result:**
- ✅ Success: `✅ Proxy Health: OK`
- ❌ Error: Shows error message

### Test 4: Test Login Endpoint (Direct Backend)
Copy and paste this, then press Enter:
```javascript
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'test' })
})
  .then(r => r.json())
  .then(result => console.log('✅ Login Success:', result))
  .catch(error => console.error('❌ Login Error:', error))
```

**Expected Result:**
- ✅ Success: `✅ Login Success: {token: "...", user_id: 1}`
- ❌ Error: Shows error message

### Test 5: Test Login Through Proxy
Copy and paste this, then press Enter:
```javascript
fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'test' })
})
  .then(r => r.json())
  .then(result => console.log('✅ Proxy Login Success:', result))
  .catch(error => console.error('❌ Proxy Login Error:', error))
```

**Expected Result:**
- ✅ Success: `✅ Proxy Login Success: {token: "...", user_id: 1}`
- ❌ Error: Shows error message

## Step 5: Check Network Tab

1. Click on the **Network** tab in DevTools

2. You'll see a list of network requests. If the list is empty or has old requests:
   - Click the **Clear** button (🚫 icon) to clear the log

3. Now try to login from the page:
   - Enter any username and password
   - Click "Login"

4. Watch the Network tab - you should see a new request appear:
   - Look for `/api/v1/auth/login`
   - It might show as "pending" then change to a status

5. Click on the `/api/v1/auth/login` request to see details:

   **Headers Tab:**
   - **Request URL**: Should be `http://localhost:3000/api/v1/auth/login`
   - **Request Method**: Should be `POST`
   - **Status Code**: Should be `200` (green) or an error code (red)

   **Preview Tab:**
   - Should show: `{token: "...", user_id: 1}`

   **Response Tab:**
   - Shows the raw response JSON

## Step 6: Check for Errors

### In Console Tab:
Look for messages in red - these are errors. Common errors:

- **ERR_CONNECTION_REFUSED**: Backend not running
- **ERR_FAILED**: Network error
- **CORS error**: Cross-origin issue
- **404 Not Found**: Route not found
- **500 Internal Server Error**: Backend error

### In Network Tab:
- **Red status codes** (4xx, 5xx): Request failed
- **Gray/Canceled**: Request was canceled
- **Pending forever**: Request stuck

## Quick Copy-Paste Test Suite

Copy all of this into the console at once (it will run all tests):

```javascript
// Complete Test Suite
console.log('🧪 Starting tests...\n');

// Test 1: Backend Health (Direct)
fetch('http://localhost:8080/api/v1/health')
  .then(r => r.text())
  .then(result => console.log('✅ Test 1 - Backend Health:', result))
  .catch(error => console.error('❌ Test 1 - Backend Error:', error));

// Test 2: Backend Health (Proxy)
fetch('/api/v1/health')
  .then(r => r.text())
  .then(result => console.log('✅ Test 2 - Proxy Health:', result))
  .catch(error => console.error('❌ Test 2 - Proxy Error:', error));

// Test 3: Login (Direct)
fetch('http://localhost:8080/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'test' })
})
  .then(r => r.json())
  .then(result => console.log('✅ Test 3 - Direct Login:', result))
  .catch(error => console.error('❌ Test 3 - Direct Login Error:', error));

// Test 4: Login (Proxy)
fetch('/api/v1/auth/login', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({ username: 'test', password: 'test' })
})
  .then(r => r.json())
  .then(result => console.log('✅ Test 4 - Proxy Login:', result))
  .catch(error => console.error('❌ Test 4 - Proxy Login Error:', error));
```

## What to Look For

### ✅ Everything Working:
- All tests show ✅ (green checkmarks)
- Network tab shows status 200
- Response contains token

### ❌ Issues Found:
- **Test 1 fails, Test 2 works**: Backend not accessible directly (CORS issue)
- **Test 1 works, Test 2 fails**: Vite proxy not working
- **Both fail**: Backend not running or wrong port
- **CORS errors**: Backend CORS middleware issue
- **404 errors**: Route not registered

## Screenshots to Share

If you need help, take screenshots of:
1. **Console tab** - showing the test results
2. **Network tab** - showing the `/api/v1/auth/login` request details
3. **Any error messages** - red text in console

## Tips

- **Clear console**: Click the 🚫 icon or type `clear()` and press Enter
- **Copy results**: Right-click on any console message → Copy
- **Filter console**: Type in the filter box to search for specific messages
- **Preserve log**: Check "Preserve log" checkbox to keep messages after page reload

## Next Steps

After running tests:
1. Share the results (what you see in console)
2. Share any error messages
3. Share Network tab details if requests are failing

This will help identify exactly what's wrong!

