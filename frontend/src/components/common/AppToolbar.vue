<script setup lang="ts">
import { ref } from 'vue'
import { NButton, NTooltip } from 'naive-ui'
import { useCalendarStore } from '@/stores/calendar'
import { useThemeStore } from '@/stores/theme'

const calendarStore = useCalendarStore()
const themeStore = useThemeStore()

const isPinned = ref(false)

function togglePin() {
  const newState = !isPinned.value
  isPinned.value = newState
  if (window.electronAPI?.setAlwaysOnTop) {
    window.electronAPI.setAlwaysOnTop(newState, 'normal')
  }
}
</script>

<template>
  <div class="toolbar">
    <div class="toolbar-left">
      <n-button quaternary size="small" @click="calendarStore.prevMonth">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="15,18 9,12 15,6" />
        </svg>
      </n-button>
      <h2 class="month-label">{{ calendarStore.monthLabel }}</h2>
      <n-button quaternary size="small" @click="calendarStore.nextMonth">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="9,18 15,12 9,6" />
        </svg>
      </n-button>
      <n-button quaternary size="tiny" @click="calendarStore.goToToday">Today</n-button>
    </div>
    <div class="toolbar-right">
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button quaternary circle size="small" @click="themeStore.showTime = !themeStore.showTime">
            <svg
              v-if="themeStore.showTime"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
              <circle cx="12" cy="12" r="3" />
            </svg>
            <svg
              v-else
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
              <line x1="1" y1="1" x2="23" y2="23" />
            </svg>
          </n-button>
        </template>
        {{ themeStore.showTime ? 'Hide time' : 'Show time' }}
      </n-tooltip>
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button quaternary circle size="small" @click="togglePin">
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              :fill="isPinned ? 'currentColor' : 'none'"
              stroke="currentColor"
              stroke-width="2"
            >
              <path d="M12 2L12 12M12 12L8 8M12 12L16 8" />
              <circle cx="12" cy="17" r="4" />
            </svg>
          </n-button>
        </template>
        {{ isPinned ? 'Unpin window' : 'Pin window on top' }}
      </n-tooltip>
      <n-tooltip trigger="hover">
        <template #trigger>
          <n-button quaternary circle size="small" @click="themeStore.toggleTheme">
            <svg v-if="themeStore.isDark" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="5" />
              <line x1="12" y1="1" x2="12" y2="3" />
              <line x1="12" y1="21" x2="12" y2="23" />
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
              <line x1="1" y1="12" x2="3" y2="12" />
              <line x1="21" y1="12" x2="23" y2="12" />
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
            </svg>
            <svg v-else width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
            </svg>
          </n-button>
        </template>
        {{ themeStore.isDark ? 'Light mode' : 'Dark mode' }}
      </n-tooltip>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 12px;
  border-bottom: 1px solid var(--border-color);
  background: var(--bg-secondary);
  -webkit-app-region: drag;
}

.toolbar-left, .toolbar-right {
  display: flex;
  align-items: center;
  gap: 6px;
  -webkit-app-region: no-drag;
}

.month-label {
  font-size: 15px;
  font-weight: 500;
  text-align: center;
  margin: 0 4px;
  flex: 0 1 auto;
  min-width: 0;
}

@media (max-width: 900px) {
  .toolbar {
    padding: 8px 10px;
  }
  .toolbar-left, .toolbar-right {
    gap: 4px;
  }
  .month-label {
    font-size: 13px;
  }
}

@media (max-width: 640px) {
  .toolbar {
    padding: 6px 8px;
  }
  .month-label {
    font-size: 12px;
    font-weight: 400;
  }
}
</style>
