<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import type { CalendarEvent } from '@/types'
import { useCalendarStore } from '@/stores/calendar'
import CalendarCell from './CalendarCell.vue'

const calendarStore = useCalendarStore()

const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

const emit = defineEmits<{
  (e: 'showDay', date: string): void
  (e: 'edit', event: CalendarEvent): void
}>()

// Zoom
const zoomLevel = ref(1)
const MAX_ZOOM = 1.5
const ZOOM_STEP = 0.1
const containerRef = ref<HTMLElement | null>(null)
const scalableRef = ref<HTMLElement | null>(null)

// Pan
const panOffset = ref({ x: 0, y: 0 })
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0 })

// Auto-fit zoom on mount/resize
function calculateAutoFitZoom() {
  if (!containerRef.value || !scalableRef.value) return 1
  const containerWidth = containerRef.value.clientWidth
  const containerHeight = containerRef.value.clientHeight
  const scalableWidth = scalableRef.value.scrollWidth
  const scalableHeight = scalableRef.value.scrollHeight

  const scaleX = containerWidth / scalableWidth
  const scaleY = containerHeight / scalableHeight
  return Math.min(scaleX, scaleY) // Auto-fit = fill container exactly
}

const MIN_ZOOM = 1 // Never shrink below 100%
const autoFitZoom = ref(1)

function handleWheel(e: WheelEvent) {
  if (e.ctrlKey) {
    e.preventDefault()
    const delta = e.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP
    const newZoom = Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, zoomLevel.value + delta))
    zoomLevel.value = Math.round(newZoom * 10) / 10
    if (scalableRef.value) {
      scalableRef.value.scrollTop += e.deltaY
    }
  } else {
    if (scalableRef.value) {
      e.preventDefault()
      scalableRef.value.scrollTop += e.deltaY
    }
  }
}

function handleMouseDown(e: MouseEvent) {
  if (e.button === 2) { // Right click
    isPanning.value = true
    panStart.value = { x: e.clientX - panOffset.value.x, y: e.clientY - panOffset.value.y }
    e.preventDefault()
  }
}

function handleMouseMove(e: MouseEvent) {
  if (isPanning.value) {
    panOffset.value = {
      x: e.clientX - panStart.value.x,
      y: e.clientY - panStart.value.y
    }
    clampPan()
  }
}

function clampPan() {
  if (!containerRef.value || !scalableRef.value) return
  const scaledWidth = scalableRef.value.scrollWidth * zoomLevel.value
  const scaledHeight = scalableRef.value.scrollHeight * zoomLevel.value
  const maxPanX = Math.max(0, scaledWidth - containerRef.value.clientWidth)
  const maxPanY = Math.max(0, scaledHeight - containerRef.value.clientHeight)
  panOffset.value = {
    x: Math.max(-maxPanX, Math.min(0, panOffset.value.x)),
    y: Math.max(-maxPanY, Math.min(0, panOffset.value.y))
  }
}

function handleMouseUp() {
  isPanning.value = false
}

function handleContextMenu(e: MouseEvent) {
  e.preventDefault()
}

function resetView() {
  zoomLevel.value = autoFitZoom.value
  panOffset.value = { x: 0, y: 0 }
  if (scalableRef.value) {
    scalableRef.value.scrollTop = 0
  }
}

let resizeObserver: ResizeObserver | null = null

onMounted(() => {
  autoFitZoom.value = calculateAutoFitZoom()
  zoomLevel.value = autoFitZoom.value
  resizeObserver = new ResizeObserver(() => {
    autoFitZoom.value = calculateAutoFitZoom()
    if (zoomLevel.value < autoFitZoom.value) {
      zoomLevel.value = autoFitZoom.value
    }
  })
  if (containerRef.value) {
    resizeObserver.observe(containerRef.value)
  }
})

onUnmounted(() => {
  resizeObserver?.disconnect()
})
</script>

<template>
  <div
    ref="containerRef"
    class="calendar"
    :class="{ panning: isPanning }"
    @wheel="handleWheel"
    @mousedown="handleMouseDown"
    @mousemove="handleMouseMove"
    @mouseup="handleMouseUp"
    @mouseleave="handleMouseUp"
    @contextmenu="handleContextMenu"
  >
    <div
      ref="scalableRef"
      class="calendar-scroll"
    >
      <div
        class="calendar-inner"
        :style="{
          transform: `translate(${panOffset.x}px, ${panOffset.y}px) scale(${zoomLevel})`
        }"
      >
        <div class="calendar-header">
          <div v-for="day in weekDays" :key="day" class="calendar-header-cell">
            {{ day }}
          </div>
        </div>
        <div class="calendar-grid">
          <template v-for="(week, weekIndex) in calendarStore.weeks" :key="weekIndex">
            <calendar-cell
              v-for="(day, dayIndex) in week"
              :key="dayIndex"
              :day="day"
              @showDay="(date) => emit('showDay', date)"
              @edit="(event) => emit('edit', event)"
            />
          </template>
        </div>
      </div>
    </div>
    <div v-if="zoomLevel !== autoFitZoom" class="zoom-indicator" @click="resetView" title="Reset view">
      {{ Math.round(zoomLevel * 100) }}%
    </div>
  </div>
</template>

<style scoped>
.calendar {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  position: relative;
  cursor: default;
  user-select: none;
}

.calendar.panning {
  cursor: grabbing;
}

.calendar-scroll {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  position: relative;
}

.calendar-inner {
  transform-origin: top left;
  transition: transform 0.1s ease;
  min-height: 100%;
}

.calendar-header {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  background: var(--border-color);
  flex-shrink: 0;
}

.calendar-header-cell {
  background: var(--bg-secondary);
  padding: 8px;
  text-align: center;
  font-weight: 600;
  font-size: 12px;
  color: var(--text-secondary);
  text-transform: uppercase;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1px;
  background: var(--border-color);
  min-width: 0;
}

.zoom-indicator {
  position: absolute;
  bottom: 12px;
  right: 12px;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  padding: 4px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  cursor: pointer;
  opacity: 0.85;
}

.zoom-indicator:hover {
  opacity: 1;
}
</style>
