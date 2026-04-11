<script setup lang="ts">
import { ref, computed } from 'vue'
import { NModal, NCard, NInput, NDatePicker, NButton, NSpace } from 'naive-ui'
import type { CalendarEvent, EventInput } from '@/types'
import ColorPalettePicker from './ColorPalettePicker.vue'

const props = defineProps<{
  show: boolean
  event?: CalendarEvent | null
  defaultDate?: string
}>()

const emit = defineEmits<{
  (e: 'update:show', value: boolean): void
  (e: 'save', input: EventInput): void
  (e: 'delete', id: number): void
}>()

const title = ref('')
const description = ref('')
const startDate = ref<number | null>(null)
const startTime = ref<string>('09:00')
const endTime = ref<string>('10:00')
const allDay = ref(true)
const selectedColor = ref('#845ec2')
const bold = ref(false)

const isEditing = computed(() => !!props.event?.id)

const startTimestamp = computed(() => {
  if (!startDate.value) return null
  return startDate.value
})

function initFromEvent() {
  if (props.event) {
    title.value = props.event.title
    description.value = props.event.description || ''
    const startDt = new Date(props.event.start)
    startDate.value = startDt.getTime()
    startTime.value = formatTime(startDt)
    const endDt = props.event.end ? new Date(props.event.end) : null
    endTime.value = endDt ? formatTime(endDt) : ''
    allDay.value = props.event.allDay
    selectedColor.value = props.event.color
    bold.value = props.event.bold || false
  } else if (props.defaultDate) {
    title.value = ''
    description.value = ''
    startDate.value = new Date(props.defaultDate).getTime()
    startTime.value = '09:00'
    endTime.value = '10:00'
    allDay.value = true
    selectedColor.value = '#845ec2'
    bold.value = false
  }
}

function formatTime(date: Date): string {
  const h = date.getHours().toString().padStart(2, '0')
  const m = date.getMinutes().toString().padStart(2, '0')
  return `${h}:${m}`
}

function handleSave() {
  if (!title.value.trim() || !startTimestamp.value) return

  const dateStr = new Date(startTimestamp.value).toISOString().split('T')[0]
  const start = `${dateStr}T${startTime.value}:00`
  const end = endTime.value ? `${dateStr}T${endTime.value}:00` : start

  emit('save', {
    title: title.value.trim(),
    description: description.value.trim(),
    start,
    end,
    allDay: allDay.value,
    color: selectedColor.value,
    bold: bold.value
  })
}

function handleDelete() {
  if (props.event?.id) {
    emit('delete', props.event.id)
  }
}

function handleClose() {
  emit('update:show', false)
}

defineExpose({ initFromEvent })
</script>

<template>
  <n-modal :show="show" @update:show="handleClose" :mask-closable="true">
    <n-card
      style="width: 480px; max-width: 90vw;"
      :title="isEditing ? 'Edit Event' : 'New Event'"
      :bordered="false"
      closable
      @close="handleClose"
    >
      <div class="form-group">
        <div class="label-row">
          <label>Title</label>
          <button
            class="bold-btn"
            :class="{ active: bold }"
            title="Bold title"
            type="button"
            @click="bold = !bold"
          >B</button>
        </div>
        <n-input
          v-model:value="title"
          placeholder="Event title"
          size="large"
        />
      </div>

      <div class="form-group">
        <label>Date</label>
        <n-date-picker
          v-model:value="startDate"
          type="date"
          style="width: 100%"
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label>Start Time</label>
          <input
            v-model="startTime"
            type="time"
            class="time-input"
          />
        </div>
        <div class="form-group">
          <label>End Time (optional)</label>
          <input
            v-model="endTime"
            type="time"
            class="time-input"
            placeholder="optional"
          />
        </div>
      </div>

      <div class="form-group">
        <label>Color</label>
        <color-palette-picker v-model="selectedColor" />
      </div>

      <div class="form-group">
        <label>Description <span class="optional">(links supported)</span></label>
        <n-input
          v-model:value="description"
          type="textarea"
          placeholder="e.g. https://meet.google.com/abc-defg-hij"
          :rows="2"
        />
      </div>

      <div class="form-actions">
        <n-space>
          <n-button @click="handleClose">Cancel</n-button>
          <n-button type="primary" @click="handleSave" :disabled="!title.trim()">
            {{ isEditing ? 'Update' : 'Create' }}
          </n-button>
          <n-button v-if="isEditing" type="error" @click="handleDelete">
            Delete
          </n-button>
        </n-space>
      </div>
    </n-card>
  </n-modal>
</template>

<style scoped>
.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 6px;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
}

.form-group label .optional {
  font-weight: 400;
  opacity: 0.6;
}

.label-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.label-row label {
  margin-bottom: 0;
}

.bold-btn {
  width: 22px;
  height: 22px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: transparent;
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.15s, color 0.15s, border-color 0.15s;
  padding: 0;
}

.bold-btn:hover {
  border-color: var(--primary-color);
  color: var(--primary-color);
}

.bold-btn.active {
  background: var(--primary-color);
  border-color: var(--primary-color);
  color: #ffffff;
}

.form-actions {
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid var(--border-color);
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-row .form-group {
  flex: 1;
}

.time-input {
  width: 100%;
  padding: 6px 12px;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 14px;
  box-sizing: border-box;
}

.time-input:focus {
  outline: none;
  border-color: #845ec2;
}
</style>
