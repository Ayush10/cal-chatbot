"use client"

import { cn } from "@/lib/utils"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import { Calendar, FileText } from "lucide-react"
import { format, parseISO } from "date-fns"
import type { Message } from "@/types/chat"

interface ChatMessageProps {
  message: Message
}

export function ChatMessage({ message }: ChatMessageProps) {
  const isUser = message.role === "user"

  return (
    <div className={cn("flex items-start gap-3 text-sm", isUser ? "flex-row-reverse" : "")}>
      <Avatar className="h-8 w-8 mt-0.5">
        {isUser ? (
          <>
            <AvatarImage src="/user-avatar.png" alt="User" />
            <AvatarFallback>U</AvatarFallback>
          </>
        ) : (
          <>
            <AvatarImage src="/bot-avatar.png" alt="Bot" />
            <AvatarFallback>CB</AvatarFallback>
          </>
        )}
      </Avatar>
      <div
        className={cn(
          "rounded-lg px-4 py-2 max-w-[80%] space-y-2",
          isUser ? "bg-primary text-primary-foreground" : "bg-muted",
        )}
      >
        <div>{message.content}</div>

        {message.attachment && (
          <div
            className={cn(
              "flex items-center gap-2 p-2 rounded-md text-xs",
              isUser ? "bg-primary-foreground/10" : "bg-background",
            )}
          >
            <FileText className="h-4 w-4" />
            <a
              href={message.attachment.url}
              target="_blank"
              rel="noopener noreferrer"
              className="underline underline-offset-2"
            >
              {message.attachment.name}
            </a>
          </div>
        )}

        {message.scheduledDate && (
          <div
            className={cn(
              "flex items-center gap-2 p-2 rounded-md text-xs",
              isUser ? "bg-primary-foreground/10" : "bg-background",
            )}
          >
            <Calendar className="h-4 w-4" />
            <span>{format(parseISO(message.scheduledDate), "PPP 'at' p")}</span>
          </div>
        )}

        {message.role === "assistant" && message.read && (
          <span className="ml-2 text-xs text-green-500 align-bottom" title="Read">✔</span>
        )}
      </div>
    </div>
  )
}
