# Troubleshooting: No Console Logs Appearing

## Quick Checks

### 1. Is Frontend Actually Running?

Check your terminal where you ran `npm run dev`. You should see:
```
  VITE v5.x.x  ready in xxx ms

  ➜  Local:   http://localhost:3000/
  ➜  Network: use --host to expose
```

If you don't see this, the frontend isn't running!

**Fix:**
```bash
cd frontend
npm run dev
```

### 2. Is the Page Actually Loading?

When you go to `http://localhost:3000`, do you see:
- ✅ The login page? → Frontend is loading
- ❌ Blank page? → Frontend not loading
- ❌ Error page? → Check terminal for errors

### 3. Check Browser Console Settings

**Make sure console is not filtered:**

1. Open DevTools (F12)
2. Click Console tab
3. Look for filter icons/buttons:
   - **All levels** should be checked (not just Errors)
   - **Hide network messages** should be OFF
   - **Preserve log** checkbox - try checking this

### 4. Check if Console is Actually Open

Sometimes console is hidden or minimized:
- Look for a small arrow/chevron to expand it
- Try dragging the console panel to make it bigger
- Check if console is docked to bottom/side (try different positions)

### 5. Check for JavaScript Errors

If there are JavaScript errors, the app might not load at all:

1. Open Console tab
2. Look for RED error messages
3. These will tell you what's wrong

Common errors:
- `Failed to load resource` → File not found
- `Uncaught SyntaxError` → JavaScript syntax error
- `Cannot find module` → Import error

## Step-by-Step Debugging

### Step 1: Verify Frontend is Running

**In Terminal:**
```bash
# Check if port 3000 is in use
lsof -i :3000

# Or check if you can access it
curl http://localhost:3000
```

**What to look for:**
- Terminal shows Vite dev server running
- Port 3000 is listening
- curl returns HTML (not connection refused)

### Step 2: Check Browser Console Settings

1. Open `http://localhost:3000` in browser
2. Press F12 to open DevTools
3. Click **Console** tab
4. Look at the top of Console panel:
   - There should be filter buttons (All, Errors, Warnings, Info)
   - Make sure **All** is selected (not just Errors)
   - Uncheck "Hide network" if checked

### Step 3: Try Simple Console Test

In Console, type this and press Enter:
```javascript
console.log('TEST - Can you see this?')
```

**Expected:** You should see `TEST - Can you see this?` appear

**If you don't see it:**
- Console might be filtered
- Try clicking "Clear" button and try again
- Check if "Preserve log" is checked

### Step 4: Check Network Tab

1. Click **Network** tab in DevTools
2. Refresh the page (F5)
3. Look for:
   - **index.html** - Should be status 200
   - **main.tsx** or **main.js** - Should be status 200
   - Any **red** entries = failed to load

### Step 5: Check Terminal for Errors

Look at the terminal where you ran `npm run dev`:

**Good signs:**
```
✓ built in xxxms
```

**Bad signs:**
```
✗ Failed to compile
Error: ...
```

If you see errors, fix them first!

## Common Issues

### Issue 1: Console Filtered

**Symptom:** Console is empty even though app is running

**Fix:**
1. In Console tab, click filter buttons
2. Make sure "All levels" is selected
3. Uncheck "Hide network messages"
4. Click "Clear" and refresh page

### Issue 2: Frontend Not Compiling

**Symptom:** Blank page, errors in terminal

**Fix:**
1. Check terminal for compilation errors
2. Fix any TypeScript/JavaScript errors
3. Restart dev server: `npm run dev`

### Issue 3: Browser Cache

**Symptom:** Old version of app, console shows old errors

**Fix:**
1. Hard refresh: `Ctrl+Shift+R` (Windows) / `Cmd+Shift+R` (Mac)
2. Or clear browser cache
3. Try incognito/private mode

### Issue 4: Console Panel Hidden

**Symptom:** Can't see console at all

**Fix:**
1. Press F12 again to toggle DevTools
2. Try different dock positions (bottom, side, separate window)
3. Look for Console tab - might be hidden behind other tabs

## Quick Diagnostic Commands

Run these in browser console (if you can access it):

```javascript
// Test 1: Is console working?
console.log('✅ Console is working!');

// Test 2: Is page loaded?
console.log('Page URL:', window.location.href);

// Test 3: Check for errors
window.onerror = function(msg, url, line) {
  console.error('JavaScript Error:', msg, 'at', url, ':', line);
};

// Test 4: Check if React is loaded
console.log('React version:', window.React?.version || 'Not found');
```

## What to Share

If still not working, share:

1. **Terminal output** - What does `npm run dev` show?
2. **Browser screenshot** - What do you see at `http://localhost:3000`?
3. **Console screenshot** - What does Console tab show?
4. **Network tab** - Any failed requests?

## Alternative: Check Terminal Output

If browser console isn't working, check the terminal where frontend is running:

- **Vite errors** will show there
- **Compilation errors** will show there
- **Network requests** might be logged there

Look for any RED error messages in the terminal!

