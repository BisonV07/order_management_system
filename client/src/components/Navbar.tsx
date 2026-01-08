import { useAuth } from '../context/AuthContext'
import { useNavigate } from 'react-router-dom'
import '../App.css'

const Navbar = () => {
  const { isAuthenticated, logout, role } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    logout()
    navigate('/login')
  }

  if (!isAuthenticated) return null

  return (
    <nav style={{
      background: 'rgba(255, 255, 255, 0.95)',
      backdropFilter: 'blur(20px)',
      color: 'var(--dark)',
      padding: '16px 24px',
      display: 'flex',
      justifyContent: 'space-between',
      alignItems: 'center',
      boxShadow: 'var(--shadow-md)',
      borderBottom: '1px solid rgba(0, 0, 0, 0.05)',
      position: 'sticky',
      top: 0,
      zIndex: 1000
    }}>
      <div style={{ display: 'flex', gap: '24px', alignItems: 'center' }}>
        <h2 style={{ 
          margin: 0, 
          fontSize: '20px',
          fontWeight: '700',
          background: 'linear-gradient(135deg, var(--primary) 0%, var(--primary-dark) 100%)',
          WebkitBackgroundClip: 'text',
          WebkitTextFillColor: 'transparent'
        }}>
          🛒 Order Management System
        </h2>
        {role === 'admin' || role === 'ADMIN' ? (
          <span className="badge badge-danger">
            👑 ADMIN
          </span>
        ) : (
          <span className="badge badge-info" style={{ background: 'rgba(59, 130, 246, 0.1)', color: 'var(--info)' }}>
            👤 User
          </span>
        )}
        <div style={{ display: 'flex', gap: '8px', marginLeft: '8px' }}>
          <a 
            href="/" 
            style={{ 
              color: 'var(--dark)', 
              textDecoration: 'none', 
              padding: '8px 16px',
              borderRadius: 'var(--radius)',
              fontWeight: '500',
              fontSize: '14px',
              transition: 'var(--transition)'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(99, 102, 241, 0.1)'
              e.currentTarget.style.color = 'var(--primary)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent'
              e.currentTarget.style.color = 'var(--dark)'
            }}
            onClick={(e) => { e.preventDefault(); navigate('/') }}
          >
            📊 Dashboard
          </a>
          <a 
            href="/orders" 
            style={{ 
              color: 'var(--dark)', 
              textDecoration: 'none', 
              padding: '8px 16px',
              borderRadius: 'var(--radius)',
              fontWeight: '500',
              fontSize: '14px',
              transition: 'var(--transition)'
            }}
            onMouseEnter={(e) => {
              e.currentTarget.style.background = 'rgba(99, 102, 241, 0.1)'
              e.currentTarget.style.color = 'var(--primary)'
            }}
            onMouseLeave={(e) => {
              e.currentTarget.style.background = 'transparent'
              e.currentTarget.style.color = 'var(--dark)'
            }}
            onClick={(e) => { e.preventDefault(); navigate('/orders') }}
          >
            📦 Orders
          </a>
          {(role === 'admin' || role === 'ADMIN') && (
            <>
              <a 
                href="/admin" 
                style={{ 
                  color: 'var(--primary)', 
                  textDecoration: 'none', 
                  padding: '8px 16px',
                  borderRadius: 'var(--radius)',
                  fontWeight: '600',
                  fontSize: '14px',
                  background: 'rgba(99, 102, 241, 0.1)',
                  transition: 'var(--transition)'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(99, 102, 241, 0.2)'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'rgba(99, 102, 241, 0.1)'
                }}
                onClick={(e) => { e.preventDefault(); navigate('/admin') }}
              >
                ⚙️ Admin Panel
              </a>
              <a 
                href="/metrics" 
                style={{ 
                  color: 'var(--primary)', 
                  textDecoration: 'none', 
                  padding: '8px 16px',
                  borderRadius: 'var(--radius)',
                  fontWeight: '600',
                  fontSize: '14px',
                  background: 'rgba(99, 102, 241, 0.1)',
                  transition: 'var(--transition)'
                }}
                onMouseEnter={(e) => {
                  e.currentTarget.style.background = 'rgba(99, 102, 241, 0.2)'
                }}
                onMouseLeave={(e) => {
                  e.currentTarget.style.background = 'rgba(99, 102, 241, 0.1)'
                }}
                onClick={(e) => { e.preventDefault(); navigate('/metrics') }}
              >
                📊 Metrics
              </a>
            </>
          )}
        </div>
      </div>
      <button
        onClick={handleLogout}
        className="btn btn-danger"
        style={{ fontSize: '14px', padding: '10px 20px' }}
      >
        🚪 Logout
      </button>
    </nav>
  )
}

export default Navbar

