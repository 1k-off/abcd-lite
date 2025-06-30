import React, { createContext, useContext, useState, useEffect } from "react"

const API_BASE_URL = import.meta.env.MODE === "development" ? "http://localhost:8900" : ""

interface AuthContextType {
  isAuthenticated: boolean
  login: (adminToken: string) => Promise<void>
  logout: () => Promise<void>
  error: string | null
  loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  // Always check auth status on mount by calling the protected endpoint
  useEffect(() => {
    const checkAuth = async () => {
      setLoading(true)
      try {
        const res = await fetch(`${API_BASE_URL}/api/auth/status`, { credentials: "include" })
        if (res.ok) {
          const data = await res.json()
          setIsAuthenticated(!!data.authenticated)
        } else {
          setIsAuthenticated(false)
        }
      } catch {
        setIsAuthenticated(false)
      } finally {
        setLoading(false)
      }
    }
    checkAuth()
  }, [])

  const login = async (adminToken: string) => {
    setError(null)
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ admin_token: adminToken }),
        credentials: "include",
      })
      if (!res.ok) {
        const data = await res.json().catch(() => ({}))
        throw new Error(data.message || data.error || "Login failed")
      }
      setIsAuthenticated(true)
    } catch (err: any) {
      setError(err.message || "Login failed")
      setIsAuthenticated(false)
    } finally {
      setLoading(false)
    }
  }

  const logout = async () => {
    setLoading(true)
    try {
      await fetch(`${API_BASE_URL}/logout`, {
        method: "POST",
        credentials: "include",
      })
    } catch {}
    setIsAuthenticated(false)
    setLoading(false)
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, login, logout, error, loading }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const ctx = useContext(AuthContext)
  if (!ctx) throw new Error("useAuth must be used within AuthProvider")
  return ctx
} 