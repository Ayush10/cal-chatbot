"use client"

import type React from "react"
import { useState, useRef, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Send } from "lucide-react"
import { ChatMessage } from "@/components/chat-message"
import { useChat } from "@/hooks/use-chat"
import { FileUpload } from "@/components/file-upload"
import { SimpleDatePicker } from "@/components/simple-date-picker"
import { format } from "date-fns"

export function ChatInterface() {
  const { currentMessages, isLoading, sendMessage, currentConversationId } = useChat()
  const [input, setInput] = useState("")
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [selectedDate, setSelectedDate] = useState<Date | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    scrollToBottom()
  }, [currentMessages])

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if ((!input.trim() && !selectedFile && !selectedDate) || isLoading) {
      return
    }

    let messageContent = input.trim()
    let attachmentData = null
    let dateString = null

    // Handle file attachment
    if (selectedFile) {
      // In a real app, you would upload the file to a server here
      // For this example, we'll just include the file name in the message
      attachmentData = {
        name: selectedFile.name,
        type: selectedFile.type,
        url: URL.createObjectURL(selectedFile), // This is temporary and will be revoked after use
      }

      if (!messageContent) {
        messageContent = `Attached file: ${selectedFile.name}`
      }
    }

    // Handle date selection
    if (selectedDate) {
      dateString = format(selectedDate, "yyyy-MM-dd'T'HH:mm:ss.SSSxxx")

      if (!messageContent) {
        messageContent = `Selected date: ${format(selectedDate, "PPP")}`
      }
    }

    await sendMessage(messageContent, attachmentData, dateString)

    // Reset form
    setInput("")
    setSelectedFile(null)
    setSelectedDate(null)
  }

  return (
    <Card className="w-full h-full">
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Avatar className="h-8 w-8">
            <AvatarImage src="/bot-avatar.png" alt="Bot" />
            <AvatarFallback>CB</AvatarFallback>
          </Avatar>
          <span>Cal Chatbot</span>
        </CardTitle>
      </CardHeader>
      <CardContent>
        <ScrollArea className="h-[500px] pr-4">
          {currentMessages.length === 0 ? (
            <div className="flex h-full items-center justify-center text-center p-8">
              <div className="space-y-2">
                <h3 className="text-lg font-semibold">Welcome to Cal Chatbot!</h3>
                <p className="text-sm text-muted-foreground">
                  Ask me about scheduling, availability, or managing your calendar.
                </p>
              </div>
            </div>
          ) : (
            <div className="space-y-4 pt-4">
              {currentMessages.map((message, index) => (
                <ChatMessage key={`${currentConversationId}-${index}`} message={message} />
              ))}
              <div ref={messagesEndRef} />
            </div>
          )}
        </ScrollArea>
      </CardContent>
      <CardFooter className="flex flex-col gap-2">
        {(selectedFile || selectedDate) && (
          <div className="flex flex-wrap gap-2 w-full">
            {selectedFile && <FileUpload onFileSelect={setSelectedFile} selectedFile={selectedFile} />}
            {selectedDate && <SimpleDatePicker onDateSelect={setSelectedDate} selectedDate={selectedDate} />}
          </div>
        )}
        <form onSubmit={handleSubmit} className="flex w-full gap-2">
          <div className="flex items-center gap-1">
            <FileUpload onFileSelect={setSelectedFile} selectedFile={null} />
            <SimpleDatePicker onDateSelect={setSelectedDate} selectedDate={null} />
          </div>
          <Input
            placeholder="Type your message..."
            value={input}
            onChange={(e) => setInput(e.target.value)}
            disabled={isLoading}
            className="flex-1"
          />
          <Button type="submit" size="icon" disabled={isLoading || (!input.trim() && !selectedFile && !selectedDate)}>
            <Send className="h-4 w-4" />
          </Button>
        </form>
      </CardFooter>
    </Card>
  )
}
