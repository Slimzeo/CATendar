export interface CalendarEvent {
  id: number
  title: string
  start: string
  end?: string
  color: string
  allDay: boolean
  description: string
  bold: boolean
  createdAt: string
  updatedAt: string
}

export interface EventInput {
  title: string
  start: string
  end?: string
  color: string
  allDay: boolean
  description: string
  bold: boolean
}

export interface DayCell {
  date: Date
  isCurrentMonth: boolean
  isToday: boolean
  isWeekend: boolean
}

export type AIProvider = 'codex' | 'claude'

export interface EmailAccount {
  id: string
  address: string
  username: string
  imapHost: string
  imapPort: number
  folder: string
  useTLS: boolean
  configured: boolean
  hasCredential: boolean
}

export interface EmailAccountInput {
  address: string
  username: string
  imapHost: string
  imapPort: number
  folder: string
  useTLS: boolean
  secret: string
}

export interface EmailMessage {
  id: string
  messageId: string
  subject: string
  sender: string
  receivedAt: string
  sourceUrl: string
  text: string
}

export interface AISettings {
  provider: AIProvider
  codexPath: string
  claudePath: string
}

export interface AgentAvailability {
  provider: AIProvider
  installed: boolean
  path: string
  version?: string
  error?: string
}

export interface AISyncSetup {
  email: EmailAccount
  agent: {
    settings: AISettings
    agents: AgentAvailability[]
  }
}

export type AISyncState = 'idle' | 'running' | 'succeeded' | 'failed' | 'cancelled'

export interface AISyncStatus {
  runId?: string
  state: AISyncState
  stage?: string
  message?: string
  provider?: AIProvider
  emailCount: number
  created: number
  skipped: number
  rejected: number
  conflicts: number
  startedAt?: string
  finishedAt?: string
}
