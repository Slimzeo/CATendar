import type {
  AISettings,
  AISyncSetup,
  AISyncStatus,
  CalendarEvent,
  EmailAccount,
  EmailAccountInput,
  EmailMessage,
  EventInput
} from '@/types'

// Wails generates individual functions in frontend/wailsjs/go/main/CalendarApp
import { GetEvents, CreateEvent, UpdateEvent, DeleteEvent } from '../../wailsjs/go/main/CalendarApp'
import {
  CancelAISync,
  GetAISyncSetup,
  GetAISyncStatus,
  GetSourceEmail,
  SaveAISettings,
  SaveEmailAccount,
  StartAISync,
  TestEmailAccount
} from '../../wailsjs/go/main/AISyncApp'

export const eventApi = {
  getEvents: async (start: string, end: string): Promise<CalendarEvent[]> => {
    if (!('go' in window)) return []
    return GetEvents(start, end)
  },

  createEvent: async (input: EventInput): Promise<CalendarEvent> => {
    return CreateEvent(input)
  },

  updateEvent: async (id: number, input: EventInput): Promise<CalendarEvent> => {
    return UpdateEvent(id, input)
  },

  deleteEvent: async (id: number): Promise<void> => {
    await DeleteEvent(id)
  }
}

const browserSetup: AISyncSetup = {
  email: {
    id: 'primary',
    address: '',
    username: '',
    imapHost: '',
    imapPort: 993,
    folder: 'INBOX',
    useTLS: true,
    configured: false,
    hasCredential: false
  },
  agent: {
    settings: { provider: 'codex', codexPath: '', claudePath: '' },
    agents: [
      { provider: 'codex', installed: false, path: '', error: 'Desktop runtime required' },
      { provider: 'claude', installed: false, path: '', error: 'Desktop runtime required' }
    ]
  }
}

export const aiSyncApi = {
  getSetup: async (): Promise<AISyncSetup> => {
    if (!('go' in window)) return browserSetup
    return GetAISyncSetup() as Promise<AISyncSetup>
  },

  saveEmail: async (input: EmailAccountInput): Promise<EmailAccount> => {
    return SaveEmailAccount(input) as Promise<EmailAccount>
  },

  testEmail: async (input: EmailAccountInput): Promise<void> => {
    await TestEmailAccount(input)
  },

  getSourceEmail: async (sourceURL: string): Promise<EmailMessage> => {
    return GetSourceEmail(sourceURL) as Promise<EmailMessage>
  },

  saveSettings: async (settings: AISettings): Promise<void> => {
    await SaveAISettings(settings)
  },

  start: async (): Promise<AISyncStatus> => {
    return StartAISync() as Promise<AISyncStatus>
  },

  status: async (): Promise<AISyncStatus> => {
    if (!('go' in window)) {
      return { state: 'idle', emailCount: 0, created: 0, skipped: 0, rejected: 0, conflicts: 0 }
    }
    return GetAISyncStatus() as Promise<AISyncStatus>
  },

  cancel: async (): Promise<boolean> => {
    return CancelAISync()
  }
}
