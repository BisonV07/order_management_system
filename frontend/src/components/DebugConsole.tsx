import { useEffect, useState } from 'react'

/**
 * DebugConsole - Displays console output on the page
 * Useful when Cursor's browser doesn't have DevTools
 * Remove this component in production
 */
const DebugConsole = () => {
  const [logs, setLogs] = useState<Array<{ type: string; message: string; timestamp: Date }>>([])
  const [isCollapsed, setIsCollapsed] = useState(false)

  useEffect(() => {
    // Override console methods to capture logs
    const originalLog = console.log
    const originalError = console.error
    const originalWarn = console.warn

    console.log = (...args: any[]) => {
      originalLog.apply(console, args)
      setLogs(prev => [...prev, { type: 'log', message: args.map(String).join(' '), timestamp: new Date() }])
    }

    console.error = (...args: any[]) => {
      originalError.apply(console, args)
      setLogs(prev => [...prev, { type: 'error', message: args.map(String).join(' '), timestamp: new Date() }])
    }

    console.warn = (...args: any[]) => {
      originalWarn.apply(console, args)
      setLogs(prev => [...prev, { type: 'warn', message: args.map(String).join(' '), timestamp: new Date() }])
    }

    // Cleanup
    return () => {
      console.log = originalLog
      console.error = originalError
      console.warn = originalWarn
    }
  }, [])

  if (logs.length === 0) return null

  return (
    <div style={{
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      maxHeight: isCollapsed ? '40px' : '150px',
      overflowY: 'auto',
      backgroundColor: '#1e1e1e',
      color: '#fff',
      padding: '8px',
      fontSize: '11px',
      fontFamily: 'monospace',
      zIndex: 9999,
      borderTop: '2px solid #444',
      transition: 'max-height 0.3s ease'
    }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: isCollapsed ? '0' : '5px' }}>
        <strong style={{ fontSize: '12px' }}>Debug Console ({logs.length} logs)</strong>
        <div style={{ display: 'flex', gap: '5px' }}>
          <button
            onClick={() => setIsCollapsed(!isCollapsed)}
            style={{
              backgroundColor: '#555',
              color: '#fff',
              border: 'none',
              padding: '4px 8px',
              borderRadius: '3px',
              cursor: 'pointer',
              fontSize: '10px'
            }}
          >
            {isCollapsed ? '▼' : '▲'}
          </button>
          <button
            onClick={() => setLogs([])}
            style={{
              backgroundColor: '#666',
              color: '#fff',
              border: 'none',
              padding: '4px 8px',
              borderRadius: '3px',
              cursor: 'pointer',
              fontSize: '10px'
            }}
          >
            Clear
          </button>
        </div>
      </div>
      {!isCollapsed && (
        <div style={{ maxHeight: '120px', overflowY: 'auto' }}>
          {logs.slice(-15).map((log, idx) => (
            <div
              key={idx}
              style={{
                color: log.type === 'error' ? '#ff6b6b' : log.type === 'warn' ? '#ffd93d' : '#51cf66',
                marginBottom: '1px',
                padding: '1px 0',
                fontSize: '10px',
                lineHeight: '1.3'
              }}
            >
              [{log.timestamp.toLocaleTimeString()}] {log.type.toUpperCase()}: {log.message.substring(0, 100)}
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default DebugConsole

