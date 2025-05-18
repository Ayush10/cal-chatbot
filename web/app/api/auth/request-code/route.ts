import { NextResponse } from "next/server"

// In a real application, you would store these codes in a database with expiration
const verificationCodes: Record<string, string> = {}

export async function POST(request: Request) {
  try {
    const { email } = await request.json()

    if (!email || !email.includes("@")) {
      return NextResponse.json({ message: "Invalid email address" }, { status: 400 })
    }

    // Generate a random 6-digit code
    const code = Math.floor(100000 + Math.random() * 900000).toString()

    // Store the code (in a real app, you would send an actual email)
    verificationCodes[email] = code

    // For demo purposes, we'll log the code to the console
    console.log(`Verification code for ${email}: ${code}`)

    // In a real application, you would send an email here
    // await sendEmail(email, `Your verification code is: ${code}`);

    return NextResponse.json({ message: "Verification code sent" })
  } catch (error) {
    console.error("Error requesting verification code:", error)
    return NextResponse.json({ message: "Failed to send verification code" }, { status: 500 })
  }
}
