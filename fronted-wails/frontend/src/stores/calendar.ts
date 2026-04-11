import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  startOfMonth,
  endOfMonth,
  startOfWeek,
  endOfWeek,
  addDays,
  addMonths,
  subMonths,
  isSameMonth,
  isSameDay,
  format
} from 'date-fns'
import type { DayCell } from '@/types'
import { useEventsStore } from './events'

export const useCalendarStore = defineStore('calendar', () => {
  const currentDate = ref(new Date())

  const currentMonth = computed(() => format(currentDate.value, 'yyyy-MM-dd'))
  const monthLabel = computed(() => format(currentDate.value, 'MMMM yyyy'))

  const monthStart = computed(() => startOfMonth(currentDate.value))
  const monthEnd = computed(() => endOfMonth(currentDate.value))

  const calendarDays = computed((): DayCell[] => {
    const eventsStore = useEventsStore()
    const start = startOfWeek(monthStart.value)
    const end = endOfWeek(monthEnd.value)

    const days: DayCell[] = []
    let day = start

    while (day <= end) {
      const dayStr = format(day, 'yyyy-MM-dd')
      days.push({
        date: new Date(day),
        isCurrentMonth: isSameMonth(day, currentDate.value),
        isToday: isSameDay(day, new Date()),
        events: eventsStore.eventsByDate[dayStr] || []
      })
      day = addDays(day, 1)
    }

    return days
  })

  const weeks = computed(() => {
    const result: DayCell[][] = []
    const days = calendarDays.value
    for (let i = 0; i < days.length; i += 7) {
      result.push(days.slice(i, i + 7))
    }
    return result
  })

  function nextMonth() {
    currentDate.value = addMonths(currentDate.value, 1)
  }

  function prevMonth() {
    currentDate.value = subMonths(currentDate.value, 1)
  }

  function goToToday() {
    currentDate.value = new Date()
  }

  return {
    currentDate,
    currentMonth,
    monthLabel,
    monthStart,
    monthEnd,
    calendarDays,
    weeks,
    nextMonth,
    prevMonth,
    goToToday
  }
})
