<script setup lang="ts">
import { computed } from 'vue'
import { format, parseISO } from 'date-fns'
import { NButton, NCard, NModal, NSpin } from 'naive-ui'
import type { EmailMessage } from '@/types'

const props = defineProps<{
  show: boolean
  loading: boolean
  message: EmailMessage | null
}>()

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void
}>()

const receivedAt = computed(() => {
  if (!props.message?.receivedAt) return ''
  return format(parseISO(props.message.receivedAt), 'MMM d, yyyy · HH:mm')
})

const sentAt = computed(() => {
  if (!props.message?.sentAt) return ''
  return format(parseISO(props.message.sentAt), 'MMM d, yyyy · HH:mm')
})
</script>

<template>
  <n-modal :show="show" @update:show="value => emit('update:show', value)">
    <n-card class="source-card" :bordered="false" closable @close="emit('update:show', false)">
      <div v-if="loading" class="source-loading">
        <n-spin size="small" />
        <span>Reading source email…</span>
      </div>
      <article v-else-if="message" class="source-email">
        <header>
          <span class="eyebrow">Source email</span>
          <h2>{{ message.subject || 'Untitled email' }}</h2>
          <div class="email-meta">
            <span>{{ message.sender }}</span>
            <span class="email-times">
              <time>Received {{ receivedAt }}</time>
              <time v-if="sentAt">Sent {{ sentAt }}</time>
            </span>
          </div>
        </header>
        <pre>{{ message.text || 'This email has no readable text body.' }}</pre>
        <footer>
          <n-button @click="emit('update:show', false)">Close</n-button>
        </footer>
      </article>
    </n-card>
  </n-modal>
</template>

<style scoped>
.source-card {
  width: min(680px, calc(100vw - 28px));
  max-height: calc(100vh - 24px);
}

.source-card :deep(.n-card__content) {
  overflow: hidden;
}

.source-loading {
  display: flex;
  min-height: 180px;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: var(--text-tertiary);
  font-size: 12px;
}

.source-email {
  display: flex;
  max-height: calc(100vh - 90px);
  flex-direction: column;
}

.source-email header {
  padding-bottom: 14px;
  border-bottom: 1px solid var(--border-soft);
}

.eyebrow {
  color: var(--primary-color);
  font-size: 9px;
  font-weight: 750;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.source-email h2 {
  margin: 4px 0 5px;
  font-family: var(--font-display);
  font-size: 19px;
  letter-spacing: -0.025em;
}

.email-meta {
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  gap: 6px 14px;
  color: var(--text-tertiary);
  font-size: 11px;
}

.email-times {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 4px 12px;
}

.source-email pre {
  flex: 1;
  min-height: 160px;
  margin: 0;
  padding: 16px 2px;
  overflow: auto;
  color: var(--text-secondary);
  font-family: var(--font-body);
  font-size: 12px;
  line-height: 1.65;
  overflow-wrap: anywhere;
  white-space: pre-wrap;
}

.source-email footer {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  border-top: 1px solid var(--border-soft);
}
</style>
