<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'
import type { CalendarEvent } from '@/types'

const props = defineProps<{
  event: CalendarEvent
}>()

const themeStore = useThemeStore()

const chipStyle = computed(() => ({
  '--event-color': props.event.color
}))

function formatEventTime(event: CalendarEvent): string {
  if (event.allDay) return 'All day'

  const startDate = new Date(event.start)
  const startTimeStr = formatTime(startDate)

  if (event.end) {
    const endDate = new Date(event.end)
    if (startDate.toDateString() === endDate.toDateString()) {
      const endTimeStr = formatTime(endDate)
      if (endTimeStr !== startTimeStr) {
        return `${startTimeStr}–${endTimeStr}`
      }
    }
  }
  return startTimeStr
}

function formatTime(date: Date): string {
  const h = date.getHours().toString().padStart(2, '0')
  const m = date.getMinutes().toString().padStart(2, '0')
  return `${h}:${m}`
}
</script>

<template>
  <button class="event-chip" type="button" :style="chipStyle" :title="event.title">
    <span class="event-dot" aria-hidden="true" />
    <span class="event-title" :class="{ bold: event.bold }">{{ event.title }}</span>
    <span v-if="themeStore.showTime" class="event-time">{{ formatEventTime(event) }}</span>
  </button>
</template>

<style scoped>
.event-chip {
  display: flex;
  width: 100%;
  min-width: 0;
  align-items: center;
  gap: 6px;
  padding: 4px 7px;
  border: 0;
  border-radius: 6px;
  background: color-mix(in srgb, var(--event-color) 14%, var(--bg-elevated));
  color: var(--text-primary);
  font: inherit;
  font-size: 11px;
  line-height: 1.25;
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  cursor: pointer;
  transition: background-color 140ms ease, transform 140ms ease;
}

.event-dot {
  width: 6px;
  height: 6px;
  flex: 0 0 auto;
  border-radius: 50%;
  background: var(--event-color);
}

.event-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.event-title.bold {
  font-weight: 700;
}

.event-time {
  flex: 0 1 auto;
  color: var(--text-secondary);
  font-size: 10px;
  font-variant-numeric: tabular-nums;
  max-width: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-chip:hover {
  background: color-mix(in srgb, var(--event-color) 21%, var(--bg-elevated));
  transform: translateY(-1px);
}

.event-chip:focus-visible {
  outline: 2px solid var(--focus-ring);
  outline-offset: -1px;
}

@media (max-width: 700px) {
  .event-time {
    display: none;
  }
}
</style>
