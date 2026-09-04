<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { format, parseISO } from 'date-fns'
import { NButton, NCard, NDatePicker, NInput, NModal, NPopconfirm, NSwitch } from 'naive-ui'
import type { CalendarEvent, EventInput } from '@/types'
import ColorPalettePicker from './ColorPalettePicker.vue'

const props = withDefaults(defineProps<{
  show: boolean
  event?: CalendarEvent | null
  defaultDate?: string
  saving?: boolean
}>(), {
  event: null,
  defaultDate: '',
  saving: false
})

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void
  (event: 'save', input: EventInput): void
  (event: 'delete', id: number): void
}>()

const titleInputRef = ref<{ focus: () => void } | null>(null)
const title = ref('')
const description = ref('')
const startDate = ref<number | null>(null)
const startTime = ref('09:00')
const endTime = ref('10:00')
const allDay = ref(true)
const selectedColor = ref('#7663b5')
const bold = ref(false)
const formError = ref('')

const isEditing = computed(() => Boolean(props.event?.id))

watch(
  () => props.show,
  show => {
    if (!show) return
    initializeForm()
    nextTick(() => titleInputRef.value?.focus())
  },
  { immediate: true }
)

function initializeForm() {
  formError.value = ''

  if (props.event) {
    const start = parseISO(props.event.start)
    const end = props.event.end ? parseISO(props.event.end) : null
    title.value = props.event.title
    description.value = props.event.description || ''
    startDate.value = start.getTime()
    startTime.value = format(start, 'HH:mm')
    endTime.value = end ? format(end, 'HH:mm') : ''
    allDay.value = props.event.allDay
    selectedColor.value = props.event.color
    bold.value = props.event.bold
    return
  }

  title.value = ''
  description.value = ''
  startDate.value = props.defaultDate ? parseISO(props.defaultDate).getTime() : Date.now()
  startTime.value = '09:00'
  endTime.value = '10:00'
  allDay.value = true
  selectedColor.value = '#7663b5'
  bold.value = false
}

function handleSave() {
  const trimmedTitle = title.value.trim()
  if (!trimmedTitle) {
    formError.value = 'Add a title before saving.'
    titleInputRef.value?.focus()
    return
  }
  if (!startDate.value) {
    formError.value = 'Choose a date before saving.'
    return
  }
  if (!allDay.value && !startTime.value) {
    formError.value = 'Choose a start time.'
    return
  }
  if (!allDay.value && endTime.value && endTime.value <= startTime.value) {
    formError.value = 'End time must be later than start time.'
    return
  }

  formError.value = ''
  const date = format(new Date(startDate.value), 'yyyy-MM-dd')
  const start = allDay.value ? `${date}T00:00:00` : `${date}T${startTime.value}:00`
  const end = allDay.value
    ? start
    : endTime.value
      ? `${date}T${endTime.value}:00`
      : start

  emit('save', {
    title: trimmedTitle,
    description: description.value.trim(),
    start,
    end,
    allDay: allDay.value,
    color: selectedColor.value,
    bold: bold.value
  })
}

function handleDelete() {
  if (props.event?.id) emit('delete', props.event.id)
}
</script>

<template>
  <n-modal :show="show" :mask-closable="!saving" @update:show="value => emit('update:show', value)">
    <n-card
      class="event-card"
      :title="isEditing ? 'Edit event' : 'New event'"
      :bordered="false"
      closable
      :closable-disabled="saving"
      @close="emit('update:show', false)"
    >
      <form class="event-form" @submit.prevent="handleSave">
        <div class="form-group">
          <div class="label-row">
            <label>Title</label>
            <button
              type="button"
              class="bold-button"
              :class="{ active: bold }"
              :aria-pressed="bold"
              title="Emphasize title"
              @click="bold = !bold"
            >
              B
            </button>
          </div>
          <n-input
            ref="titleInputRef"
            v-model:value="title"
            placeholder="Event title"
            size="large"
            :disabled="saving"
            @input="formError = ''"
          />
        </div>

        <div class="form-row date-row">
          <div class="form-group date-field">
            <label>Date</label>
            <n-date-picker
              v-model:value="startDate"
              type="date"
              :disabled="saving"
              style="width: 100%"
            />
          </div>
          <label class="all-day-control">
            <span>
              <strong>All day</strong>
              <small>Hide start and end times</small>
            </span>
            <n-switch v-model:value="allDay" :disabled="saving" />
          </label>
        </div>

        <div v-if="!allDay" class="form-row time-row">
          <div class="form-group">
            <label for="event-start-time">Start time</label>
            <input id="event-start-time" v-model="startTime" type="time" class="time-input" :disabled="saving" />
          </div>
          <div class="form-group">
            <label for="event-end-time">End time <span class="optional">Optional</span></label>
            <input id="event-end-time" v-model="endTime" type="time" class="time-input" :disabled="saving" />
          </div>
        </div>

        <div class="form-group">
          <label>Color</label>
          <color-palette-picker v-model="selectedColor" />
        </div>

        <div class="form-group description-field">
          <label>Description <span class="optional">Links are clickable</span></label>
          <n-input
            v-model:value="description"
            type="textarea"
            placeholder="Notes, location, or meeting link"
            :rows="3"
            :disabled="saving"
          />
        </div>

        <p v-if="formError" class="form-error" role="alert">{{ formError }}</p>

        <footer class="form-actions">
          <n-popconfirm v-if="isEditing" @positive-click="handleDelete">
            <template #trigger>
              <n-button type="error" ghost :disabled="saving">Delete</n-button>
            </template>
            Delete this event? This cannot be undone.
          </n-popconfirm>
          <span v-else />

          <div class="primary-actions">
            <n-button :disabled="saving" @click="emit('update:show', false)">Cancel</n-button>
            <n-button type="primary" attr-type="submit" :loading="saving">
              {{ isEditing ? 'Save changes' : 'Create event' }}
            </n-button>
          </div>
        </footer>
      </form>
    </n-card>
  </n-modal>
</template>

<style scoped>
.event-card {
  width: min(520px, calc(100vw - 28px));
  max-height: calc(100vh - 24px);
}

.event-card :deep(.n-card__content) {
  overflow-y: auto;
}

.event-form {
  display: flex;
  flex-direction: column;
}

.form-group {
  min-width: 0;
  margin-bottom: 17px;
}

.form-group > label,
.label-row label {
  display: block;
  margin-bottom: 7px;
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 650;
}

.optional {
  margin-left: 4px;
  color: var(--text-tertiary);
  font-size: 10px;
  font-weight: 500;
}

.label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.label-row label {
  margin-bottom: 0;
}

.bold-button {
  display: grid;
  width: 28px;
  height: 28px;
  padding: 0;
  place-items: center;
  border: 1px solid var(--border-strong);
  border-radius: 7px;
  background: transparent;
  color: var(--text-secondary);
  font: inherit;
  font-size: 12px;
  font-weight: 800;
  cursor: pointer;
}

.bold-button:hover,
.bold-button.active {
  border-color: var(--primary-color);
  background: var(--primary-muted);
  color: var(--primary-color);
}

.bold-button:focus-visible {
  outline: 2px solid var(--focus-ring);
  outline-offset: 2px;
}

.form-row {
  display: flex;
  gap: 14px;
}

.form-row .form-group {
  flex: 1;
}

.date-row {
  align-items: center;
}

.date-field {
  flex: 1.2 !important;
}

.all-day-control {
  display: flex;
  flex: 1;
  min-width: 0;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin: 4px 0 17px;
  padding: 8px 10px;
  border: 1px solid var(--border-soft);
  border-radius: 9px;
  background: var(--bg-secondary);
  color: var(--text-primary);
  cursor: pointer;
}

.all-day-control span {
  display: flex;
  min-width: 0;
  flex-direction: column;
}

.all-day-control strong {
  font-size: 12px;
  font-weight: 650;
}

.all-day-control small {
  overflow: hidden;
  color: var(--text-tertiary);
  font-size: 10px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.time-input {
  width: 100%;
  height: 34px;
  padding: 0 10px;
  border: 1px solid var(--border-strong);
  border-radius: 8px;
  outline: none;
  background: var(--bg-elevated);
  color: var(--text-primary);
  color-scheme: light;
  font: inherit;
  font-size: 13px;
  font-variant-numeric: tabular-nums;
}

[data-theme='dark'] .time-input {
  color-scheme: dark;
}

.time-input:hover {
  border-color: var(--primary-color);
}

.time-input:focus-visible {
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px var(--focus-ring);
}

.description-field {
  margin-bottom: 10px;
}

.form-error {
  margin: 0 0 10px;
  color: var(--danger-color);
  font-size: 12px;
  line-height: 1.45;
}

.form-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 8px;
  padding-top: 16px;
  border-top: 1px solid var(--border-soft);
}

.primary-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 520px) {
  .form-row {
    flex-direction: column;
    gap: 0;
  }

  .all-day-control {
    margin-top: 0;
  }

  .form-actions {
    align-items: stretch;
    flex-direction: column-reverse;
  }

  .primary-actions {
    display: grid;
    grid-template-columns: 1fr 1fr;
  }
}
</style>
