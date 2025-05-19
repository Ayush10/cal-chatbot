import { NextResponse } from "next/server"

export async function POST(request: Request) {
  try {
    const { messages } = await request.json()
    // Remove frontend-only fields (like 'read') from each message
    const sanitizedMessages = Array.isArray(messages)
      ? messages.map(({ read, ...rest }) => rest)
      : []

    // Call the Go backend API
    const response = await fetch("http://localhost:8080/api/chat", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ messages: sanitizedMessages }),
    })

    if (!response.ok) {
      throw new Error("Failed to get response from backend")
    }

    const data = await response.json()
    console.log("Backend returned:", data)

    // Always return a 'response' field for the frontend
    return NextResponse.json({ response: data.response || data.message || JSON.stringify(data) })
  } catch (error) {
    console.error("Error in chat API:", error)
    return NextResponse.json({ error: "Failed to process chat request" }, { status: 500 })
  }
}
