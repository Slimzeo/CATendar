import { defineStore } from 'pinia'
import { ref } from 'vue'
import { darkTheme, lightTheme, type GlobalTheme } from 'naive-ui'

export const THEME_TRANSITION = '180ms cubic-bezier(0.2, 0, 0, 1)'

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(false)
  const theme = ref<GlobalTheme>(lightTheme)
  const showTime = ref(true)

  function toggleTheme() {
    isDark.value = !isDark.value
    theme.value = isDark.value ? darkTheme : lightTheme
    document.documentElement.setAttribute('data-theme', isDark.value ? 'dark' : 'light')
  }

  return {
    isDark,
    theme,
    showTime,
    toggleTheme
  }
})
