import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import { authService } from '../services/api'
import '../App.css'

const Login = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const { login } = useAuth()
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)

    try {
      // Use the API service for better error handling
      const data = await authService.login(username, password)

      // Store token, user ID, and role
      login(data.token, data.user_id, data.role)

      // Redirect to dashboard
      navigate('/')
    } catch (err: any) {
      // Provide more specific error messages
      if (err.code === 'ERR_NETWORK' || err.message?.includes('Network Error') || err.message?.includes('Failed to fetch')) {
        setError('Cannot connect to server. Please make sure the backend is running on http://localhost:8080')
      } else if (err.response) {
        // Server responded with error
        const errorMsg = err.response.data?.message || err.response.data?.error || 'Login failed'
        setError(errorMsg)
      } else {
        setError(err.message || 'Login failed. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="app" style={{ display: 'flex', alignItems: 'center', justifyContent: 'center', minHeight: '100vh', padding: '20px' }}>
      <div className="card" style={{
        width: '100%',
        maxWidth: '440px',
        background: 'rgba(255, 255, 255, 0.95)',
        backdropFilter: 'blur(20px)',
        boxShadow: 'var(--shadow-lg)'
      }}>
        <div style={{ textAlign: 'center', marginBottom: '32px' }}>
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>🔐</div>
          <h1 style={{ 
            margin: 0,
            fontSize: '2rem',
            fontWeight: '700',
            background: 'linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            marginBottom: '8px'
          }}>
            Welcome Back
          </h1>
          <p style={{ color: 'var(--gray)', fontSize: '14px', margin: 0 }}>
            Sign in to your account to continue
          </p>
        </div>
        
        <form onSubmit={handleSubmit}>
          <div style={{ marginBottom: '20px' }}>
            <label htmlFor="username">Username</label>
            <input
              id="username"
              type="text"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              required
              placeholder="Enter your username"
            />
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              placeholder="Enter your password"
            />
          </div>

          {error && (
            <div className="alert alert-error" style={{ marginBottom: '20px' }}>
              <span>⚠️</span>
              <span>{error}</span>
            </div>
          )}

          <button
            type="submit"
            disabled={loading}
            className="btn btn-primary"
            style={{ width: '100%', fontSize: '16px', padding: '14px' }}
          >
            {loading ? (
              <>
                <span className="loading" style={{ width: '16px', height: '16px', borderWidth: '2px' }}></span>
                <span>Logging in...</span>
              </>
            ) : (
              <>
                <span>🚀</span>
                <span>Login</span>
              </>
            )}
          </button>
        </form>

        <div style={{ marginTop: '24px', textAlign: 'center' }}>
          <p style={{ fontSize: '14px', color: 'var(--gray)', margin: 0 }}>
            Don't have an account?{' '}
            <Link to="/signup" style={{ 
              color: 'var(--primary)', 
              textDecoration: 'none',
              fontWeight: '600'
            }}>
              Sign Up
            </Link>
          </p>
        </div>
      </div>
    </div>
  )
}

export default Login

