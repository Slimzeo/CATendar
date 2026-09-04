<script setup lang="ts">
import { computed, ref } from 'vue'
import { isSameMonth } from 'date-fns'
import { NButton, NTooltip } from 'naive-ui'
import { useCalendarStore } from '@/stores/calendar'
import { useThemeStore } from '@/stores/theme'
import { setAlwaysOnTop } from '@/wails'

defineProps<{
  syncing: boolean
}>()

const emit = defineEmits<{
  (event: 'create'): void
  (event: 'ai-sync'): void
  (event: 'ai-settings'): void
}>()

const calendarStore = useCalendarStore()
const themeStore = useThemeStore()
const isPinned = ref(false)
const isCurrentMonth = computed(() => isSameMonth(calendarStore.currentDate, new Date()))

function togglePin() {
  const nextState = !isPinned.value
  if (setAlwaysOnTop(nextState)) {
    isPinned.value = nextState
  }
}
</script>

<template>
  <header class="toolbar">
    <button class="brand" type="button" title="Go to today" @click="calendarStore.goToToday">
      <img src="/catendar.png" alt="" />
      <span class="brand-name"><strong>CAT</strong>endar</span>
    </button>

    <nav class="month-nav" aria-label="Month navigation">
      <n-button quaternary circle size="small" aria-label="Previous month" @click="calendarStore.prevMonth">
        <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m12.5 4.5-5 5.5 5 5.5" /></svg>
      </n-button>
      <h1 class="month-label" aria-live="polite">{{ calendarStore.monthLabel }}</h1>
      <n-button quaternary circle size="small" aria-label="Next month" @click="calendarStore.nextMonth">
        <svg viewBox="0 0 20 20" aria-hidden="true"><path d="m7.5 4.5 5 5.5-5 5.5" /></svg>
      </n-button>
      <n-button
        quaternary
        size="small"
        class="today-button"
        :disabled="isCurrentMonth"
        @click="calendarStore.goToToday"
      >
        Today
      </n-button>
    </nav>

    <div class="toolbar-actions">
      <n-button
        secondary
        type="primary"
        size="small"
        class="ai-sync-button"
        :loading="syncing"
        :disabled="syncing"
        @click="emit('ai-sync')"
      >
        <template #icon>
          <svg viewBox="0 0 20 20" aria-hidden="true">
            <path d="m10 2 1.1 4.2L15 8l-3.9 1.8L10 14l-1.1-4.2L5 8l3.9-1.8L10 2ZM15.5 13l.6 2.1L18 16l-1.9.9-.6 2.1-.6-2.1L13 16l1.9-.9.6-2.1Z" />
          </svg>
        </template>
        <span class="ai-sync-label">{{ syncing ? 'Syncing' : 'AI Sync' }}</span>
      </n-button>

      <n-button type="primary" size="small" class="new-event-button" @click="emit('create')">
        <template #icon>
          <svg viewBox="0 0 20 20" aria-hidden="true"><path d="M10 4v12M4 10h12" /></svg>
        </template>
        <span class="new-event-label">New event</span>
      </n-button>

      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="small"
            :class="{ active: !themeStore.showTime }"
            :aria-label="themeStore.showTime ? 'Hide event times' : 'Show event times'"
            :aria-pressed="!themeStore.showTime"
            @click="themeStore.toggleShowTime"
          >
            <svg v-if="themeStore.showTime" viewBox="0 0 20 20" aria-hidden="true">
              <path d="M2.5 10s2.8-5 7.5-5 7.5 5 7.5 5-2.8 5-7.5 5-7.5-5-7.5-5Z" />
              <circle cx="10" cy="10" r="2.25" />
            </svg>
            <svg v-else viewBox="0 0 20 20" aria-hidden="true">
              <path d="M7 5.7A7.2 7.2 0 0 1 10 5c4.7 0 7.5 5 7.5 5a11.6 11.6 0 0 1-2.1 2.7M12.2 14.7A7.4 7.4 0 0 1 10 15c-4.7 0-7.5-5-7.5-5a11.9 11.9 0 0 1 2-2.6M3 3l14 14" />
            </svg>
          </n-button>
        </template>
        {{ themeStore.showTime ? 'Hide event times' : 'Show event times' }}
      </n-tooltip>

      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="small"
            :class="{ active: isPinned }"
            aria-label="Keep window on top"
            :aria-pressed="isPinned"
            @click="togglePin"
          >
            <svg viewBox="0 0 20 20" aria-hidden="true">
              <path d="M7 3.5h6l-.8 4 2.3 2.4v1.4h-9V9.9l2.3-2.4-.8-4ZM10 11.3v5.2" />
            </svg>
          </n-button>
        </template>
        {{ isPinned ? 'Stop keeping on top' : 'Keep window on top' }}
      </n-tooltip>

      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="small"
            :aria-label="themeStore.isDark ? 'Use light theme' : 'Use dark theme'"
            @click="themeStore.toggleTheme"
          >
            <svg v-if="themeStore.isDark" viewBox="0 0 20 20" aria-hidden="true">
              <circle cx="10" cy="10" r="3.5" />
              <path d="M10 2v2M10 16v2M2 10h2M16 10h2M4.3 4.3l1.4 1.4M14.3 14.3l1.4 1.4M15.7 4.3l-1.4 1.4M5.7 14.3l-1.4 1.4" />
            </svg>
            <svg v-else viewBox="0 0 20 20" aria-hidden="true">
              <path d="M16.8 12.2A7 7 0 0 1 7.8 3.2 7 7 0 1 0 16.8 12.2Z" />
            </svg>
          </n-button>
        </template>
        {{ themeStore.isDark ? 'Use light theme' : 'Use dark theme' }}
      </n-tooltip>

      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button
            quaternary
            circle
            size="small"
            :disabled="syncing"
            aria-label="AI Sync settings"
            @click="emit('ai-settings')"
          >
            <svg viewBox="0 0 20 20" aria-hidden="true">
              <circle cx="10" cy="10" r="2.4" />
              <path d="M16.4 11.2v-2.4l-2-.5a5 5 0 0 0-.5-1.1l1.1-1.8-1.7-1.7-1.8 1.1a5 5 0 0 0-1.1-.5L9.9 2H7.5L7 4.1a5 5 0 0 0-1.1.5L4.1 3.5 2.4 5.2 3.5 7a5 5 0 0 0-.5 1.1l-2 .5V11l2 .5a5 5 0 0 0 .5 1.1l-1.1 1.8 1.7 1.7L6 15a5 5 0 0 0 1.1.5l.5 2.1H10l.5-2.1a5 5 0 0 0 1.1-.5l1.8 1.1 1.7-1.7-1.1-1.8a5 5 0 0 0 .5-1.1l1.9-.3Z" />
            </svg>
          </n-button>
        </template>
        {{ syncing ? 'Cancel AI Sync before changing settings' : 'AI Sync settings' }}
      </n-tooltip>
    </div>
  </header>
</template>

<style scoped>
.toolbar {
  position: relative;
  z-index: 5;
  display: grid;
  grid-template-columns: minmax(128px, 1fr) auto minmax(360px, 1fr);
  min-height: 58px;
  align-items: center;
  gap: 16px;
  padding: 8px 14px;
  border-bottom: 1px solid var(--border-soft);
  background: var(--bg-toolbar);
  box-shadow: 0 1px 0 color-mix(in srgb, var(--bg-elevated) 65%, transparent);
  -webkit-app-region: drag;
}

.brand,
.month-nav,
.toolbar-actions {
  -webkit-app-region: no-drag;
}

.brand {
  display: inline-flex;
  width: fit-content;
  align-items: center;
  gap: 9px;
  padding: 3px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  color: var(--text-primary);
  cursor: pointer;
}

.brand:hover {
  color: var(--primary-color);
}

.brand:focus-visible {
  outline: 2px solid var(--focus-ring);
  outline-offset: 2px;
}

.brand img {
  width: 30px;
  height: 30px;
  border-radius: 8px;
  box-shadow: var(--shadow-sm);
}

.brand-name {
  font-family: var(--font-display);
  font-size: 14px;
  font-weight: 500;
  letter-spacing: -0.015em;
}

.brand-name strong {
  color: var(--primary-color);
  font-weight: 750;
}

.month-nav {
  display: flex;
  align-items: center;
  gap: 4px;
}

.month-label {
  min-width: 156px;
  margin: 0 3px;
  color: var(--text-primary);
  font-family: var(--font-display);
  font-size: 17px;
  font-weight: 650;
  letter-spacing: -0.025em;
  line-height: 1.2;
  text-align: center;
  text-wrap: balance;
}

.toolbar-actions {
  display: flex;
  min-width: 0;
  align-items: center;
  justify-content: flex-end;
  gap: 4px;
}

.toolbar :deep(.n-button) {
  min-width: 34px;
}

.toolbar :deep(.n-button.active) {
  background: var(--primary-muted);
  color: var(--primary-color);
}

.toolbar svg {
  width: 17px;
  height: 17px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 1.65;
}

.new-event-button {
  margin-right: 3px;
}

.ai-sync-button {
  margin-right: 2px;
}

@media (max-width: 860px) {
  .toolbar {
    grid-template-columns: auto 1fr auto;
    gap: 8px;
    padding-inline: 10px;
  }

  .brand-name,
  .new-event-label,
  .ai-sync-label {
    display: none;
  }

  .month-nav {
    justify-self: center;
  }

  .month-label {
    min-width: 132px;
    font-size: 15px;
  }
}

@media (max-width: 580px) {
  .toolbar {
    grid-template-columns: auto 1fr auto;
    min-height: 50px;
    gap: 4px;
    padding: 6px;
  }

  .brand img {
    width: 27px;
    height: 27px;
  }

  .today-button,
  .toolbar-actions :deep(.n-button:not(.new-event-button):first-of-type) {
    display: none;
  }

  .month-label {
    min-width: 104px;
    font-size: 13px;
  }
}
</style>
