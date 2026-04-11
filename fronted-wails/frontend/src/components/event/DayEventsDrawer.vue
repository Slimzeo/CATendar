<script setup lang="ts">
import { computed, ref } from 'vue'
import { NDrawer, NDrawerContent, NButton, NEmpty } from 'naive-ui'
import { format } from 'date-fns'
import { useThemeStore } from '@/stores/theme'
import type { CalendarEvent } from '@/types'
import { openExternal } from '@/wails'

const props = defineProps<{
  show: boolean
  date: string
  events: CalendarEvent[]
  drawerSide: 'left' | 'right'
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'create', date: string): void
  (e: 'edit', event: CalendarEvent): void
  (e: 'toggle-side'): void
}>()

const themeStore = useThemeStore()

// Track which event items are expanded
const expandedId = ref<number | null>(null)

const formattedDate = computed(() => {
  if (!props.date) return ''
  return format(new Date(props.date), 'EEEE, MMMM d, yyyy')
})

function handleClose() {
  emit('update:show', false)
}

function handleAdd() {
  emit('create', props.date)
}

function handleEventClick(event: CalendarEvent) {
  // Toggle expand/collapse; if already expanded, collapse
  if (expandedId.value === event.id) {
    expandedId.value = null
  } else {
    expandedId.value = event.id
  }
}

function handleEditClick(event: CalendarEvent, e: MouseEvent) {
  e.stopPropagation()
  emit('edit', event)
}

function hexToRgb(hex: string): [number, number, number] {
  const n = parseInt(hex.replace('#', ''), 16)
  return [(n >> 16) & 255, (n >> 8) & 255, n & 255]
}

function rgbToLinear(c: number): number {
  const v = c / 255
  return v <= 0.03928 ? v / 12.92 : Math.pow((v + 0.055) / 1.055, 2.4)
}

function relativeLuminance(r: number, g: number, b: number): number {
  return 0.2126 * rgbToLinear(r) + 0.7152 * rgbToLinear(g) + 0.0722 * rgbToLinear(b)
}

function wcagContrast(l1: number, l2: number): number {
  const lighter = Math.max(l1, l2)
  const darker = Math.min(l1, l2)
  return (lighter + 0.05) / (darker + 0.05)
}

function blendColor(eventHex: string, bgHex: string, alpha: number): [number, number, number] {
  const [er, eg, eb] = hexToRgb(eventHex)
  const [br, bg, bb] = hexToRgb(bgHex)
  return [
    Math.round(er * alpha + br * (1 - alpha)),
    Math.round(eg * alpha + bg * (1 - alpha)),
    Math.round(eb * alpha + bb * (1 - alpha))
  ]
}

// DayEventsDrawer uses 20% opacity background (event.color + '20')
function adjustColor(hex: string): string {
  const bg = themeStore.isDark ? '#1a1a1a' : '#ffffff'
  const [r, g, b] = blendColor(hex, bg, 0.2)
  const bgLuminance = relativeLuminance(r, g, b)

  const blackLuminance = relativeLuminance(26, 26, 26)
  const whiteLuminance = relativeLuminance(255, 255, 255)

  const contrastBlack = wcagContrast(bgLuminance, blackLuminance)
  const contrastWhite = wcagContrast(bgLuminance, whiteLuminance)

  return contrastWhite >= contrastBlack ? '#ffffff' : '#1a1a1a'
}

function formatTime(dateStr: string): string {
  const d = new Date(dateStr)
  const h = d.getHours().toString().padStart(2, '0')
  const m = d.getMinutes().toString().padStart(2, '0')
  return `${h}:${m}`
}

// Split description into text/URL segments for safe Vue rendering (no v-html)
interface DescSegment {
  type: 'text' | 'url'
  value: string
}

function handleLinkClick(url: string) {
  openExternal(url)
}

function parseDescription(text: string): DescSegment[] {
  // Match protocol URLs first, then bare domains (www.xxx.com)
  const protocolRegex = /(https?:\/\/[^\s]+)/g
  const segments: DescSegment[] = []
  let lastIndex = 0
  let match

  // First pass: protocol URLs
  while ((match = protocolRegex.exec(text)) !== null) {
    if (match.index > lastIndex) {
      // Check if this gap contains a bare domain
      const gap = text.slice(lastIndex, match.index)
      let bareMatch
      const bareRegex = /\bwww\.[a-zA-Z0-9-]+\.[a-zA-Z]{2,}[^\s]*/g
      let foundBare = false
      while ((bareMatch = bareRegex.exec(gap)) !== null) {
        segments.push({ type: 'url', value: bareMatch[0].startsWith('http') ? bareMatch[0] : 'https://' + bareMatch[0] })
        foundBare = true
      }
      // If no bare domain found, the gap is plain text
      if (!foundBare && gap.trim()) {
        segments.push({ type: 'text', value: gap })
      }
    }
    segments.push({ type: 'url', value: match[1] })
    lastIndex = protocolRegex.lastIndex
  }

  // Remaining text: bare domains only
  const remaining = text.slice(lastIndex)
  if (remaining) {
    let bareMatch
    const bareRegex = /(\bwww\.[a-zA-Z0-9-]+\.[a-zA-Z]{2,}[^\s]*)/g
    let bareLastIndex = 0
    while ((bareMatch = bareRegex.exec(remaining)) !== null) {
      if (bareMatch.index > bareLastIndex) {
        segments.push({ type: 'text', value: remaining.slice(bareLastIndex, bareMatch.index) })
      }
      segments.push({ type: 'url', value: 'https://' + bareMatch[1] })
      bareLastIndex = bareRegex.lastIndex
    }
    if (bareLastIndex < remaining.length) {
      segments.push({ type: 'text', value: remaining.slice(bareLastIndex) })
    }
  }

  return segments
}
</script>

<template>
  <n-drawer
    :show="show"
    :placement="drawerSide"
    :width="340"
    :resizable="true"
    @update:show="handleClose"
  >
    <n-drawer-content closable @close="handleClose">
      <template #header>
        <div class="drawer-header">
          <span class="drawer-date">{{ formattedDate }}</span>
          <div class="drawer-actions">
            <n-button
              @click="emit('toggle-side')"
              title="Switch side"
              class="switch-side-btn"
              text
              size="small"
            >
              <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                <circle cx="5" cy="12" r="2" />
                <circle cx="12" cy="12" r="2" />
                <circle cx="19" cy="12" r="2" />
              </svg>
            </n-button>
          </div>
        </div>
      </template>

      <div class="drawer-body">
        <div v-if="events.length === 0" class="empty-state">
          <n-empty description="No events" size="small" />
        </div>

        <div v-else class="events-list">
          <div
            v-for="event in events"
            :key="event.id"
            class="event-item"
            :class="{ expanded: expandedId === event.id }"
            :style="{
              backgroundColor: event.color + '20',
              borderLeft: `3px solid ${event.color}`,
              color: adjustColor(event.color)
            }"
            @click="handleEventClick(event)"
          >
            <div class="event-header">
              <div class="event-time-row">
                <span class="event-time">{{ formatTime(event.start) }}</span>
                <span v-if="event.end && event.end !== event.start" class="event-time-range">
                  - {{ formatTime(event.end) }}
                </span>
              </div>
              <button
                class="edit-btn"
                :style="{ color: adjustColor(event.color) }"
                :title="expandedId === event.id ? 'Collapse' : 'Edit'"
                @click="(e) => expandedId === event.id ? (expandedId = null, e.stopPropagation()) : handleEditClick(event, e)"
              >
                <svg v-if="expandedId !== event.id" width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
                  <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
                </svg>
                <svg v-else width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
                  <line x1="18" y1="6" x2="6" y2="18" />
                  <line x1="6" y1="6" x2="18" y2="18" />
                </svg>
              </button>
            </div>
            <div class="event-title">{{ event.title }}</div>

            <!-- Expanded: description + edit shortcut -->
            <div v-if="expandedId === event.id && event.description" class="event-description">
              <template v-for="(seg, i) in parseDescription(event.description)" :key="i">
                <a
                  v-if="seg.type === 'url'"
                  :href="seg.value"
                  class="event-link"
                  @click.prevent="handleLinkClick(seg.value)"
                >{{ seg.value }}</a>
                <span v-else>{{ seg.value }}</span>
              </template>
            </div>
          </div>
        </div>

        <div class="drawer-footer">
          <n-button type="primary" block @click="handleAdd">
            Add Event
          </n-button>
        </div>
      </div>
    </n-drawer-content>
  </n-drawer>
</template>

<style scoped>
.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

.drawer-date {
  font-size: 14px;
  font-weight: 500;
}

.drawer-actions {
  display: flex;
  gap: 4px;
  align-items: center;
}

.switch-side-btn {
  color: var(--text-secondary);
  opacity: 0.5;
  padding: 2px 4px;
  transition: opacity 0.15s;
}

.switch-side-btn:hover {
  opacity: 0.9;
  color: var(--text-primary);
}

.drawer-body {
  display: flex;
  flex-direction: column;
  height: 100%;
  gap: 16px;
}

.empty-state {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.events-list {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
  overflow-y: auto;
}

.event-item {
  padding: 10px 12px;
  border-radius: 6px;
  cursor: pointer;
  transition: opacity 0.15s, transform 0.1s;
}

.event-item:hover {
  opacity: 0.85;
  transform: translateX(2px);
}

.event-item.expanded {
  transform: none;
}

.event-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2px;
}

.event-time-row {
  display: flex;
  align-items: center;
  gap: 4px;
}

.edit-btn {
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 2px 4px;
  border-radius: 4px;
  opacity: 0;
  transition: opacity 0.15s, background 0.1s;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.event-item:hover .edit-btn,
.event-item.expanded .edit-btn {
  opacity: 0.6;
}

.edit-btn:hover {
  opacity: 1 !important;
  background: rgba(128, 128, 128, 0.15);
}

.event-time {
  font-size: 12px;
  font-weight: 600;
  opacity: 0.8;
}

.event-time-range {
  font-size: 12px;
  font-weight: 500;
  opacity: 0.7;
}

.event-title {
  font-size: 13px;
  font-weight: 700;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.event-description {
  font-size: 11px;
  margin-top: 4px;
  opacity: 0.8;
  word-break: break-all;
  line-height: 1.5;
  max-height: none;
  overflow: visible;
  white-space: pre-wrap;
}

.event-link {
  color: inherit;
  text-decoration: underline;
  text-underline-offset: 2px;
  pointer-events: auto;
}

.event-link:hover {
  text-decoration-style: solid;
  opacity: 0.85;
}

.drawer-footer {
  padding-top: 12px;
  border-top: 1px solid var(--border-color);
}
</style>
