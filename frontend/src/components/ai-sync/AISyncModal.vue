<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue'
import {
  NAlert,
  NButton,
  NCard,
  NInput,
  NInputNumber,
  NModal,
  NSelect,
  NSwitch,
  useMessage
} from 'naive-ui'
import { aiSyncApi } from '@/services/api'
import type { AIProvider, AISyncSetup, EmailAccountInput } from '@/types'

const props = defineProps<{ show: boolean }>()

const emit = defineEmits<{
  (event: 'update:show', value: boolean): void
  (event: 'saved'): void
}>()

const message = useMessage()
const loading = ref(false)
const saving = ref(false)
const testing = ref(false)
const loadError = ref('')
const setup = ref<AISyncSetup | null>(null)

const email = reactive<EmailAccountInput>({
  address: '',
  username: '',
  imapHost: '',
  imapPort: 993,
  folder: 'INBOX',
  useTLS: true,
  secret: ''
})

const agent = reactive({
  provider: 'codex' as AIProvider,
  codexPath: '',
  claudePath: ''
})

const providerOptions = [
  { label: 'Codex', value: 'codex' },
  { label: 'Claude Code', value: 'claude' }
]

const selectedAvailability = computed(() => {
  return setup.value?.agent.agents.find(item => item.provider === agent.provider)
})

const selectedPath = computed({
  get: () => agent.provider === 'codex' ? agent.codexPath : agent.claudePath,
  set: value => {
    if (agent.provider === 'codex') agent.codexPath = value
    else agent.claudePath = value
  }
})

watch(() => props.show, show => {
  if (show) void loadSetup()
})

async function loadSetup() {
  loading.value = true
  loadError.value = ''
  try {
    const result = await aiSyncApi.getSetup()
    setup.value = result
    Object.assign(email, {
      address: result.email.address,
      username: result.email.username,
      imapHost: result.email.imapHost,
      imapPort: result.email.imapPort || 993,
      folder: result.email.folder || 'INBOX',
      useTLS: result.email.useTLS,
      secret: ''
    })
    Object.assign(agent, result.agent.settings)
  } catch (error) {
    loadError.value = errorMessage(error, 'Could not load AI Sync settings')
  } finally {
    loading.value = false
  }
}

function applyEmailPreset() {
  const domain = email.address.trim().toLowerCase().split('@')[1]
  if (!domain) return
  const presets: Record<string, string> = {
    'qq.com': 'imap.qq.com',
    'gmail.com': 'imap.gmail.com',
    'outlook.com': 'outlook.office365.com',
    'hotmail.com': 'outlook.office365.com',
    'live.com': 'outlook.office365.com',
    'icloud.com': 'imap.mail.me.com'
  }
  if (!email.username) email.username = email.address.trim()
  if (!email.imapHost && presets[domain]) email.imapHost = presets[domain]
}

async function testConnection() {
  testing.value = true
  try {
    applyEmailPreset()
    await aiSyncApi.testEmail({ ...email })
    message.success('Email connection succeeded')
  } catch (error) {
    message.error(errorMessage(error, 'Email connection failed'))
  } finally {
    testing.value = false
  }
}

async function save() {
  saving.value = true
  try {
    applyEmailPreset()
    await aiSyncApi.saveEmail({ ...email })
    await aiSyncApi.saveSettings({
      provider: agent.provider,
      codexPath: agent.codexPath.trim(),
      claudePath: agent.claudePath.trim()
    })
    message.success('AI Sync settings saved')
    emit('saved')
    emit('update:show', false)
  } catch (error) {
    message.error(errorMessage(error, 'Could not save AI Sync settings'))
  } finally {
    saving.value = false
  }
}

function errorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error) return error.message
  if (typeof error === 'string' && error) return error
  return fallback
}
</script>

<template>
  <n-modal :show="show" :mask-closable="!saving && !testing" @update:show="value => emit('update:show', value)">
    <n-card class="ai-sync-card" :bordered="false" closable @close="emit('update:show', false)">
      <template #header>
        <div class="modal-heading">
          <span class="eyebrow">Local agent workflow</span>
          <h2>AI Sync</h2>
          <p>Read recent email through CATendar Core, then let your local agent plan the calendar.</p>
        </div>
      </template>

      <div v-if="loading" class="loading-state">Loading local configuration…</div>
      <n-alert v-else-if="loadError" type="error" :show-icon="false">{{ loadError }}</n-alert>

      <form v-else class="settings-form" @submit.prevent="save">
        <section class="settings-section">
          <div class="section-heading">
            <span class="step">01</span>
            <div>
              <h3>Email source</h3>
              <p>Standard IMAP. Message bodies stay out of CATendar's database.</p>
            </div>
          </div>

          <div class="field-grid two-columns">
            <label class="field">
              <span>Email address</span>
              <n-input v-model:value="email.address" placeholder="you@example.com" @blur="applyEmailPreset" />
            </label>
            <label class="field">
              <span>Username</span>
              <n-input v-model:value="email.username" placeholder="Defaults to email address" />
            </label>
          </div>

          <div class="field-grid server-grid">
            <label class="field host-field">
              <span>IMAP host</span>
              <n-input v-model:value="email.imapHost" placeholder="imap.example.com" />
            </label>
            <label class="field">
              <span>Port</span>
              <n-input-number v-model:value="email.imapPort" :min="1" :max="65535" />
            </label>
            <label class="field">
              <span>Folder</span>
              <n-input v-model:value="email.folder" placeholder="INBOX" />
            </label>
          </div>

          <div class="credential-row">
            <label class="field secret-field">
              <span>Password / app authorization code</span>
              <n-input
                v-model:value="email.secret"
                type="password"
                show-password-on="click"
                :placeholder="setup?.email.hasCredential ? 'Leave blank to keep saved credential' : 'Stored in the system keychain'"
              />
            </label>
            <label class="tls-control">
              <span>TLS</span>
              <n-switch v-model:value="email.useTLS" />
            </label>
          </div>

          <div class="section-action">
            <n-button :loading="testing" :disabled="saving" @click="testConnection">Test connection</n-button>
            <span v-if="setup?.email.configured" class="configured-mark">Credential saved in keychain</span>
          </div>
        </section>

        <section class="settings-section">
          <div class="section-heading">
            <span class="step">02</span>
            <div>
              <h3>Local agent</h3>
              <p>CATendar invokes the selected CLI automatically when you click AI Sync.</p>
            </div>
          </div>

          <div class="field-grid two-columns agent-grid">
            <label class="field">
              <span>Agent</span>
              <n-select v-model:value="agent.provider" :options="providerOptions" />
            </label>
            <label class="field">
              <span>Executable path <small>optional</small></span>
              <n-input v-model:value="selectedPath" placeholder="Auto detect" />
            </label>
          </div>

          <div class="agent-status" :class="{ ready: selectedAvailability?.installed }">
            <span class="status-dot" />
            <span v-if="selectedAvailability?.installed">
              {{ selectedAvailability.version || selectedAvailability.path }}
            </span>
            <span v-else>{{ selectedAvailability?.error || 'CLI not detected' }}</span>
          </div>
        </section>

        <div class="privacy-note">
          Email credentials remain in CATendar Core. Email text selected for planning is sent to the chosen agent's model service.
        </div>

        <footer class="form-actions">
          <n-button :disabled="saving || testing" @click="emit('update:show', false)">Cancel</n-button>
          <n-button type="primary" attr-type="submit" :loading="saving" :disabled="testing">Save settings</n-button>
        </footer>
      </form>
    </n-card>
  </n-modal>
</template>

<style scoped>
.ai-sync-card {
  width: min(680px, calc(100vw - 28px));
  max-height: calc(100vh - 24px);
}

.ai-sync-card :deep(.n-card__content) {
  overflow-y: auto;
}

.modal-heading h2,
.modal-heading p {
  margin: 0;
}

.modal-heading h2 {
  font-family: var(--font-display);
  font-size: 22px;
  letter-spacing: -0.035em;
}

.modal-heading p {
  max-width: 54ch;
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 12px;
  font-weight: 450;
}

.eyebrow {
  display: block;
  margin-bottom: 2px;
  color: var(--primary-color);
  font-size: 9px;
  font-weight: 750;
  letter-spacing: 0.12em;
  text-transform: uppercase;
}

.loading-state {
  padding: 48px 0;
  color: var(--text-tertiary);
  text-align: center;
}

.settings-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings-section {
  padding: 16px;
  border: 1px solid var(--border-soft);
  border-radius: 12px;
  background: color-mix(in srgb, var(--bg-secondary) 55%, transparent);
}

.section-heading {
  display: flex;
  align-items: flex-start;
  gap: 11px;
  margin-bottom: 15px;
}

.section-heading h3,
.section-heading p {
  margin: 0;
}

.section-heading h3 {
  font-size: 14px;
  font-weight: 700;
}

.section-heading p {
  margin-top: 2px;
  color: var(--text-tertiary);
  font-size: 11px;
}

.step {
  display: grid;
  width: 27px;
  height: 27px;
  flex: 0 0 auto;
  place-items: center;
  border-radius: 8px;
  background: var(--primary-muted);
  color: var(--primary-color);
  font-size: 9px;
  font-weight: 800;
}

.field-grid {
  display: grid;
  gap: 12px;
}

.two-columns {
  grid-template-columns: 1fr 1fr;
  margin-bottom: 12px;
}

.server-grid {
  grid-template-columns: minmax(0, 1.6fr) 100px minmax(110px, 0.7fr);
  margin-bottom: 12px;
}

.field {
  display: flex;
  min-width: 0;
  flex-direction: column;
  gap: 6px;
}

.field > span,
.tls-control > span {
  color: var(--text-secondary);
  font-size: 11px;
  font-weight: 650;
}

.field small {
  color: var(--text-tertiary);
  font-size: 9px;
  font-weight: 500;
}

.credential-row {
  display: flex;
  align-items: end;
  gap: 12px;
}

.secret-field {
  flex: 1;
}

.tls-control {
  display: flex;
  height: 34px;
  align-items: center;
  gap: 8px;
  padding: 0 9px;
  border: 1px solid var(--border-soft);
  border-radius: 8px;
  background: var(--bg-elevated);
}

.section-action {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 12px;
}

.configured-mark {
  color: var(--text-tertiary);
  font-size: 10px;
}

.agent-grid {
  margin-bottom: 10px;
}

.agent-status {
  display: flex;
  align-items: center;
  gap: 7px;
  color: var(--danger-color);
  font-size: 10px;
}

.agent-status.ready {
  color: #3d8b62;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: currentColor;
  box-shadow: 0 0 0 3px color-mix(in srgb, currentColor 13%, transparent);
}

.privacy-note {
  padding: 10px 12px;
  border-left: 2px solid var(--primary-color);
  background: var(--primary-muted);
  color: var(--text-secondary);
  font-size: 11px;
  line-height: 1.5;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding-top: 2px;
}

@media (max-width: 620px) {
  .two-columns,
  .server-grid {
    grid-template-columns: 1fr;
  }

  .credential-row {
    align-items: stretch;
    flex-direction: column;
  }

  .tls-control {
    width: fit-content;
  }
}
</style>
