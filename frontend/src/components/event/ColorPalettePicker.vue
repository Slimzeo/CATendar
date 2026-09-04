<script setup lang="ts">
import { computed } from 'vue'
import { eventColors } from '@/constants/eventColors'

const model = defineModel<string>({ required: true })

const availableColors = computed(() => {
  if (eventColors.some(color => color.value === model.value)) return eventColors
  return [{ name: 'Current color', value: model.value }, ...eventColors]
})
</script>

<template>
  <div class="color-grid" role="radiogroup" aria-label="Event color">
    <button
      v-for="item in availableColors"
      :key="item.value"
      type="button"
      class="color-swatch"
      :class="{ selected: model === item.value }"
      :style="{ '--swatch-color': item.value }"
      :aria-label="item.name"
      :aria-checked="model === item.value"
      role="radio"
      :title="item.name"
      @click="model = item.value"
    >
      <svg v-if="model === item.value" viewBox="0 0 16 16" aria-hidden="true">
        <path d="m3.2 8.1 3 3.1 6.6-6.6" />
      </svg>
    </button>
  </div>
</template>

<style scoped>
.color-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  padding: 2px;
}

.color-swatch {
  display: grid;
  width: 30px;
  height: 30px;
  padding: 0;
  place-items: center;
  border: 1px solid color-mix(in srgb, var(--swatch-color) 75%, black);
  border-radius: 50%;
  background: var(--swatch-color);
  color: white;
  cursor: pointer;
  box-shadow: 0 0 0 2px var(--bg-elevated);
  transition: transform 140ms ease, box-shadow 140ms ease;
}

.color-swatch:hover {
  transform: translateY(-1px);
}

.color-swatch.selected {
  box-shadow: 0 0 0 2px var(--bg-elevated), 0 0 0 4px var(--primary-color);
}

.color-swatch:focus-visible {
  outline: 3px solid var(--focus-ring);
  outline-offset: 3px;
}

.color-swatch svg {
  width: 16px;
  height: 16px;
  fill: none;
  stroke: currentColor;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-width: 2.2;
}
</style>
