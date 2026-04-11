import axios from 'axios'
import type { CalendarEvent, EventInput } from '@/types'

const apiBaseURL = import.meta.env.DEV ? '/api' : 'http://localhost:18082/api'

const api = axios.create({
  baseURL: apiBaseURL
})

export const eventApi = {
  getEvents: async (start?: string, end?: string): Promise<CalendarEvent[]> => {
    const params = start && end ? { start, end } : {}
    const response = await api.get<CalendarEvent[]>('/events', { params })
    return response.data
  },

  createEvent: async (input: EventInput): Promise<CalendarEvent> => {
    const response = await api.post<CalendarEvent>('/events', input)
    return response.data
  },

  updateEvent: async (id: number, input: EventInput): Promise<CalendarEvent> => {
    const response = await api.put<CalendarEvent>(`/events/${id}`, input)
    return response.data
  },

  deleteEvent: async (id: number): Promise<void> => {
    await api.delete(`/events/${id}`)
  }
}
