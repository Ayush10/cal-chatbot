"use client"

import type React from "react"

import { useState } from "react"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Input } from "@/components/ui/input"
import { Button } from "@/components/ui/button"
import { useAuth } from "../contexts/auth-context"
import { CheckCircle, AlertCircle, Mail, KeyRound } from "lucide-react"
import { useToast } from "@/hooks/use-toast"

enum VerificationStep {
  EMAIL_INPUT = 0,
  CODE_INPUT = 1,
  SUCCESS = 2,
  ERROR = 3,
}

export function EmailVerification() {
  const { setAuthenticated } = useAuth()
  const { toast } = useToast()
  const [email, setEmail] = useState("")
  const [verificationCode, setVerificationCode] = useState("")
  const [step, setStep] = useState<VerificationStep>(VerificationStep.EMAIL_INPUT)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleEmailSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!email || !email.includes("@")) {
      setError("Please enter a valid email address")
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      // Call Go backend to request verification code
      const response = await fetch("/api/cal/request-verification-code", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to send verification code")
      }

      // Move to code input step
      setStep(VerificationStep.CODE_INPUT)
      toast({
        title: "Verification code sent",
        description: "Please check your email for the verification code",
      })
    } catch (error) {
      console.error("Error requesting verification code:", error)
      setError(error instanceof Error ? error.message : "Failed to send verification code")
    } finally {
      setIsLoading(false)
    }
  }

  const handleCodeSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (!verificationCode) {
      setError("Please enter the verification code")
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      // Call Go backend to verify code
      const response = await fetch("/api/cal/verify-email-code", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ email, code: verificationCode }),
      })

      if (!response.ok) {
        const data = await response.json()
        throw new Error(data.message || "Failed to verify code")
      }

      // Set authenticated state
      setAuthenticated(true, email)
      setStep(VerificationStep.SUCCESS)
      toast({
        title: "Verification successful",
        description: "You can now access your scheduled events",
      })
    } catch (error) {
      console.error("Error verifying code:", error)
      setError(error instanceof Error ? error.message : "Failed to verify code")
      setStep(VerificationStep.ERROR)
    } finally {
      setIsLoading(false)
    }
  }

  const handleRetry = () => {
    setVerificationCode("")
    setStep(VerificationStep.CODE_INPUT)
    setError(null)
  }

  const handleRestart = () => {
    setEmail("")
    setVerificationCode("")
    setStep(VerificationStep.EMAIL_INPUT)
    setError(null)
  }

  return (
    <Card className="w-full max-w-md mx-auto">
      <CardHeader>
        <CardTitle>Verify Your Email</CardTitle>
        <CardDescription>
          {step === VerificationStep.EMAIL_INPUT && "Enter your email to access the Cal Chatbot"}
          {step === VerificationStep.CODE_INPUT && "Enter the verification code sent to your email"}
          {step === VerificationStep.SUCCESS && "Email verified successfully"}
          {step === VerificationStep.ERROR && "Verification failed"}
        </CardDescription>
      </CardHeader>
      <CardContent>
        {step === VerificationStep.EMAIL_INPUT && (
          <form onSubmit={handleEmailSubmit} className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <Mail className="h-4 w-4 text-muted-foreground" />
                <label htmlFor="email" className="text-sm font-medium">
                  Email Address
                </label>
              </div>
              <Input
                id="email"
                type="email"
                placeholder="Enter your email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={isLoading}
                required
              />
            </div>
            {error && (
              <div className="flex items-center space-x-2 text-destructive text-sm">
                <AlertCircle className="h-4 w-4" />
                <span>{error}</span>
              </div>
            )}
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? "Sending..." : "Send Verification Code"}
            </Button>
          </form>
        )}

        {step === VerificationStep.CODE_INPUT && (
          <form onSubmit={handleCodeSubmit} className="space-y-4">
            <div className="space-y-2">
              <div className="flex items-center space-x-2">
                <KeyRound className="h-4 w-4 text-muted-foreground" />
                <label htmlFor="code" className="text-sm font-medium">
                  Verification Code
                </label>
              </div>
              <Input
                id="code"
                type="text"
                placeholder="Enter verification code"
                value={verificationCode}
                onChange={(e) => setVerificationCode(e.target.value)}
                disabled={isLoading}
                required
              />
              <p className="text-sm text-muted-foreground">
                A verification code has been sent to <span className="font-medium">{email}</span>
              </p>
            </div>
            {error && (
              <div className="flex items-center space-x-2 text-destructive text-sm">
                <AlertCircle className="h-4 w-4" />
                <span>{error}</span>
              </div>
            )}
            <div className="flex space-x-2">
              <Button type="button" variant="outline" className="flex-1" onClick={handleRestart} disabled={isLoading}>
                Change Email
              </Button>
              <Button type="submit" className="flex-1" disabled={isLoading}>
                {isLoading ? "Verifying..." : "Verify Code"}
              </Button>
            </div>
          </form>
        )}

        {step === VerificationStep.SUCCESS && (
          <div className="flex flex-col items-center justify-center py-4 space-y-4">
            <div className="rounded-full bg-green-100 p-3">
              <CheckCircle className="h-8 w-8 text-green-600" />
            </div>
            <p className="text-center">Your email has been verified successfully. You can now access the chatbot.</p>
          </div>
        )}

        {step === VerificationStep.ERROR && (
          <div className="flex flex-col items-center justify-center py-4 space-y-4">
            <div className="rounded-full bg-red-100 p-3">
              <AlertCircle className="h-8 w-8 text-red-600" />
            </div>
            <p className="text-center">Verification failed. Please try again.</p>
            <div className="flex space-x-2 w-full">
              <Button variant="outline" className="flex-1" onClick={handleRestart}>
                Change Email
              </Button>
              <Button className="flex-1" onClick={handleRetry}>
                Try Again
              </Button>
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  )
}
