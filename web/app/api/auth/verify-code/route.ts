import { NextResponse } from "next/server"

// This should be imported from a shared location in a real app
const verificationCodes: Record<string, string> = {}

export async function POST(request: Request) {
  try {
    const { email, code } = await request.json()

    if (!email || !code) {
      return NextResponse.json({ message: "Email and code are required" }, { status: 400 })
    }

    // Check if the code matches
    const storedCode = verificationCodes[email]

    if (!storedCode) {
      return NextResponse.json({ message: "No verification code found for this email" }, { status: 400 })
    }

    if (storedCode !== code) {
      return NextResponse.json({ message: "Invalid verification code" }, { status: 400 })
    }

    // Clear the code after successful verification
    delete verificationCodes[email]

    return NextResponse.json({ message: "Verification successful" })
  } catch (error) {
    console.error("Error verifying code:", error)
    return NextResponse.json({ message: "Failed to verify code" }, { status: 500 })
  }
}
