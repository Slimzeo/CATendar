import { defineStore } from 'pinia'
import { nextTick, ref } from 'vue'

const THEME_KEY = 'catendar-theme'
const SHOW_TIME_KEY = 'catendar-show-time'

function prefersDarkTheme(): boolean {
  const savedTheme = localStorage.getItem(THEME_KEY)
  if (savedTheme) return savedTheme === 'dark'
  return window.matchMedia('(prefers-color-scheme: dark)').matches
}

export const useThemeStore = defineStore('theme', () => {
  const isDark = ref(prefersDarkTheme())
  const showTime = ref(localStorage.getItem(SHOW_TIME_KEY) !== 'false')

  function applyTheme() {
    document.documentElement.dataset.theme = isDark.value ? 'dark' : 'light'
  }

  function toggleTheme() {
    const updateTheme = async () => {
      isDark.value = !isDark.value
      localStorage.setItem(THEME_KEY, isDark.value ? 'dark' : 'light')
      applyTheme()
      await nextTick()
    }

    const reduceMotion = window.matchMedia('(prefers-reduced-motion: reduce)').matches
    const transitionDocument = document as Document & {
      startViewTransition?: (update: () => Promise<void>) => unknown
    }

    if (reduceMotion || !transitionDocument.startViewTransition) {
      void updateTheme()
      return
    }
    transitionDocument.startViewTransition(updateTheme)
  }

  function toggleShowTime() {
    showTime.value = !showTime.value
    localStorage.setItem(SHOW_TIME_KEY, String(showTime.value))
  }

  applyTheme()

  return {
    isDark,
    showTime,
    toggleTheme,
    toggleShowTime
  }
})
