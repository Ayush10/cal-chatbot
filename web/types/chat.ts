export interface Message {
  role: "user" | "assistant"
  content: string
  timestamp?: number
  attachment?: {
    name: string
    type: string
    url: string
  } | null
  scheduledDate?: string | null
}

export interface Conversation {
  id: string
  title: string
  messages: Message[]
  timestamp: number
}
