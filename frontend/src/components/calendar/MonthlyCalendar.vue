<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import type { CalendarEvent } from '@/types'
import { useCalendarStore } from '@/stores/calendar'
import CalendarCell from './CalendarCell.vue'

const calendarStore = useCalendarStore()

const weekDays = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
const MIN_ZOOM = 1
const MAX_ZOOM = 1.5
const ZOOM_STEP = 0.1

const emit = defineEmits<{
  (event: 'showDay', date: string): void
  (event: 'create', date: string): void
  (event: 'edit', calendarEvent: CalendarEvent): void
}>()

const scrollRef = ref<HTMLElement | null>(null)
const zoomLevel = ref(1)
const isPanning = ref(false)
const panStart = ref({ x: 0, y: 0, left: 0, top: 0 })

const canvasStyle = computed(() => ({
  width: `${zoomLevel.value * 100}%`,
  height: `${zoomLevel.value * 100}%`
}))

const innerStyle = computed<Record<string, string | number>>(() => ({
  width: `${100 / zoomLevel.value}%`,
  height: `${100 / zoomLevel.value}%`,
  transform: `scale(${zoomLevel.value})`,
  '--week-count': calendarStore.weeks.length
}))

const zoomLabel = computed(() => `${Math.round(zoomLevel.value * 100)}%`)
const monthDirection = ref<'forward' | 'backward'>('forward')
const monthTransitionName = computed(() => `month-${monthDirection.value}`)

watch(
  () => calendarStore.currentDate.getTime(),
  (nextDate, previousDate) => {
    monthDirection.value = nextDate >= previousDate ? 'forward' : 'backward'
  }
)

function setZoom(value: number) {
  const nextZoom = Math.round(Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, value)) * 10) / 10
  if (nextZoom === zoomLevel.value) return

  const scroller = scrollRef.value
  const previousZoom = zoomLevel.value
  const contentCenter = scroller
    ? {
        x: (scroller.scrollLeft + scroller.clientWidth / 2) / previousZoom,
        y: (scroller.scrollTop + scroller.clientHeight / 2) / previousZoom
      }
    : null

  zoomLevel.value = nextZoom

  if (scroller && contentCenter) {
    nextTick(() => {
      scroller.scrollLeft = contentCenter.x * nextZoom - scroller.clientWidth / 2
      scroller.scrollTop = contentCenter.y * nextZoom - scroller.clientHeight / 2
    })
  }
}

function handleWheel(event: WheelEvent) {
  if (!event.ctrlKey && !event.metaKey) return
  event.preventDefault()
  setZoom(zoomLevel.value + (event.deltaY > 0 ? -ZOOM_STEP : ZOOM_STEP))
}

function startPan(event: PointerEvent) {
  if (event.button !== 2 || !scrollRef.value) return
  isPanning.value = true
  panStart.value = {
    x: event.clientX,
    y: event.clientY,
    left: scrollRef.value.scrollLeft,
    top: scrollRef.value.scrollTop
  }
  scrollRef.value.setPointerCapture(event.pointerId)
  event.preventDefault()
}

function movePan(event: PointerEvent) {
  if (!isPanning.value || !scrollRef.value) return
  scrollRef.value.scrollLeft = panStart.value.left - (event.clientX - panStart.value.x)
  scrollRef.value.scrollTop = panStart.value.top - (event.clientY - panStart.value.y)
}

function stopPan(event: PointerEvent) {
  if (!isPanning.value || !scrollRef.value) return
  isPanning.value = false
  if (scrollRef.value.hasPointerCapture(event.pointerId)) {
    scrollRef.value.releasePointerCapture(event.pointerId)
  }
}

function resetZoom() {
  setZoom(1)
}
</script>

<template>
  <section class="calendar" aria-label="Monthly calendar" @wheel="handleWheel" @contextmenu.prevent>
    <div
      ref="scrollRef"
      class="calendar-scroll"
      :class="{ panning: isPanning }"
      @pointerdown="startPan"
      @pointermove="movePan"
      @pointerup="stopPan"
      @pointercancel="stopPan"
    >
      <div class="calendar-canvas" :style="canvasStyle">
        <div class="calendar-inner" :style="innerStyle">
          <div class="calendar-header" aria-hidden="true">
            <div
              v-for="(day, index) in weekDays"
              :key="day"
              class="calendar-header-cell"
              :class="{ weekend: index === 0 || index === 6 }"
            >
              {{ day }}
            </div>
          </div>

          <transition :name="monthTransitionName">
            <div :key="calendarStore.currentMonth" class="calendar-grid">
              <template v-for="week in calendarStore.weeks" :key="week[0].date.getTime()">
                <calendar-cell
                  v-for="day in week"
                  :key="day.date.getTime()"
                  :day="day"
                  @show-day="date => emit('showDay', date)"
                  @create="date => emit('create', date)"
                  @edit="event => emit('edit', event)"
                />
              </template>
            </div>
          </transition>
        </div>
      </div>
    </div>

    <div class="zoom-controls" aria-label="Calendar zoom controls">
      <button
        type="button"
        :disabled="zoomLevel <= MIN_ZOOM"
        aria-label="Zoom out"
        title="Zoom out"
        @click="setZoom(zoomLevel - ZOOM_STEP)"
      >
        <svg viewBox="0 0 16 16" aria-hidden="true"><path d="M3 8h10" /></svg>
      </button>
      <button type="button" class="zoom-value" title="Reset zoom" @click="resetZoom">
        {{ zoomLabel }}
      </button>
      <button
        type="button"
        :disabled="zoomLevel >= MAX_ZOOM"
        aria-label="Zoom in"
        title="Zoom in"
        @click="setZoom(zoomLevel + ZOOM_STEP)"
      >
        <svg viewBox="0 0 16 16" aria-hidden="true"><path d="M8 3v10M3 8h10" /></svg>
      </button>
    </div>
  </section>
</template>

<style scoped>
.calendar {
  position: relative;
  display: flex;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  background: var(--bg-primary);
  user-select: none;
}

.calendar-scroll {
  position: relative;
  flex: 1;
  min-width: 0;
  min-height: 0;
  overflow: auto;
  scrollbar-color: var(--border-strong) transparent;
  scrollbar-width: thin;
}

.calendar-scroll.panning {
  cursor: grabbing;
}

.calendar-canvas {
  position: relative;
  min-width: 100%;
  min-height: 100%;
}

.calendar-inner {
  position: absolute;
  inset: 0 auto auto 0;
  display: grid;
  grid-template-rows: 38px minmax(0, 1fr);
  transform-origin: top left;
}

.calendar-header,
.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, minmax(0, 1fr));
  gap: 1px;
  background: var(--border-soft);
}

.calendar-header {
  grid-row: 1;
  grid-column: 1;
}

.calendar-header-cell {
  display: grid;
  place-items: center;
  background: var(--bg-secondary);
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.09em;
  text-transform: uppercase;
}

.calendar-header-cell.weekend {
  color: var(--primary-color);
}

.calendar-grid {
  grid-row: 2;
  grid-column: 1;
  min-height: 0;
  grid-template-rows: repeat(var(--week-count), minmax(88px, 1fr));
}

.month-forward-enter-active,
.month-forward-leave-active,
.month-backward-enter-active,
.month-backward-leave-active {
  transition: opacity 160ms ease, transform 220ms cubic-bezier(0.22, 1, 0.36, 1);
  will-change: opacity, transform;
}

.month-forward-enter-from,
.month-backward-leave-to {
  opacity: 0;
  transform: translateX(22px);
}

.month-forward-leave-to,
.month-backward-enter-from {
  opacity: 0;
  transform: translateX(-22px);
}

.zoom-controls {
  position: absolute;
  right: 14px;
  bottom: 14px;
  z-index: 4;
  display: flex;
  align-items: center;
  padding: 3px;
  border: 1px solid var(--border-soft);
  border-radius: 9px;
  background: color-mix(in srgb, var(--bg-elevated) 90%, transparent);
  box-shadow: var(--shadow-md);
  backdrop-filter: blur(12px);
}

.zoom-controls button {
  display: grid;
  min-width: 30px;
  height: 30px;
  padding: 0 7px;
  place-items: center;
  border: 0;
  border-radius: 6px;
  background: transparent;
  color: var(--text-secondary);
  font: inherit;
  font-size: 11px;
  font-variant-numeric: tabular-nums;
  cursor: pointer;
}

.zoom-controls button:hover:not(:disabled) {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.zoom-controls button:focus-visible {
  outline: 2px solid var(--focus-ring);
}

.zoom-controls button:disabled {
  cursor: default;
  opacity: 0.35;
}

.zoom-controls .zoom-value {
  min-width: 48px;
  color: var(--text-primary);
  font-weight: 650;
}

.zoom-controls svg {
  width: 14px;
  height: 14px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-width: 1.8;
}

@media (max-width: 640px), (max-height: 460px) {
  .calendar-inner {
    grid-template-rows: 30px minmax(0, 1fr);
  }

  .calendar-header-cell {
    font-size: 9px;
  }

  .calendar-grid {
    grid-template-rows: repeat(var(--week-count), minmax(70px, 1fr));
  }

  .zoom-controls {
    right: 8px;
    bottom: 8px;
  }
}
</style>
