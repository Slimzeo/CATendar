<script setup lang="ts">
import { computed } from 'vue'
import { format } from 'date-fns'
import type { DayCell, CalendarEvent } from '@/types'
import { useEventsStore } from '@/stores/events'
import EventChip from './EventChip.vue'

const props = defineProps<{
  day: DayCell
}>()

const emit = defineEmits<{
  (e: 'showDay', date: string): void
  (e: 'create', date: string): void
  (e: 'edit', event: CalendarEvent): void
}>()

const eventsStore = useEventsStore()
const dateLabel = computed(() => format(props.day.date, 'd'))
const dateStr = computed(() => format(props.day.date, 'yyyy-MM-dd'))
const dayEvents = computed(() => eventsStore.eventsByDate[dateStr.value] || [])
const accessibleLabel = computed(() => {
  const date = format(props.day.date, 'EEEE, MMMM d, yyyy')
  const count = dayEvents.value.length
  return `${date}, ${count === 0 ? 'no events' : `${count} event${count === 1 ? '' : 's'}`}`
})

const maxVisibleEvents = 3
const visibleEvents = computed(() => dayEvents.value.slice(0, maxVisibleEvents))
const overflowCount = computed(() => Math.max(0, dayEvents.value.length - maxVisibleEvents))

function handleCellClick(e: MouseEvent) {
  if ((e.target as HTMLElement).closest('.event-chip')) return
  emit('showDay', dateStr.value)
}

function handleEventClick(event: CalendarEvent) {
  emit('edit', event)
}

function handleCellKeydown(event: KeyboardEvent) {
  if (event.key !== 'Enter' && event.key !== ' ') return
  event.preventDefault()
  emit('showDay', dateStr.value)
}
</script>

<template>
  <div
    class="calendar-cell"
    :class="{
      'other-month': !day.isCurrentMonth,
      'today': day.isToday,
      'weekend': day.isWeekend
    }"
    role="button"
    tabindex="0"
    :aria-label="accessibleLabel"
    @click="handleCellClick"
    @keydown="handleCellKeydown"
  >
    <div class="cell-header">
      <span v-if="day.isToday" class="cell-date today">{{ dateLabel }}</span>
      <span v-else class="cell-date">{{ dateLabel }}</span>
      <button
        type="button"
        class="cell-add"
        :aria-label="`Add event on ${format(day.date, 'MMMM d')}`"
        title="Add event"
        @click.stop="emit('create', dateStr)"
      >
        <svg viewBox="0 0 16 16" aria-hidden="true">
          <path d="M8 3v10M3 8h10" />
        </svg>
      </button>
    </div>
    <transition-group name="event-list" tag="div" class="events-container">
      <event-chip
        v-for="event in visibleEvents"
        :key="event.id"
        :event="event"
        @click.stop="handleEventClick(event)"
      />
      <div v-if="overflowCount > 0" key="event-overflow" class="more-events">
        +{{ overflowCount }} more
      </div>
    </transition-group>
  </div>
</template>

<style scoped>
.calendar-cell {
  background: var(--bg-elevated);
  padding: 6px;
  min-height: 88px;
  min-width: 0;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  overflow: hidden;
  transition: background-color 140ms ease, box-shadow 140ms ease;
}

.calendar-cell:hover {
  background: var(--bg-hover);
}

.calendar-cell:focus-visible {
  position: relative;
  z-index: 1;
  outline: 2px solid var(--focus-ring);
  outline-offset: -2px;
}

.calendar-cell.other-month {
  opacity: 1;
  background: var(--bg-other-month);
}

.calendar-cell.other-month .cell-date,
.calendar-cell.other-month .more-events {
  color: var(--text-other-month);
}

.calendar-cell.other-month :deep(.event-chip) {
  opacity: 0.82;
}

.calendar-cell.today {
  background: var(--bg-today);
  box-shadow: inset 0 2px 0 var(--primary-color);
}

.calendar-cell.weekend:not(.today):not(.other-month) {
  background: var(--bg-weekend);
}

.cell-header {
  display: flex;
  min-height: 26px;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 3px;
}

.cell-date {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
}

.cell-date.today {
  background: var(--primary-color);
  color: white;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  box-shadow: var(--shadow-sm);
}

.cell-add {
  display: grid;
  width: 24px;
  height: 24px;
  padding: 0;
  place-items: center;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--text-tertiary);
  cursor: pointer;
  opacity: 0;
  transition: opacity 140ms ease, background-color 140ms ease, color 140ms ease;
}

.calendar-cell:hover .cell-add,
.calendar-cell:focus-within .cell-add {
  opacity: 1;
}

.cell-add:hover {
  background: var(--primary-muted);
  color: var(--primary-color);
}

.cell-add:focus-visible {
  opacity: 1;
  outline: 2px solid var(--focus-ring);
}

.cell-add svg {
  width: 14px;
  height: 14px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-width: 1.8;
}

.events-container {
  position: relative;
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
  min-width: 0;
}

.event-list-enter-active,
.event-list-leave-active,
.event-list-move {
  transition: opacity 150ms ease, transform 190ms cubic-bezier(0.22, 1, 0.36, 1);
}

.event-list-enter-from,
.event-list-leave-to {
  opacity: 0;
  transform: translateY(5px) scale(0.98);
}

.more-events {
  font-size: 11px;
  color: var(--text-secondary);
  padding: 2px 7px;
  font-weight: 600;
}

@media (max-width: 700px), (max-height: 520px) {
  .calendar-cell {
    padding: 4px;
    min-height: 70px;
  }

  .cell-add {
    display: none;
  }
}
</style>
