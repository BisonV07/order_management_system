import { useState, useEffect } from 'react'
import { useAuth } from '../context/AuthContext'
import { adminService } from '../services/api'
import type { SystemMetrics, DockerMetrics, PostgreSQLMetrics } from '../types'
import '../App.css'

const Metrics = () => {
  const { role } = useAuth()
  const isAdmin = role === 'ADMIN' || role === 'admin'

  const [metrics, setMetrics] = useState<SystemMetrics | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [autoRefresh, setAutoRefresh] = useState(false)
  const [refreshInterval, setRefreshInterval] = useState<NodeJS.Timeout | null>(null)

  const loadMetrics = async () => {
    try {
      setLoading(true)
      const data = await adminService.getMetrics()
      setMetrics(data)
      setError(null)
    } catch (err: any) {
      setError(`Failed to load metrics: ${err?.response?.data?.message || err?.message}`)
      console.error(err)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (isAdmin) {
      loadMetrics()
    }
    return () => {
      if (refreshInterval) {
        clearInterval(refreshInterval)
      }
    }
  }, [isAdmin])

  useEffect(() => {
    if (autoRefresh) {
      const interval = setInterval(() => {
        loadMetrics()
      }, 5000) // Refresh every 5 seconds
      setRefreshInterval(interval)
      return () => clearInterval(interval)
    } else {
      if (refreshInterval) {
        clearInterval(refreshInterval)
        setRefreshInterval(null)
      }
    }
  }, [autoRefresh])

  const formatBytes = (bytes: number | undefined): string => {
    if (!bytes) return '0 B'
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    if (bytes === 0) return '0 B'
    const i = Math.floor(Math.log(bytes) / Math.log(1024))
    return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
  }

  const parsePercent = (value: string | undefined): number => {
    if (!value) return 0
    const cleaned = value.replace('%', '').trim()
    const parsed = parseFloat(cleaned)
    return isNaN(parsed) ? 0 : parsed
  }

  if (!isAdmin) {
    return (
      <div className="app">
        <div className="container">
          <div className="card" style={{ textAlign: 'center', padding: '60px 20px' }}>
            <div style={{ fontSize: '64px', marginBottom: '20px' }}>🚫</div>
            <h1 style={{ color: 'var(--danger)', marginBottom: '16px' }}>Access Denied</h1>
            <p style={{ color: 'var(--gray)', marginBottom: '8px' }}>You must be an admin to access this page.</p>
            <p className="badge badge-info" style={{ fontSize: '12px', marginTop: '16px', display: 'inline-block' }}>
              Your current role: {role || 'not set'}
            </p>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="app">
      <div className="container" style={{ paddingBottom: '200px' }}>
        <div className="card mb-3" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
          <div className="flex-between">
            <div>
              <h1 style={{ 
                margin: 0,
                background: 'linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%)',
                WebkitBackgroundClip: 'text',
                WebkitTextFillColor: 'transparent'
              }}>
                📊 System Metrics
              </h1>
              <p style={{ margin: '8px 0 0', color: 'var(--gray)', fontSize: '14px' }}>
                Monitor Docker and PostgreSQL performance metrics
              </p>
            </div>
            <div style={{ display: 'flex', gap: '12px', alignItems: 'center' }}>
              <label style={{ display: 'flex', alignItems: 'center', gap: '8px', cursor: 'pointer' }}>
                <input
                  type="checkbox"
                  checked={autoRefresh}
                  onChange={(e) => setAutoRefresh(e.target.checked)}
                  style={{ cursor: 'pointer' }}
                />
                <span style={{ fontSize: '14px', color: 'var(--gray)' }}>Auto-refresh (5s)</span>
              </label>
              <button
                onClick={loadMetrics}
                className="btn btn-primary"
                style={{ fontSize: '14px', padding: '10px 20px' }}
                disabled={loading}
              >
                {loading ? '⏳ Loading...' : '🔄 Refresh'}
              </button>
            </div>
          </div>
        </div>

        {error && (
          <div className="alert alert-error" style={{ marginBottom: '20px' }}>
            <span>⚠️</span>
            <span>{error}</span>
          </div>
        )}

        {loading && !metrics ? (
          <div className="card" style={{ textAlign: 'center', padding: '60px 20px' }}>
            <div className="loading" style={{ margin: '0 auto 20px' }}></div>
            <p style={{ color: 'var(--gray)' }}>Loading metrics...</p>
          </div>
        ) : metrics ? (
          <div className="grid" style={{ gridTemplateColumns: 'repeat(auto-fit, minmax(400px, 1fr))', gap: '24px' }}>
            {/* Docker Metrics */}
            <div className="card" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
              <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                🐳 Docker Metrics
              </h2>
              
              {metrics.docker.error ? (
                <div className="alert alert-error">
                  <span>⚠️</span>
                  <span>{metrics.docker.error}</span>
                </div>
              ) : (
                <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                  <div>
                    <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Container Status</div>
                    <div className="badge" style={{
                      background: metrics.docker.status === 'running' ? 'rgba(16, 185, 129, 0.1)' : 'rgba(239, 68, 68, 0.1)',
                      color: metrics.docker.status === 'running' ? 'var(--success)' : 'var(--danger)',
                      fontSize: '14px',
                      padding: '6px 12px',
                      display: 'inline-block'
                    }}>
                      {metrics.docker.status || 'Unknown'}
                    </div>
                  </div>

                  {metrics.docker.cpu_percent && (
                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>CPU Usage</div>
                      <div style={{ fontSize: '24px', fontWeight: '700', color: 'var(--primary)' }}>
                        {metrics.docker.cpu_percent}
                      </div>
                      <div style={{ width: '100%', height: '8px', background: 'var(--gray-lighter)', borderRadius: '4px', marginTop: '8px', overflow: 'hidden' }}>
                        <div style={{
                          width: `${Math.min(parsePercent(metrics.docker.cpu_percent), 100)}%`,
                          height: '100%',
                          background: 'linear-gradient(90deg, var(--primary) 0%, var(--primary-dark) 100%)',
                          transition: 'width 0.3s ease'
                        }}></div>
                      </div>
                    </div>
                  )}

                  {metrics.docker.memory_percent && (
                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Memory Usage</div>
                      <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--warning)' }}>
                        {metrics.docker.memory_percent}
                      </div>
                      <div style={{ fontSize: '14px', color: 'var(--gray)', marginTop: '4px' }}>
                        {metrics.docker.memory_usage}
                      </div>
                      <div style={{ width: '100%', height: '8px', background: 'var(--gray-lighter)', borderRadius: '4px', marginTop: '8px', overflow: 'hidden' }}>
                        <div style={{
                          width: `${Math.min(parsePercent(metrics.docker.memory_percent), 100)}%`,
                          height: '100%',
                          background: 'linear-gradient(90deg, var(--warning) 0%, #f59e0b 100%)',
                          transition: 'width 0.3s ease'
                        }}></div>
                      </div>
                    </div>
                  )}

                  {metrics.docker.network_io && (
                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Network I/O</div>
                      <div style={{ fontSize: '16px', fontWeight: '600', color: 'var(--dark)' }}>
                        {metrics.docker.network_io}
                      </div>
                    </div>
                  )}

                  {metrics.docker.block_io && (
                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Block I/O</div>
                      <div style={{ fontSize: '16px', fontWeight: '600', color: 'var(--dark)' }}>
                        {metrics.docker.block_io}
                      </div>
                    </div>
                  )}

                  {metrics.docker.started_at && (
                    <div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Started At</div>
                      <div style={{ fontSize: '14px', color: 'var(--dark)' }}>
                        {new Date(metrics.docker.started_at).toLocaleString()}
                      </div>
                    </div>
                  )}
                </div>
              )}
            </div>

            {/* PostgreSQL Metrics */}
            <div className="card" style={{ background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
              <h2 style={{ marginTop: 0, fontSize: '1.5rem', fontWeight: '700', color: 'var(--dark)', marginBottom: '20px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                🐘 PostgreSQL Metrics
              </h2>
              
              <div style={{ display: 'flex', flexDirection: 'column', gap: '16px' }}>
                <div>
                  <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Database Size</div>
                  <div style={{ fontSize: '24px', fontWeight: '700', color: 'var(--info)' }}>
                    {metrics.postgresql.database_size || 'N/A'}
                  </div>
                  {metrics.postgresql.database_size_bytes && (
                    <div style={{ fontSize: '12px', color: 'var(--gray)', marginTop: '4px' }}>
                      ({formatBytes(metrics.postgresql.database_size_bytes)})
                    </div>
                  )}
                </div>

                <div>
                  <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Connections</div>
                  <div style={{ display: 'flex', gap: '16px', alignItems: 'baseline' }}>
                    <div>
                      <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--success)' }}>
                        {metrics.postgresql.active_connections || 0}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)' }}>Active</div>
                    </div>
                    <div>
                      <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--gray)' }}>
                        {metrics.postgresql.total_connections || 0}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)' }}>Total</div>
                    </div>
                    <div>
                      <div style={{ fontSize: '20px', fontWeight: '700', color: 'var(--primary)' }}>
                        {metrics.postgresql.max_connections || 0}
                      </div>
                      <div style={{ fontSize: '12px', color: 'var(--gray)' }}>Max</div>
                    </div>
                  </div>
                  {metrics.postgresql.max_connections && metrics.postgresql.total_connections && (
                    <div style={{ width: '100%', height: '8px', background: 'var(--gray-lighter)', borderRadius: '4px', marginTop: '12px', overflow: 'hidden' }}>
                      <div style={{
                        width: `${Math.min((metrics.postgresql.total_connections / metrics.postgresql.max_connections) * 100, 100)}%`,
                        height: '100%',
                        background: 'linear-gradient(90deg, var(--success) 0%, var(--warning) 100%)',
                        transition: 'width 0.3s ease'
                      }}></div>
                    </div>
                  )}
                </div>

                {metrics.postgresql.cache_hit_ratio !== undefined && (
                  <div>
                    <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '4px', fontWeight: '600' }}>Cache Hit Ratio</div>
                    <div style={{ fontSize: '24px', fontWeight: '700', color: metrics.postgresql.cache_hit_ratio > 90 ? 'var(--success)' : metrics.postgresql.cache_hit_ratio > 70 ? 'var(--warning)' : 'var(--danger)' }}>
                      {metrics.postgresql.cache_hit_ratio.toFixed(2)}%
                    </div>
                    <div style={{ width: '100%', height: '8px', background: 'var(--gray-lighter)', borderRadius: '4px', marginTop: '8px', overflow: 'hidden' }}>
                      <div style={{
                        width: `${Math.min(metrics.postgresql.cache_hit_ratio, 100)}%`,
                        height: '100%',
                        background: metrics.postgresql.cache_hit_ratio > 90 
                          ? 'linear-gradient(90deg, var(--success) 0%, var(--secondary) 100%)'
                          : metrics.postgresql.cache_hit_ratio > 70
                          ? 'linear-gradient(90deg, var(--warning) 0%, #f59e0b 100%)'
                          : 'linear-gradient(90deg, var(--danger) 0%, #dc2626 100%)',
                        transition: 'width 0.3s ease'
                      }}></div>
                    </div>
                  </div>
                )}

                {metrics.postgresql.tables && metrics.postgresql.tables.length > 0 && (
                  <div>
                    <div style={{ fontSize: '12px', color: 'var(--gray)', marginBottom: '8px', fontWeight: '600' }}>Table Sizes</div>
                    <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid var(--gray-lighter)', borderRadius: 'var(--radius)', padding: '8px' }}>
                      {metrics.postgresql.tables.map((table, idx) => (
                        <div key={idx} style={{ padding: '8px', borderBottom: idx < metrics.postgresql.tables!.length - 1 ? '1px solid var(--gray-lighter)' : 'none', fontSize: '13px' }}>
                          <div style={{ fontWeight: '600', color: 'var(--dark)' }}>{table.tablename}</div>
                          <div style={{ color: 'var(--gray)', fontSize: '12px' }}>
                            {table.size} • {table.row_count?.toLocaleString() || 0} rows
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        ) : null}

        {metrics && (
          <div className="card" style={{ marginTop: '24px', background: 'rgba(255, 255, 255, 0.95)', backdropFilter: 'blur(10px)' }}>
            <div style={{ fontSize: '12px', color: 'var(--gray)', textAlign: 'center' }}>
              Last updated: {new Date((metrics.timestamp || 0) * 1000).toLocaleString()}
              {autoRefresh && <span style={{ marginLeft: '12px', color: 'var(--success)' }}>• Auto-refreshing every 5 seconds</span>}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}

export default Metrics

