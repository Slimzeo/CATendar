import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { format } from 'date-fns'
import type { CalendarEvent, EventInput } from '@/types'
import { eventApi } from '@/services/api'
import { useCalendarStore } from './calendar'

export const useEventsStore = defineStore('events', () => {
  const events = ref<CalendarEvent[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const eventsByDate = computed(() => {
    const map: Record<string, CalendarEvent[]> = {}
    const eventList = events.value || []
    for (const event of eventList) {
      const key = event.start.split('T')[0].split(' ')[0]
      if (!map[key]) map[key] = []
      map[key].push(event)
    }
    return map
  })

  async function fetchEvents() {
    loading.value = true
    error.value = null
    try {
      const calendarStore = useCalendarStore()
      const start = format(calendarStore.monthStart, 'yyyy-MM-dd')
      const end = format(calendarStore.monthEnd, 'yyyy-MM-dd')
      events.value = await eventApi.getEvents(start, end)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to fetch events'
    } finally {
      loading.value = false
    }
  }

  async function createEvent(input: EventInput): Promise<CalendarEvent> {
    try {
      const newEvent = await eventApi.createEvent(input)
      events.value.push(newEvent)
      return newEvent
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to create event'
      throw e
    }
  }

  async function updateEvent(id: number, input: EventInput): Promise<CalendarEvent> {
    try {
      const updated = await eventApi.updateEvent(id, input)
      const index = events.value.findIndex(e => e.id === id)
      if (index !== -1) {
        events.value[index] = updated
      }
      return updated
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to update event'
      throw e
    }
  }

  async function deleteEvent(id: number): Promise<void> {
    try {
      await eventApi.deleteEvent(id)
      events.value = events.value.filter(e => e.id !== id)
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Failed to delete event'
      throw e
    }
  }

  return {
    events,
    loading,
    error,
    eventsByDate,
    fetchEvents,
    createEvent,
    updateEvent,
    deleteEvent
  }
})
