<script setup lang="ts">
import { computed } from 'vue'
import { colorPalettes } from '@/constants/colorPalettes'

const model = defineModel<string>({ required: true })

const flatColors = computed(() => {
  const seen = new Set<string>()
  const result: { color: string; paletteId: string }[] = []
  for (const palette of colorPalettes) {
    for (const color of palette.colors) {
      if (!seen.has(color)) {
        seen.add(color)
        result.push({ color, paletteId: palette.id })
      }
    }
  }
  return result
})
</script>

<template>
  <div class="palette-picker">
    <div class="color-grid">
      <button
        v-for="item in flatColors"
        :key="item.color"
        class="color-swatch"
        :class="{ selected: model === item.color }"
        :style="{ backgroundColor: item.color }"
        :title="item.color"
        @click="model = item.color"
      />
    </div>
  </div>
</template>

<style scoped>
.palette-picker {
  padding: 8px 0;
}

.color-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.color-swatch {
  width: 28px;
  height: 28px;
  border-radius: 6px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: transform 0.15s, border-color 0.15s;
}

.color-swatch:hover {
  transform: scale(1.1);
}

.color-swatch.selected {
  border-color: #333;
  box-shadow: 0 0 0 2px white inset;
}
</style>
