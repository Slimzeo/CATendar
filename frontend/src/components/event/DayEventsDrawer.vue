<script setup lang="ts">
import { computed } from 'vue'
import { format, parseISO } from 'date-fns'
import { NButton, NDrawer, NDrawerContent, NEmpty, NTooltip } from 'naive-ui'
import type { CalendarEvent } from '@/types'
import { openExternal } from '@/wails'

interface DescriptionSegment {
  type: 'text' | 'url'
  value: string
  href?: string
}

const props = defineProps<{
  show: boolean
  date: string
  events: CalendarEvent[]
  drawerSide: 'left' | 'right'
}>()

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void
  (event: 'create', date: string): void
  (event: 'edit', calendarEvent: CalendarEvent): void
  (event: 'toggle-side'): void
  (event: 'open-email', sourceURL: string): void
}>()

const formattedDate = computed(() => {
  if (!props.date) return ''
  return format(parseISO(props.date), 'EEEE, MMMM d')
})

function formatTime(date: Date): string {
  return format(date, 'HH:mm')
}

function eventTime(event: CalendarEvent): string {
  if (event.allDay) return 'All day'

  const start = parseISO(event.start)
  if (!event.end) return formatTime(start)

  const end = parseISO(event.end)
  const startLabel = formatTime(start)
  const endLabel = formatTime(end)
  return startLabel === endLabel ? startLabel : `${startLabel}–${endLabel}`
}

function parseDescription(text: string): DescriptionSegment[] {
  const urlPattern = /(https?:\/\/[^\s]+|catendar:\/\/[^\s]+|www\.[^\s]+)/gi
  const segments: DescriptionSegment[] = []
  let previousIndex = 0

  for (const match of text.matchAll(urlPattern)) {
    const index = match.index ?? 0
    if (index > previousIndex) {
      segments.push({ type: 'text', value: text.slice(previousIndex, index) })
    }
    const value = match[0]
    segments.push({
      type: 'url',
      value,
      href: value.startsWith('www.') ? `https://${value}` : value
    })
    previousIndex = index + value.length
  }

  if (previousIndex < text.length) {
    segments.push({ type: 'text', value: text.slice(previousIndex) })
  }
  return segments
}

function editEvent(event: CalendarEvent) {
  emit('edit', event)
}

function handleEventKeydown(event: KeyboardEvent, calendarEvent: CalendarEvent) {
  if (event.key !== 'Enter' && event.key !== ' ') return
  event.preventDefault()
  editEvent(calendarEvent)
}

function openLink(url: string) {
  if (url.startsWith('catendar://email/')) {
    emit('open-email', url)
    return
  }
  openExternal(url)
}
</script>

<template>
  <n-drawer
    :show="show"
    :placement="drawerSide"
    :width="360"
    :resizable="true"
    display-directive="show"
    @update:show="value => emit('update:show', value)"
  >
    <n-drawer-content closable @close="emit('update:show', false)">
      <template #header>
        <div class="drawer-header">
          <div>
            <span class="drawer-kicker">Day plan</span>
            <h2>{{ formattedDate }}</h2>
          </div>
          <n-tooltip trigger="hover">
            <template #trigger>
              <n-button
                quaternary
                circle
                size="small"
                :aria-label="drawerSide === 'left' ? 'Move panel to the right' : 'Move panel to the left'"
                @click="emit('toggle-side')"
              >
                <svg viewBox="0 0 20 20" aria-hidden="true">
                  <path v-if="drawerSide === 'left'" d="M4 4h12v12H4zM11 7l3 3-3 3M6 10h8" />
                  <path v-else d="M4 4h12v12H4zM9 7l-3 3 3 3M6 10h8" />
                </svg>
              </n-button>
            </template>
            {{ drawerSide === 'left' ? 'Move panel right' : 'Move panel left' }}
          </n-tooltip>
        </div>
      </template>

      <div class="drawer-body">
        <div v-if="events.length === 0" class="empty-state">
          <n-empty description="No events yet" size="small" />
          <p>Select “Add event” to plan this day.</p>
        </div>

        <div v-else class="events-list">
          <article
            v-for="event in events"
            :key="event.id"
            class="event-row"
            :style="{ '--event-color': event.color }"
            role="button"
            tabindex="0"
            :aria-label="`Edit ${event.title}, ${eventTime(event)}`"
            @click="editEvent(event)"
            @keydown="keyboardEvent => handleEventKeydown(keyboardEvent, event)"
          >
            <span class="event-marker" aria-hidden="true" />
            <div class="event-content">
              <div class="event-meta">
                <span>{{ eventTime(event) }}</span>
                <span v-if="event.bold" class="emphasis-label">Emphasized</span>
              </div>
              <h3 :class="{ bold: event.bold }">{{ event.title }}</h3>
              <p v-if="event.description" class="event-description">
                <template v-for="(segment, index) in parseDescription(event.description)" :key="index">
                  <a
                    v-if="segment.type === 'url'"
                    :href="segment.href"
                    @click.stop.prevent="openLink(segment.href!)"
                  >{{ segment.value }}</a>
                  <span v-else>{{ segment.value }}</span>
                </template>
              </p>
            </div>
            <span class="edit-hint" aria-hidden="true">
              <svg viewBox="0 0 20 20"><path d="m12.8 4.2 3 3L7.5 15.5 4 16l.5-3.5 8.3-8.3Z" /></svg>
            </span>
          </article>
        </div>

        <footer class="drawer-footer">
          <n-button type="primary" block size="large" @click="emit('create', date)">
            <template #icon>
              <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M10 4v12M4 10h12" /></svg>
            </template>
            Add event
          </n-button>
        </footer>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped>
.drawer-header {
  display: flex;
  width: 100%;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.drawer-kicker {
  display: block;
  margin-bottom: 2px;
  color: var(--primary-color);
  font-size: 10px;
  font-weight: 750;
  letter-spacing: 0.11em;
  text-transform: uppercase;
}

.drawer-header h2 {
  margin: 0;
  color: var(--text-primary);
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 650;
  letter-spacing: -0.02em;
}

.drawer-header svg,
.drawer-footer svg,
.edit-hint svg {
  width: 17px;
  height: 17px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.7;
}

.drawer-body {
  display: flex;
  height: 100%;
  min-height: 0;
  flex-direction: column;
}

.empty-state {
  display: grid;
  flex: 1;
  align-content: center;
  justify-items: center;
  gap: 9px;
  color: var(--text-tertiary);
  text-align: center;
}

.empty-state p {
  max-width: 24ch;
  margin: 0;
  font-size: 12px;
  line-height: 1.5;
}

.events-list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.event-row {
  display: grid;
  grid-template-columns: 8px minmax(0, 1fr) 24px;
  gap: 10px;
  align-items: start;
  padding: 14px 6px;
  border-bottom: 1px solid var(--border-soft);
  border-radius: 7px;
  cursor: pointer;
  transition: background-color 140ms ease;
}

:global(.n-drawer.slide-in-from-left-transition-enter-active),
:global(.n-drawer.slide-in-from-right-transition-enter-active) {
  transition-duration: 170ms !important;
  transition-timing-function: cubic-bezier(0.22, 1, 0.36, 1) !important;
}

:global(.n-drawer.slide-in-from-left-transition-leave-active),
:global(.n-drawer.slide-in-from-right-transition-leave-active) {
  transition-duration: 130ms !important;
}

.event-row:hover {
  background: var(--bg-hover);
}

.event-row:focus-visible {
  outline: 2px solid var(--focus-ring);
  outline-offset: -2px;
}

.event-marker {
  width: 7px;
  height: 7px;
  margin-top: 6px;
  border-radius: 50%;
  background: var(--event-color);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--event-color) 14%, transparent);
}

.event-content {
  min-width: 0;
}

.event-meta {
  display: flex;
  align-items: center;
  gap: 8px;
  color: var(--text-secondary);
  font-size: 11px;
  font-variant-numeric: tabular-nums;
}

.emphasis-label {
  color: var(--primary-color);
  font-size: 9px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
}

.event-row h3 {
  margin: 3px 0 0;
  overflow: hidden;
  color: var(--text-primary);
  font-size: 14px;
  font-weight: 550;
  line-height: 1.4;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event-row h3.bold {
  font-weight: 750;
}

.event-description {
  display: -webkit-box;
  margin: 5px 0 0;
  overflow: hidden;
  color: var(--text-secondary);
  font-size: 12px;
  line-height: 1.55;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  text-wrap: pretty;
}

.event-description a {
  color: var(--primary-color);
  text-decoration: underline;
  text-decoration-color: color-mix(in srgb, var(--primary-color) 40%, transparent);
  text-underline-offset: 2px;
}

.edit-hint {
  display: grid;
  width: 24px;
  height: 24px;
  place-items: center;
  color: var(--text-tertiary);
  opacity: 0;
  transition: opacity 140ms ease, color 140ms ease;
}

.event-row:hover .edit-hint,
.event-row:focus-visible .edit-hint {
  color: var(--primary-color);
  opacity: 1;
}

.drawer-footer {
  flex: 0 0 auto;
  padding-top: 14px;
  border-top: 1px solid var(--border-soft);
  background: var(--bg-elevated);
}
</style>
