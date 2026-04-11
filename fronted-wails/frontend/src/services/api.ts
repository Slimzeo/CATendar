import type { CalendarEvent, EventInput } from '@/types'

// Wails generates individual functions in frontend/wailsjs/go/main/CalendarApp
import { GetEvents, CreateEvent, UpdateEvent, DeleteEvent } from '../../wailsjs/go/main/CalendarApp'

export const eventApi = {
  getEvents: async (start?: string, end?: string): Promise<CalendarEvent[]> => {
    return await GetEvents(start || '', end || '')
  },

  createEvent: async (input: EventInput): Promise<CalendarEvent> => {
    return await CreateEvent(input)
  },

  updateEvent: async (id: number, input: EventInput): Promise<CalendarEvent> => {
    return await UpdateEvent(id, input)
  },

  deleteEvent: async (id: number): Promise<void> => {
    await DeleteEvent(id)
  }
}