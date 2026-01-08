import { useEffect } from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider, useAuth } from './context/AuthContext'
import Dashboard from './pages/Dashboard'
import Orders from './pages/Orders'
import Login from './pages/Login'
import Signup from './pages/Signup'
import AdminPanel from './pages/AdminPanel'
import Metrics from './pages/Metrics'
import Navbar from './components/Navbar'
import ErrorBoundary from './components/ErrorBoundary'
// import DebugConsole from './components/DebugConsole' // Temporarily disabled
import './App.css'

// Protected Route component - requires authentication
const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuth()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" replace />
}

// Admin Route component - requires admin role
const AdminRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated, role } = useAuth()
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }
  if (role !== 'admin' && role !== 'ADMIN') {
    return <Navigate to="/" replace />
  }
  return <>{children}</>
}

function AppRoutes() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<Signup />} />
        <Route path="/" element={<ProtectedRoute><Dashboard /></ProtectedRoute>} />
        <Route path="/orders" element={<ProtectedRoute><Orders /></ProtectedRoute>} />
        <Route path="/admin" element={<AdminRoute><AdminPanel /></AdminRoute>} />
        <Route path="/metrics" element={<AdminRoute><Metrics /></AdminRoute>} />
      </Routes>
      {/* Debug console - temporarily disabled */}
      {/* <DebugConsole /> */}
    </>
  )
}

function App() {
  // Minimal logging on app load
  useEffect(() => {
    // Only log if there's an issue with API URL
    const apiUrl = import.meta.env.VITE_API_BASE_URL || '/api/v1'
    if (apiUrl.startsWith('http')) {
      console.warn('⚠️ API URL is using full URL, should use proxy')
    }
  }, [])

  return (
    <ErrorBoundary>
      <AuthProvider>
        <Router>
          <AppRoutes />
        </Router>
      </AuthProvider>
    </ErrorBoundary>
  )
}

export default App

