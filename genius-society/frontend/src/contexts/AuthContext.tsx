import { createContext, useContext, useEffect, useState, type ReactNode } from 'react'
import { api, type User, type ApiResponse } from '../services/api'

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (login: string, password: string) => Promise<User>
  logout: () => Promise<void>
  refreshUser: () => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)

  const refreshUser = async () => {
    try {
      const res = await api.get<ApiResponse<User>>('/auth/me')
      if (res.data.success && res.data.data) {
        setUser(res.data.data)
      } else {
        setUser(null)
      }
    } catch {
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    refreshUser()
  }, [])

  const login = async (loginId: string, password: string): Promise<User> => {
    const res = await api.post<ApiResponse<User>>('/auth/login', {
      login: loginId,
      password,
    })

    if (!res.data.success || !res.data.data) {
      throw new Error(res.data.error?.message || 'Login gagal')
    }

    const authenticatedUser = res.data.data
    setUser(authenticatedUser)
    return authenticatedUser
  }

  const logout = async () => {
    try {
      await api.post('/auth/logout')
    } finally {
      setUser(null)
    }
  }

  return (
    <AuthContext.Provider value={{ user, loading, login, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
