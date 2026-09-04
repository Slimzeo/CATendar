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
  isWeekend,
  isSameMonth,
  isSameDay,
  format
} from 'date-fns'
import type { DayCell } from '@/types'

export const useCalendarStore = defineStore('calendar', () => {
  const currentDate = ref(new Date())

  const currentMonth = computed(() => format(currentDate.value, 'yyyy-MM'))
  const monthLabel = computed(() => format(currentDate.value, 'MMMM yyyy'))

  const monthStart = computed(() => startOfMonth(currentDate.value))
  const monthEnd = computed(() => endOfMonth(currentDate.value))
  const visibleStart = computed(() => startOfWeek(monthStart.value))
  const visibleEnd = computed(() => endOfWeek(monthEnd.value))

  const calendarDays = computed((): DayCell[] => {
    const days: DayCell[] = []
    let day = visibleStart.value

    while (day <= visibleEnd.value) {
      days.push({
        date: new Date(day),
        isCurrentMonth: isSameMonth(day, currentDate.value),
        isToday: isSameDay(day, new Date()),
        isWeekend: isWeekend(day)
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
    visibleStart,
    visibleEnd,
    weeks,
    nextMonth,
    prevMonth,
    goToToday
  }
})
