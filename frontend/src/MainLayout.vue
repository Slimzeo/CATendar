<script setup lang="ts">
import { computed, defineAsyncComponent, onMounted, onUnmounted, ref, watch } from 'vue'
import { format } from 'date-fns'
import { NButton, useMessage } from 'naive-ui'
import MonthlyCalendar from './components/calendar/MonthlyCalendar.vue'
import AppToolbar from './components/common/AppToolbar.vue'
import DayEventsDrawer from './components/event/DayEventsDrawer.vue'
import { useCalendarStore } from './stores/calendar'
import { useEventsStore } from './stores/events'
import { aiSyncApi } from './services/api'
import type { AISyncStatus, CalendarEvent, EmailMessage, EventInput } from './types'

const loadEventModal = () => import('./components/event/EventModal.vue')
const loadAISyncModal = () => import('./components/ai-sync/AISyncModal.vue')
const loadEmailSourceModal = () => import('./components/ai-sync/EmailSourceModal.vue')
const EventModal = defineAsyncComponent(loadEventModal)
const AISyncModal = defineAsyncComponent(loadAISyncModal)
const EmailSourceModal = defineAsyncComponent(loadEmailSourceModal)

const calendarStore = useCalendarStore()
const eventsStore = useEventsStore()
const message = useMessage()

const showEventModal = ref(false)
const editingEvent = ref<CalendarEvent | null>(null)
const defaultDate = ref('')
const savingEvent = ref(false)

const showAISyncModal = ref(false)
const showAISyncStatus = ref(false)
const aiSyncStatus = ref<AISyncStatus>({
  state: 'idle',
  emailCount: 0,
  readEmails: 0,
  unreadableEmails: 0,
  created: 0,
  skipped: 0,
  rejected: 0,
  conflicts: 0
})
let aiSyncPoller: ReturnType<typeof setInterval> | null = null
let modalPreloadTimer: ReturnType<typeof setTimeout> | null = null
const isAISyncing = computed(() => aiSyncStatus.value.state === 'running')
const showEmailSource = ref(false)
const emailSourceLoading = ref(false)
const sourceEmail = ref<EmailMessage | null>(null)
const aiSyncCounts = computed(() => {
  const status = aiSyncStatus.value
  const parts = [`${status.created} created`, `${status.skipped} skipped`]
  if (status.rejected) parts.push(`${status.rejected} rejected`)
  if (status.conflicts) parts.push(`${status.conflicts} conflicts`)
  if (status.unreadableEmails) parts.push(`${status.unreadableEmails} unreadable`)
  return parts.join(' · ')
})

const showDayDrawer = ref(false)
const selectedDayDate = ref('')
const drawerSide = ref<'left' | 'right'>('left')

const selectedDayEvents = computed(() => {
  if (!selectedDayDate.value) return []
  return eventsStore.eventsByDate[selectedDayDate.value] || []
})

function openDayDrawer(date: string) {
  selectedDayDate.value = date
  showDayDrawer.value = true
}

function openCreateModal(date = format(new Date(), 'yyyy-MM-dd')) {
  showDayDrawer.value = false
  editingEvent.value = null
  defaultDate.value = date
  showEventModal.value = true
}

function openEditModal(event: CalendarEvent) {
  showDayDrawer.value = false
  editingEvent.value = event
  defaultDate.value = ''
  showEventModal.value = true
}

async function handleSave(input: EventInput) {
  if (savingEvent.value) return
  savingEvent.value = true
  try {
    if (editingEvent.value?.id) {
      await eventsStore.updateEvent(editingEvent.value.id, input)
      message.success('Changes saved')
    } else {
      await eventsStore.createEvent(input)
      message.success('Event created')
    }
    showEventModal.value = false
  } catch {
    message.error(eventsStore.error || 'Could not save the event')
  } finally {
    savingEvent.value = false
  }
}

async function handleDelete(id: number) {
  if (savingEvent.value) return
  savingEvent.value = true
  try {
    await eventsStore.deleteEvent(id)
    message.success('Event deleted')
    showEventModal.value = false
  } catch {
    message.error(eventsStore.error || 'Could not delete the event')
  } finally {
    savingEvent.value = false
  }
}

async function startAISync() {
  if (isAISyncing.value) return
  showAISyncStatus.value = true
  try {
    aiSyncStatus.value = await aiSyncApi.start()
    startAISyncPolling()
  } catch (error) {
    const detail = errorMessage(error, 'Could not start AI Sync')
    message.error(detail)
    showAISyncStatus.value = false
    if (/not configured|not available|not installed/i.test(detail)) {
      showAISyncModal.value = true
    }
  }
}

async function cancelAISync() {
  try {
    await aiSyncApi.cancel()
  } catch (error) {
    message.error(errorMessage(error, 'Could not cancel AI Sync'))
  }
}

async function openSourceEmail(sourceURL: string) {
  showEmailSource.value = true
  emailSourceLoading.value = true
  sourceEmail.value = null
  try {
    sourceEmail.value = await aiSyncApi.getSourceEmail(sourceURL)
  } catch (error) {
    showEmailSource.value = false
    message.error(errorMessage(error, 'Could not read source email'))
  } finally {
    emailSourceLoading.value = false
  }
}

function startAISyncPolling() {
  stopAISyncPolling()
  aiSyncPoller = setInterval(() => void pollAISyncStatus(), 650)
}

function stopAISyncPolling() {
  if (aiSyncPoller) clearInterval(aiSyncPoller)
  aiSyncPoller = null
}

async function pollAISyncStatus() {
  try {
    const previousState = aiSyncStatus.value.state
    const status = await aiSyncApi.status()
    aiSyncStatus.value = status
    if (status.state === 'running') {
      showAISyncStatus.value = true
      return
    }
    stopAISyncPolling()
    if (previousState !== 'running') return
    if (status.state === 'succeeded') {
      await eventsStore.fetchEvents()
      message.success(`AI Sync complete · ${aiSyncCounts.value}`)
    } else if (status.state === 'cancelled') {
      message.info('AI Sync cancelled')
    } else if (status.state === 'failed') {
      message.error(status.message || 'AI Sync failed')
    }
  } catch (error) {
    stopAISyncPolling()
    message.error(errorMessage(error, 'Could not read AI Sync status'))
  }
}

function errorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error) return error.message
  if (typeof error === 'string' && error) return error
  return fallback
}

function handleKeyboardShortcut(event: KeyboardEvent) {
  const target = event.target as HTMLElement | null
  if (target?.closest('input, textarea, select, button, [contenteditable="true"]')) return
  if (showEventModal.value || showDayDrawer.value) return

  if (event.key === 'ArrowLeft') {
    calendarStore.prevMonth()
  } else if (event.key === 'ArrowRight') {
    calendarStore.nextMonth()
  } else if (event.key.toLowerCase() === 't') {
    calendarStore.goToToday()
  } else if (event.key.toLowerCase() === 'n') {
    openCreateModal()
  } else {
    return
  }
  event.preventDefault()
}

watch(() => calendarStore.currentMonth, () => {
  eventsStore.fetchEvents()
})

onMounted(() => {
  window.addEventListener('keydown', handleKeyboardShortcut)
  eventsStore.fetchEvents()
  requestAnimationFrame(() => {
    void loadEventModal()
  })
  modalPreloadTimer = setTimeout(() => {
    void loadAISyncModal()
    void loadEmailSourceModal()
  }, 300)
  void aiSyncApi.status().then(status => {
    aiSyncStatus.value = status
    if (status.state === 'running') {
      showAISyncStatus.value = true
      startAISyncPolling()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyboardShortcut)
  stopAISyncPolling()
  if (modalPreloadTimer) clearTimeout(modalPreloadTimer)
})
</script>

<template>
  <div class="app-container">
    <app-toolbar
      :syncing="isAISyncing"
      @create="openCreateModal()"
      @ai-sync="startAISync"
      @ai-settings="showAISyncModal = true"
    />

    <transition name="sync-status">
      <div
        v-if="showAISyncStatus && aiSyncStatus.state !== 'idle'"
        class="ai-sync-status"
        :class="aiSyncStatus.state"
        role="status"
      >
        <span class="sync-indicator" aria-hidden="true" />
        <div class="sync-copy">
          <strong>{{ aiSyncStatus.message || 'AI Sync' }}</strong>
          <span v-if="aiSyncStatus.state === 'running' && aiSyncStatus.emailCount">
            {{ aiSyncStatus.emailCount }} recent emails
            <span v-if="aiSyncStatus.readEmails"> · {{ aiSyncStatus.readEmails }} bodies read</span>
            <span v-if="aiSyncStatus.unreadableEmails"> · {{ aiSyncStatus.unreadableEmails }} unreadable</span>
          </span>
          <span v-else-if="aiSyncStatus.state === 'succeeded'">{{ aiSyncCounts }}</span>
        </div>
        <n-button v-if="isAISyncing" text size="small" @click="cancelAISync">Cancel</n-button>
        <n-button v-else-if="aiSyncStatus.state === 'failed'" text size="small" @click="showAISyncModal = true">
          Settings
        </n-button>
        <n-button v-if="!isAISyncing" text size="small" aria-label="Dismiss" @click="showAISyncStatus = false">×</n-button>
      </div>
    </transition>

    <div class="calendar-area">
      <div v-if="eventsStore.loading" class="loading-line" role="progressbar" aria-label="Loading events" />

      <transition name="status">
        <div v-if="eventsStore.error" class="error-banner" role="alert">
          <span>{{ eventsStore.error }}</span>
          <n-button text size="small" @click="eventsStore.fetchEvents">Retry</n-button>
        </div>
      </transition>

      <monthly-calendar
        @show-day="openDayDrawer"
        @create="openCreateModal"
        @edit="openEditModal"
      />
    </div>

    <day-events-drawer
      v-model:show="showDayDrawer"
      :date="selectedDayDate"
      :events="selectedDayEvents"
      :drawer-side="drawerSide"
      @create="openCreateModal"
      @edit="openEditModal"
      @toggle-side="drawerSide = drawerSide === 'left' ? 'right' : 'left'"
      @open-email="openSourceEmail"
    />

    <event-modal
      v-if="showEventModal"
      v-model:show="showEventModal"
      :event="editingEvent"
      :default-date="defaultDate"
      :saving="savingEvent"
      @save="handleSave"
      @delete="handleDelete"
    />

    <AISyncModal
      v-model:show="showAISyncModal"
      @saved="showAISyncStatus = false"
    />

    <EmailSourceModal
      v-model:show="showEmailSource"
      :loading="emailSourceLoading"
      :message="sourceEmail"
    />
  </div>
</template>

<style scoped>
.app-container {
  display: flex;
  width: 100%;
  height: 100vh;
  min-width: 0;
  flex-direction: column;
  overflow: hidden;
  background: var(--bg-primary);
  color: var(--text-primary);
}

.calendar-area {
  position: relative;
  display: flex;
  flex: 1;
  min-height: 0;
  flex-direction: column;
}

.ai-sync-status {
  display: grid;
  z-index: 4;
  grid-template-columns: 9px minmax(0, 1fr) auto auto;
  min-height: 40px;
  align-items: center;
  gap: 10px;
  padding: 6px 13px;
  border-bottom: 1px solid var(--border-soft);
  background: color-mix(in srgb, var(--primary-muted) 66%, var(--bg-toolbar));
  color: var(--text-secondary);
}

.sync-indicator {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: var(--primary-color);
  box-shadow: 0 0 0 4px var(--primary-muted);
}

.ai-sync-status.running .sync-indicator {
  animation: sync-pulse 1.2s ease-in-out infinite;
}

.ai-sync-status.succeeded .sync-indicator {
  background: #4a9a6b;
}

.ai-sync-status.failed .sync-indicator,
.ai-sync-status.cancelled .sync-indicator {
  background: var(--danger-color);
}

.sync-copy {
  display: flex;
  min-width: 0;
  align-items: baseline;
  gap: 9px;
}

.sync-copy strong {
  overflow: hidden;
  color: var(--text-primary);
  font-size: 11px;
  font-weight: 650;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.sync-copy span {
  color: var(--text-tertiary);
  font-size: 10px;
  white-space: nowrap;
}

.sync-status-enter-active,
.sync-status-leave-active {
  transition: opacity 160ms ease, transform 180ms cubic-bezier(0.22, 1, 0.36, 1);
}

.sync-status-enter-from,
.sync-status-leave-to {
  opacity: 0;
  transform: translateY(-7px);
}

.loading-line {
  position: absolute;
  inset: 0 0 auto;
  z-index: 6;
  height: 2px;
  overflow: hidden;
  background: var(--primary-muted);
}

.loading-line::after {
  position: absolute;
  width: 32%;
  height: 100%;
  background: var(--primary-color);
  content: '';
  animation: loading 1s ease-in-out infinite;
}

.error-banner {
  position: absolute;
  top: 10px;
  left: 50%;
  z-index: 6;
  display: flex;
  max-width: min(520px, calc(100% - 24px));
  align-items: center;
  gap: 14px;
  padding: 8px 10px 8px 13px;
  border: 1px solid color-mix(in srgb, var(--danger-color) 32%, var(--border-soft));
  border-radius: 9px;
  background: var(--bg-elevated);
  box-shadow: var(--shadow-md);
  color: var(--danger-color);
  font-size: 12px;
  transform: translateX(-50%);
}

.error-banner span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-enter-active,
.status-leave-active {
  transition: opacity 160ms ease, transform 160ms ease;
}

.status-enter-from,
.status-leave-to {
  opacity: 0;
  transform: translate(-50%, -6px);
}

@keyframes loading {
  from { transform: translateX(-110%); }
  to { transform: translateX(330%); }
}

@keyframes sync-pulse {
  50% {
    box-shadow: 0 0 0 7px color-mix(in srgb, var(--primary-color) 3%, transparent);
    opacity: 0.72;
  }
}

@media (max-width: 560px) {
  .sync-copy span {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .loading-line::after {
    animation: none;
    transform: none;
  }
}
</style>
