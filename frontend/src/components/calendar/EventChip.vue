<script setup lang="ts">
import { computed } from 'vue'
import { useThemeStore } from '@/stores/theme'
import type { CalendarEvent } from '@/types'

const props = defineProps<{
  event: CalendarEvent
}>()

const themeStore = useThemeStore()

const chipStyle = computed(() => ({
  backgroundColor: props.event.color + '30',
  borderLeft: `3px solid ${props.event.color}`,
  color: adjustColor(props.event.color)
}))

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

function adjustColor(hex: string): string {
  // Theme background colors
  const bgLight = '#ffffff'
  const bgDark = '#1a1a1a'
  const bg = themeStore.isDark ? bgDark : bgLight

  // Blend event color at 30% opacity with theme background
  const [r, g, b] = blendColor(hex, bg, 0.3)
  const bgLuminance = relativeLuminance(r, g, b)

  // Contrast against black (#1a1a1a) and white (#ffffff)
  const blackLuminance = relativeLuminance(26, 26, 26)
  const whiteLuminance = relativeLuminance(255, 255, 255)

  const contrastBlack = wcagContrast(bgLuminance, blackLuminance)
  const contrastWhite = wcagContrast(bgLuminance, whiteLuminance)

  // Pick the color with higher contrast
  return contrastWhite >= contrastBlack ? '#ffffff' : '#1a1a1a'
}

function formatEventTime(event: CalendarEvent): string {
  const startDate = new Date(event.start)
  const startTimeStr = formatTime(startDate)

  // If end exists and is different from start, show range
  if (event.end) {
    const endDate = new Date(event.end)
    // Only show end time if it's on the same day and different time
    if (startDate.toDateString() === endDate.toDateString()) {
      const endTimeStr = formatTime(endDate)
      if (endTimeStr !== startTimeStr) {
        return `${startTimeStr}-${endTimeStr}`
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
  <div class="event-chip" :style="chipStyle" :title="event.title">
    <span class="event-title" :style="{ fontWeight: event.bold ? '700' : '400' }">{{ event.title }}</span>
    <span v-if="themeStore.showTime" class="event-time">{{ formatEventTime(event) }}</span>
  </div>
</template>

<style scoped>
.event-chip {
  padding: 2px 6px;
  border-radius: 3px;
  font-size: 11px;
  white-space: nowrap;
  overflow: hidden;
  cursor: pointer;
  transition: transform 0.1s, opacity 0.15s;
  display: flex;
  gap: 4px;
  align-items: center;
}

.event-title {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.event-time {
  font-weight: 400;
  opacity: 0.75;
  flex-shrink: 1;
  font-size: 10px;
  max-width: 70px;
  overflow: hidden;
  text-overflow: ellipsis;
}

.event-chip:hover {
  transform: scale(1.02);
  opacity: 0.9;
}
</style>
