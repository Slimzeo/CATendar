<script setup lang="ts">
import { computed } from 'vue'
import { format } from 'date-fns'
import type { DayCell, CalendarEvent } from '@/types'
import EventChip from './EventChip.vue'

const props = defineProps<{
  day: DayCell
}>()

const emit = defineEmits<{
  (e: 'showDay', date: string): void
  (e: 'edit', event: CalendarEvent): void
}>()

const dateLabel = computed(() => format(props.day.date, 'd'))
const dateStr = computed(() => format(props.day.date, 'yyyy-MM-dd'))

const maxVisibleEvents = 3
const visibleEvents = computed(() => props.day.events.slice(0, maxVisibleEvents))
const overflowCount = computed(() => Math.max(0, props.day.events.length - maxVisibleEvents))

function handleCellClick(e: MouseEvent) {
  if ((e.target as HTMLElement).closest('.event-chip')) return
  emit('showDay', dateStr.value)
}

function handleEventClick(event: CalendarEvent) {
  emit('edit', event)
}
</script>

<template>
  <div
    class="calendar-cell"
    :class="{
      'other-month': !day.isCurrentMonth,
      'today': day.isToday
    }"
    @click="handleCellClick"
  >
    <div class="cell-header">
      <span v-if="day.isToday" class="cell-date today">{{ dateLabel }}</span>
      <span v-else class="cell-date">{{ dateLabel }}</span>
    </div>
    <div class="events-container">
      <event-chip
        v-for="event in visibleEvents"
        :key="event.id"
        :event="event"
        @click.stop="handleEventClick(event)"
      />
      <div v-if="overflowCount > 0" class="more-events">
        +{{ overflowCount }} more
      </div>
    </div>
  </div>
</template>

<style scoped>
.calendar-cell {
  background: var(--bg-elevated);
  padding: 4px;
  min-height: 100px;
  min-width: 0;
  display: flex;
  flex-direction: column;
  cursor: pointer;
  overflow: hidden;
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
}

.cell-header {
  margin-bottom: 2px;
}

.cell-date {
  font-size: 12px;
  font-weight: 500;
  color: var(--text-primary);
}

.cell-date.today {
  background: #845ec2;
  color: white;
  width: 22px;
  height: 22px;
  border-radius: 50%;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.events-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  overflow: hidden;
  min-width: 0;
}

.more-events {
  font-size: 10px;
  color: var(--text-secondary);
  padding: 2px 4px;
}
</style>
