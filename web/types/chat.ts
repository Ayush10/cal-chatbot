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
  read?: boolean
  booking?: Record<string, any>
  listEvents?: Record<string, any>
}

export interface Conversation {
  id: string
  title: string
  messages: Message[]
  timestamp: number
  lastReadTimestamp?: number
}
