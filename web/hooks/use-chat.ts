"use client"

import { useState, useEffect } from "react"
import type { Message, Conversation } from "@/types/chat"
import { v4 as uuidv4 } from "uuid"

export function useChat() {
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [currentConversationId, setCurrentConversationId] = useState<string>("")
  const [isLoading, setIsLoading] = useState(false)

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
    setConversations((prev) => {
      const filtered = prev.filter((conv) => conv.id !== id)
      // If deleting current, pick next or start new
      if (id === currentConversationId) {
        if (filtered.length > 0) {
          setCurrentConversationId(filtered[0].id)
        } else {
          const newId = startNewConversation()
          setCurrentConversationId(newId)
        }
      }
      return filtered
    })
  }

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

  // Send a message (fix: always use up-to-date conversation ID)
  const sendMessage = async (
    content: string,
    attachment: { name: string; type: string; url: string } | null = null,
    scheduledDate: string | null = null,
  ) => {
    let convId = currentConversationId
    if (!convId) {
      convId = startNewConversation()
    }

    const userMessage: Message = {
      role: "user",
      content,
      timestamp: Date.now(),
      attachment,
      scheduledDate,
    }

    setConversations((prev) =>
      prev.map((conv) => {
        if (conv.id === convId) {
          return {
            ...conv,
            messages: [...conv.messages, userMessage],
            timestamp: Date.now(),
          }
        }
        return conv
      })
    )

    // Update the title if this is the first message
    const conv = conversations.find((c) => c.id === convId)
    if (!conv || conv.messages.length === 0) {
      updateConversationTitle(convId, content)
    }

    setIsLoading(true)

    try {
      const requestBody = {
        message: content,
        ...(attachment && { attachment: { name: attachment.name, type: attachment.type } }),
        ...(scheduledDate && { scheduledDate }),
      }
      const response = await fetch("/api/chat", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(requestBody),
      })
      if (!response.ok) throw new Error("Failed to send message")
      const data = await response.json()
      const botMessage: Message = {
        role: "assistant",
        content: data.response,
        timestamp: Date.now(),
      }
      setConversations((prev) =>
        prev.map((conv) => {
          if (conv.id === convId) {
            return {
              ...conv,
              messages: [...conv.messages, botMessage],
              timestamp: Date.now(),
            }
          }
          return conv
        })
      )
    } catch (error) {
      console.error("Error sending message:", error)
      const errorMessage: Message = {
        role: "assistant",
        content: "Sorry, I encountered an error. Please try again.",
        timestamp: Date.now(),
      }
      setConversations((prev) =>
        prev.map((conv) => {
          if (conv.id === convId) {
            return {
              ...conv,
              messages: [...conv.messages, errorMessage],
              timestamp: Date.now(),
            }
          }
          return conv
        })
      )
    } finally {
      setIsLoading(false)
    }
  }

  return {
    conversations,
    currentConversationId,
    currentMessages,
    isLoading,
    sendMessage,
    startNewConversation,
    loadConversation,
    deleteConversation,
  }
}
