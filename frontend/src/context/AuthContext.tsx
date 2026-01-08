import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react'

interface AuthContextType {
  token: string | null
  userId: number | null
  role: string | null
  login: (token: string, userId: number, role: string) => void
  logout: () => void
  isAuthenticated: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export const useAuth = () => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}

// Export AuthContext for direct access if needed
export { AuthContext }

interface AuthProviderProps {
  children: ReactNode
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [token, setToken] = useState<string | null>(null)
  const [userId, setUserId] = useState<number | null>(null)
  const [role, setRole] = useState<string | null>(null)

  useEffect(() => {
    // Load token from localStorage on mount
    const storedToken = localStorage.getItem('auth_token')
    const storedUserId = localStorage.getItem('user_id')
    const storedRole = localStorage.getItem('user_role')
    if (storedToken && storedUserId) {
      setToken(storedToken)
      setUserId(parseInt(storedUserId, 10))
      setRole(storedRole)
    }
  }, [])

  const login = (newToken: string, newUserId: number, newRole: string) => {
    setToken(newToken)
    setUserId(newUserId)
    setRole(newRole)
    localStorage.setItem('auth_token', newToken)
    localStorage.setItem('user_id', newUserId.toString())
    localStorage.setItem('user_role', newRole)
  }

  const logout = () => {
    setToken(null)
    setUserId(null)
    setRole(null)
    localStorage.removeItem('auth_token')
    localStorage.removeItem('user_id')
    localStorage.removeItem('user_role')
  }

  return (
    <AuthContext.Provider
      value={{
        token,
        userId,
        role,
        login,
        logout,
        isAuthenticated: !!token,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

