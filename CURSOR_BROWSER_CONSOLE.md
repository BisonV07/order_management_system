# How to Access Console in Cursor Browser

## Method 1: Right-Click Menu
1. Right-click anywhere on the page in Cursor's browser
2. Look for "Inspect" or "Inspect Element" option
3. Click it to open DevTools

## Method 2: Keyboard Shortcut
Try these keyboard shortcuts:
- **Mac**: `Cmd + Option + I` or `Cmd + Shift + I`
- **Windows/Linux**: `Ctrl + Shift + I` or `F12`

## Method 3: Menu Bar
1. Look for a menu icon (three dots or hamburger menu) in Cursor's browser
2. Navigate to: **View** → **Developer** → **Developer Tools**
3. Or: **Tools** → **Developer Tools**

## Method 4: Command Palette
1. Press `Cmd + Shift + P` (Mac) or `Ctrl + Shift + P` (Windows/Linux)
2. Type "Developer Tools" or "Console"
3. Select the option to open DevTools

## Method 5: If Console Not Available
If Cursor's browser doesn't have DevTools, you have these options:

### Option A: Use External Browser (Recommended)
1. Copy the URL from Cursor's browser
2. Paste it in Chrome/Safari
3. Use full DevTools there (F12)

### Option B: Use Browser Console via Code
Add this to your React app temporarily:

```typescript
// In your component or App.tsx
useEffect(() => {
  // Log to see if it works
  console.log('App loaded');
  
  // You can also use window.console methods
  window.console.log('Window console available');
}, []);
```

### Option C: Check Terminal Output
Some errors might appear in the terminal where you ran `npm run dev`:
- Vite compilation errors
- Network errors
- Build warnings

## What to Look For in Console

Once you have console access, check for:

1. **Errors** (red text):
   - Network errors
   - CORS errors
   - JavaScript errors
   - React errors

2. **Warnings** (yellow text):
   - Deprecation warnings
   - React warnings

3. **Network Tab**:
   - Failed requests
   - Request/response details
   - Status codes

## Quick Test

Once console is open, try this:

```javascript
// Test if console works
console.log('Console is working!');

// Test API connection
fetch('/api/v1/health')
  .then(r => r.text())
  .then(console.log)
  .catch(console.error);
```

## Alternative: Use External Browser DevTools

If Cursor's browser console is limited:

1. **Open in Chrome/Safari:**
   - Right-click the preview URL
   - Select "Open in External Browser"
   - Or manually open: `http://localhost:3000`

2. **Use Full DevTools:**
   - Press `F12` or `Cmd+Option+I` (Mac) / `Ctrl+Shift+I` (Windows)
   - Full Chrome DevTools available
   - Better debugging experience

## Pro Tip

For web development, it's often better to:
- ✅ Use Cursor for code editing
- ✅ Use external browser (Chrome/Safari) for testing
- ✅ Use browser DevTools for debugging

This gives you the best of both worlds!

