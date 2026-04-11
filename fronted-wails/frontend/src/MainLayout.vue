<script setup lang="ts">
import { onMounted, ref, watch, computed } from 'vue'
import { useMessage } from 'naive-ui'
import MonthlyCalendar from './components/calendar/MonthlyCalendar.vue'
import AppToolbar from './components/common/AppToolbar.vue'
import EventModal from './components/event/EventModal.vue'
import DayEventsDrawer from './components/event/DayEventsDrawer.vue'
import { useCalendarStore } from './stores/calendar'
import { useEventsStore } from './stores/events'
import type { CalendarEvent, EventInput } from './types'

const calendarStore = useCalendarStore()
const eventsStore = useEventsStore()
const message = useMessage()

// Event modal state
const showEventModal = ref(false)
const editingEvent = ref<CalendarEvent | null>(null)
const defaultDate = ref<string>('')
const eventModalRef = ref<InstanceType<typeof EventModal> | null>(null)

// Day drawer state
const showDayDrawer = ref(false)
const selectedDayDate = ref<string>('')
const drawerSide = ref<'left' | 'right'>('left')

// Events for selected day
const selectedDayEvents = computed(() => {
  if (!selectedDayDate.value) return []
  return eventsStore.eventsByDate[selectedDayDate.value] || []
})

function openDayDrawer(date: string) {
  selectedDayDate.value = date
  showDayDrawer.value = true
}

function openCreateModal(date: string) {
  editingEvent.value = null
  defaultDate.value = date
  showEventModal.value = true
  setTimeout(() => eventModalRef.value?.initFromEvent(), 0)
}

function openEditModal(event: CalendarEvent) {
  editingEvent.value = event
  defaultDate.value = ''
  showEventModal.value = true
  setTimeout(() => eventModalRef.value?.initFromEvent(), 0)
}

async function handleSave(input: EventInput) {
  try {
    if (editingEvent.value?.id) {
      await eventsStore.updateEvent(editingEvent.value.id, input)
      message.success('Event updated')
    } else {
      await eventsStore.createEvent(input)
      message.success('Event created')
    }
    showEventModal.value = false
  } catch {
    message.error('Failed to save event')
  }
}

async function handleDelete(id: number) {
  try {
    await eventsStore.deleteEvent(id)
    message.success('Event deleted')
    showEventModal.value = false
  } catch {
    message.error('Failed to delete event')
  }
}

watch(() => calendarStore.currentMonth, () => {
  eventsStore.fetchEvents()
})

onMounted(async () => {
  await eventsStore.fetchEvents()
})
</script>

<template>
  <div class="app-container">
    <app-toolbar />
    <monthly-calendar @showDay="openDayDrawer" @edit="openEditModal" />
    <day-events-drawer
      v-model:show="showDayDrawer"
      :date="selectedDayDate"
      :events="selectedDayEvents"
      :drawer-side="drawerSide"
      @create="openCreateModal"
      @edit="openEditModal"
      @toggle-side="drawerSide = drawerSide === 'left' ? 'right' : 'left'"
    />
    <event-modal
      ref="eventModalRef"
      v-model:show="showEventModal"
      :event="editingEvent"
      :default-date="defaultDate"
      @save="handleSave"
      @delete="handleDelete"
    />
  </div>
</template>

<style scoped>
.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-primary);
  color: var(--text-primary);
}
</style>
