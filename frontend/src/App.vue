<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider, darkTheme } from 'naive-ui'
import { useThemeStore } from './stores/theme'
import MainLayout from './MainLayout.vue'

const themeStore = useThemeStore()
const theme = computed(() => themeStore.isDark ? darkTheme : null)

const themeOverrides = computed(() => {
  const dark = themeStore.isDark
  const primary = dark ? '#b3a2df' : '#7663b5'
  const primaryHover = dark ? '#c3b5e7' : '#6955a8'
  const surface = dark ? '#26232d' : '#ffffff'
  const elevated = dark ? '#2e2a37' : '#ffffff'
  const text = dark ? '#f1eef5' : '#24212b'
  const secondaryText = dark ? '#b7b0c0' : '#686371'
  const mutedText = dark ? '#8f8799' : '#928c9b'
  const border = dark ? 'rgba(221, 211, 232, 0.17)' : 'rgba(55, 47, 68, 0.16)'

  return {
    common: {
      primaryColor: primary,
      primaryColorHover: primaryHover,
      primaryColorPressed: primaryHover,
      primaryColorSuppl: primary,
      bodyColor: dark ? '#19181e' : '#f7f6fa',
      cardColor: surface,
      modalColor: elevated,
      popoverColor: elevated,
      inputColor: surface,
      inputColorDisabled: dark ? '#302c39' : '#f0eef4',
      actionColor: dark ? '#302c39' : '#f0eef4',
      hoverColor: dark ? '#302c39' : '#f0edf6',
      borderColor: border,
      dividerColor: border,
      textColorBase: text,
      textColor1: text,
      textColor2: secondaryText,
      textColor3: mutedText,
      placeholderColor: mutedText,
      iconColor: secondaryText,
      borderRadius: '9px',
      fontFamily: 'var(--font-body)'
    },
    Button: {
      borderRadiusMedium: '8px',
      borderRadiusSmall: '7px',
      textColorPrimary: '#ffffff'
    },
    Input: {
      borderColor: border,
      borderColorHover: primary,
      borderColorFocus: primary,
      boxShadowFocus: `0 0 0 3px ${dark ? 'rgba(179, 162, 223, 0.28)' : 'rgba(118, 99, 181, 0.22)'}`
    },
    Card: {
      color: surface,
      colorModal: elevated,
      borderRadius: '14px'
    },
    Drawer: {
      color: surface,
      headerBorderBottom: border,
      bodyBorderBottom: border
    },
    Tooltip: {
      color: elevated,
      textColor: text,
      borderRadius: '7px'
    },
    Empty: {
      textColor: mutedText
    }
  }
})
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-message-provider placement="bottom-right" :duration="2400">
      <main-layout />
    </n-message-provider>
  </n-config-provider>
</template>
