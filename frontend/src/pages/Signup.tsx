import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { authService } from '../services/api'
import '../App.css'

const Signup = () => {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const navigate = useNavigate()

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    if (password !== confirmPassword) {
      setError('Passwords do not match')
      return
    }

    if (password.length < 4) {
      setError('Password must be at least 4 characters')
      return
    }

    setLoading(true)

    try {
      await authService.signup(username, password)
      // Redirect to login after successful signup
      navigate('/login', { state: { message: 'Account created successfully! Please login.' } })
    } catch (err: any) {
      if (err.code === 'ERR_NETWORK' || err.message?.includes('Network Error')) {
        setError('Cannot connect to server. Please make sure the backend is running on http://localhost:8080')
      } else if (err.response) {
        const errorMsg = err.response.data?.message || err.response.data?.error || 'Signup failed'
        setError(errorMsg)
      } else {
        setError(err.message || 'Signup failed. Please try again.')
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
          <div style={{ fontSize: '48px', marginBottom: '16px' }}>✨</div>
          <h1 style={{ 
            margin: 0,
            fontSize: '2rem',
            fontWeight: '700',
            background: 'linear-gradient(135deg, var(--success) 0%, var(--secondary) 100%)',
            WebkitBackgroundClip: 'text',
            WebkitTextFillColor: 'transparent',
            marginBottom: '8px'
          }}>
            Create Account
          </h1>
          <p style={{ color: 'var(--gray)', fontSize: '14px', margin: 0 }}>
            Join us and start managing your orders
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
              placeholder="Choose a username"
            />
          </div>

          <div style={{ marginBottom: '20px' }}>
            <label htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={4}
              placeholder="Enter password (min 4 characters)"
            />
          </div>

          <div style={{ marginBottom: '24px' }}>
            <label htmlFor="confirmPassword">Confirm Password</label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              placeholder="Confirm your password"
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
            className="btn btn-success"
            style={{ width: '100%', fontSize: '16px', padding: '14px' }}
          >
            {loading ? (
              <>
                <span className="loading" style={{ width: '16px', height: '16px', borderWidth: '2px' }}></span>
                <span>Creating account...</span>
              </>
            ) : (
              <>
                <span>🎉</span>
                <span>Sign Up</span>
              </>
            )}
          </button>
        </form>

        <p style={{ marginTop: '24px', textAlign: 'center', fontSize: '14px', color: 'var(--gray)' }}>
          Already have an account?{' '}
          <Link to="/login" style={{ 
            color: 'var(--primary)', 
            textDecoration: 'none',
            fontWeight: '600'
          }}>
            Login
          </Link>
        </p>
      </div>
    </div>
  )
}

export default Signup

