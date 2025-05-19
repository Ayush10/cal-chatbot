"use client"

import { useState, useEffect } from "react"
import type { Message, Conversation } from "@/types/chat"
import { v4 as uuidv4 } from "uuid"
import { parseISO, isValid } from "date-fns"
import { useAuth } from "@/contexts/auth-context"

// Helper: Detect booking intent
function isBookingIntent(content: string): boolean {
  return /book( a| an)? meeting|schedule|create.*booking|random details/i.test(content)
}

// Helper: Detect list events intent
function isListEventsIntent(content: string): boolean {
  return /show.*events|list.*events|scheduled events/i.test(content)
}

// Helper: Extract email from user message
function extractEmail(content: string): string | null {
  const match = content.match(/[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}/)
  return match ? match[0] : null
}

// Helper: Try to extract booking details from user message
function extractBookingDetails(content: string): any | null {
  // Very basic extraction for demo; in production, use NLP or a form
  if (/random details/i.test(content)) {
    const now = new Date()
    const tomorrow = new Date(now.getTime() + 24 * 60 * 60 * 1000)
    const startTime = new Date(tomorrow.setHours(10, 0, 0, 0))
    const endTime = new Date(tomorrow.setHours(11, 0, 0, 0))
    return {
      eventTypeId: 0,
      startTime: startTime.toISOString(),
      endTime: endTime.toISOString(),
      name: "Random User",
      email: `randomuser+${Date.now()}@example.com`,
      notes: "Random details for the meeting",
    }
  }
  // Add more extraction logic as needed
  return null
}

export function useChat() {
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [currentConversationId, setCurrentConversationId] = useState<string>("")
  const [isLoading, setIsLoading] = useState(false)
  const { email: userEmail } = useAuth()

  // Initialize or load from localStorage
  useEffect(() => {
    const savedConversations = localStorage.getItem("chatConversations")
    const savedCurrentId = localStorage.getItem("currentConversationId")

    if (savedConversations) {
      setConversations(JSON.parse(savedConversations))
    }

    if (savedCurrentId) {
      setCurrentConversationId(savedCurrentId)
    } else {
      startNewConversation()
    }
  }, [])

  // Save to localStorage whenever conversations change
  useEffect(() => {
    if (conversations.length > 0) {
      localStorage.setItem("chatConversations", JSON.stringify(conversations))
    }
    if (currentConversationId) {
      localStorage.setItem("currentConversationId", currentConversationId)
    }
  }, [conversations, currentConversationId])

  // Get current conversation messages
  const currentMessages = conversations.find((conv) => conv.id === currentConversationId)?.messages || []

  // Start a new conversation and return its ID
  const startNewConversation = () => {
    const newId = uuidv4()
    const newConversation: Conversation = {
      id: newId,
      title: "New Conversation",
      messages: [],
      timestamp: Date.now(),
    }
    setConversations((prev) => [newConversation, ...prev])
    setCurrentConversationId(newId)
    return newId
  }

  // Load an existing conversation, fallback to first if not found
  const loadConversation = (id: string) => {
    const exists = conversations.some((conv) => conv.id === id)
    if (exists) {
      setCurrentConversationId(id)
    } else if (conversations.length > 0) {
      setCurrentConversationId(conversations[0].id)
    } else {
      const newId = startNewConversation()
      setCurrentConversationId(newId)
    }
  }

  // Delete a conversation and always select a valid one
  const deleteConversation = (id: string) => {
    setConversations((prev) => prev.filter((conv) => conv.id !== id))
    if (id === currentConversationId) {
      setCurrentConversationId("")
    }
  }

  // Ensure a new conversation is started if all are deleted
  useEffect(() => {
    if (conversations.length === 0) {
      startNewConversation()
    }
  }, [conversations])

  // Update conversation title based on first message
  const updateConversationTitle = (id: string, message: string) => {
    setConversations((prev) =>
      prev.map((conv) => {
        if (conv.id === id && (conv.title === "New Conversation" || !conv.title)) {
          const title = message.length > 30 ? `${message.substring(0, 30)}...` : message
          return { ...conv, title }
        }
        return conv
      })
    )
  }

  // Send a message
  const sendMessage = async (
    content: string,
    attachment: { name: string; type: string; url: string } | null = null,
    scheduledDate: string | null = null,
  ) => {
    let convId = currentConversationId
    if (!convId) {
      convId = startNewConversation()
    }

    // Build the user message
    const userMessage: Message = {
      role: "user",
      content,
      timestamp: Date.now(),
      attachment,
      scheduledDate,
      read: true,
    }

    // Attach only one of booking or listEvents
    if (isBookingIntent(content)) {
      userMessage.booking = extractBookingDetails(content)
    } else if (isListEventsIntent(content)) {
      // Use explicit email if present, else fallback to authenticated email
      const explicitEmail = extractEmail(content)
      const emailToUse = explicitEmail || userEmail
      if (!emailToUse) {
        // Prompt for email if not available
        setConversations((prev) =>
          prev.map((conv) =>
            conv.id === convId
              ? {
                  ...conv,
                  messages: [
                    ...conv.messages,
                    {
                      role: "assistant",
                      content: "Please provide your email address to view your scheduled events.",
                      timestamp: Date.now(),
                      read: false,
                    },
                  ],
                }
              : conv
          )
        )
        return
      }
      userMessage.listEvents = { email: emailToUse }
    }

    // Add user message to conversation
    setConversations((prev) =>
      prev.map((conv) =>
        conv.id === convId
          ? { ...conv, messages: [...conv.messages, userMessage] }
          : conv
      )
    )
    setIsLoading(true)

    // Prepare payload for backend (strip frontend-only fields)
    const payloadMessages = conversations.find((c) => c.id === convId)?.messages.concat(userMessage) || [userMessage]
    const sanitizedMessages = payloadMessages.map(({ read, ...rest }) => rest)

    try {
      const response = await fetch("/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ messages: sanitizedMessages }),
      })
      const data = await response.json()
      const botMessage: Message = {
        role: "assistant",
        content: data.response || data.message || JSON.stringify(data),
        timestamp: Date.now(),
        read: false,
      }
      setConversations((prev) =>
        prev.map((conv) =>
          conv.id === convId
            ? { ...conv, messages: [...conv.messages, botMessage] }
            : conv
        )
      )
    } catch (err) {
      setConversations((prev) =>
        prev.map((conv) =>
          conv.id === convId
            ? {
                ...conv,
                messages: [
                  ...conv.messages,
                  {
                    role: "assistant",
                    content: "Sorry, I encountered an error. Please try again.",
                    timestamp: Date.now(),
                    read: false,
                  },
                ],
              }
            : conv
        )
      )
    } finally {
      setIsLoading(false)
    }
  }

  // Get unread count for each conversation
  const getUnreadCount = (conv: Conversation) =>
    conv.messages.filter((msg) => msg.role === "assistant" && !msg.read).length

  // Mark all messages in the current conversation as read
  const markAllAsRead = () => {
    setConversations((prev) =>
      prev.map((conv) => {
        if (conv.id === currentConversationId) {
          const now = Date.now()
          return {
            ...conv,
            messages: conv.messages.map((msg) => ({ ...msg, read: true })),
            lastReadTimestamp: now,
          }
        }
        return conv
      })
    )
  }

  // Listen for tab visibility and mark as read if active
  useEffect(() => {
    const handleVisibility = () => {
      if (document.visibilityState === "visible") {
        markAllAsRead()
      }
    }
    document.addEventListener("visibilitychange", handleVisibility)
    // Also mark as read on mount and when conversation changes
    markAllAsRead()
    return () => {
      document.removeEventListener("visibilitychange", handleVisibility)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentConversationId])

  return {
    conversations,
    currentConversationId,
    currentMessages,
    isLoading,
    sendMessage,
    startNewConversation,
    loadConversation,
    deleteConversation,
    markAllAsRead,
    getUnreadCount,
  }
}
