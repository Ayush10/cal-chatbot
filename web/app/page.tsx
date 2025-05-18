"use client"

import { useAuth } from "@/contexts/auth-context"
import { ChatInterface } from "@/components/chat-interface"
import { ChatHistory } from "@/components/chat-history"
import { EmailVerification } from "@/components/email-verification"
import { Button } from "@/components/ui/button"
import { LogOut } from "lucide-react"

export default function Home() {
  const { isAuthenticated, email, logout } = useAuth()

  return (
    <main className="flex min-h-screen flex-col p-4 md:p-8">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Cal Chatbot</h1>
        {isAuthenticated && (
          <div className="flex items-center gap-4">
            <span className="text-sm text-muted-foreground hidden md:inline-block">{email}</span>
            <Button variant="outline" size="sm" onClick={logout} className="flex items-center gap-2">
              <LogOut className="h-4 w-4" />
              <span>Logout</span>
            </Button>
          </div>
        )}
      </div>

      {isAuthenticated ? (
        <div className="flex flex-col lg:flex-row gap-4 w-full max-w-7xl mx-auto">
          <div className="w-full lg:w-80">
            <ChatHistory />
          </div>
          <div className="flex-1">
            <ChatInterface />
          </div>
        </div>
      ) : (
        <div className="flex items-center justify-center py-12">
          <EmailVerification />
        </div>
      )}
    </main>
  )
}
