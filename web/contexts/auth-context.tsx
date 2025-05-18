"use client"

import type React from "react"
import { createContext, useContext, useState, useEffect } from "react"

interface AuthContextType {
  isAuthenticated: boolean
  email: string | null
  setAuthenticated: (value: boolean, email?: string) => void
  logout: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuthenticated, setIsAuthenticated] = useState(false)
  const [email, setEmail] = useState<string | null>(null)

  // Check localStorage on initial load
  useEffect(() => {
    const storedAuth = localStorage.getItem("auth")
    if (storedAuth) {
      try {
        const authData = JSON.parse(storedAuth)
        if (authData.isAuthenticated && authData.email) {
          setIsAuthenticated(true)
          setEmail(authData.email)
        }
      } catch (error) {
        console.error("Error parsing auth data:", error)
        localStorage.removeItem("auth")
      }
    }
  }, [])

  const setAuthenticated = (value: boolean, email?: string) => {
    setIsAuthenticated(value)
    if (email) {
      setEmail(email)
      // Store in localStorage
      localStorage.setItem("auth", JSON.stringify({ isAuthenticated: value, email }))
    } else if (!value) {
      setEmail(null)
      localStorage.removeItem("auth")
    }
  }

  const logout = () => {
    setIsAuthenticated(false)
    setEmail(null)
    localStorage.removeItem("auth")
  }

  return (
    <AuthContext.Provider value={{ isAuthenticated, email, setAuthenticated, logout }}>{children}</AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}
