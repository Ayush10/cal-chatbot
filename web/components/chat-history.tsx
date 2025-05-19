"use client"

import { useChat } from "@/hooks/use-chat"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { History, Plus, Trash2 } from "lucide-react"
import { cn } from "@/lib/utils"
import { formatDistanceToNow } from "date-fns"

export function ChatHistory() {
  const { conversations, currentConversationId, startNewConversation, loadConversation, deleteConversation, getUnreadCount } = useChat()

  return (
    <Card className="h-full">
      <CardHeader className="px-4 py-3">
        <CardTitle className="text-base flex items-center gap-2">
          <History className="h-4 w-4" />
          <span>Chat History</span>
        </CardTitle>
      </CardHeader>
      <CardContent className="p-0">
        <div className="p-2">
          <Button variant="outline" className="w-full justify-start gap-2" onClick={startNewConversation}>
            <Plus className="h-4 w-4" />
            New Chat
          </Button>
        </div>
        <ScrollArea className="h-[500px]">
          <div className="space-y-1 p-2">
            {conversations.length === 0 ? (
              <p className="text-sm text-muted-foreground text-center py-4">No conversation history</p>
            ) : (
              conversations.map((conversation) => {
                const unread = getUnreadCount(conversation)
                return (
                  <div
                    key={conversation.id}
                    className={cn(
                      "flex items-center justify-between group rounded-md px-3 py-2 text-sm hover:bg-accent cursor-pointer",
                      conversation.id === currentConversationId && "bg-accent",
                    )}
                    onClick={() => loadConversation(conversation.id)}
                  >
                    <div className="truncate">
                      <div className="font-medium truncate flex items-center gap-2">
                        {conversation.title || "New Conversation"}
                        {unread > 0 && (
                          <span className="ml-2 inline-block min-w-[18px] px-1 py-0.5 rounded-full bg-green-500 text-white text-xs text-center">
                            {unread}
                          </span>
                        )}
                      </div>
                      <div className="text-xs text-muted-foreground">
                        {formatDistanceToNow(new Date(conversation.timestamp), { addSuffix: true })}
                      </div>
                    </div>
                    <Button
                      variant="ghost"
                      size="icon"
                      className="h-6 w-6 opacity-0 group-hover:opacity-100"
                      onClick={(e) => {
                        e.stopPropagation()
                        deleteConversation(conversation.id)
                      }}
                    >
                      <Trash2 className="h-3 w-3" />
                      <span className="sr-only">Delete</span>
                    </Button>
                  </div>
                )
              })
            )}
          </div>
        </ScrollArea>
      </CardContent>
    </Card>
  )
}
