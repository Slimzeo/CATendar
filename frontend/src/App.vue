<script setup lang="ts">
import { computed } from 'vue'
import { NConfigProvider, NMessageProvider, NDialogProvider, darkTheme, lightTheme } from 'naive-ui'
import { useThemeStore } from './stores/theme'
import MainLayout from './MainLayout.vue'

const themeStore = useThemeStore()
const theme = computed(() => themeStore.isDark ? darkTheme : lightTheme)

const themeOverrides = computed(() => {
  if (themeStore.isDark) {
    return {
      common: {
        primaryColor: '#9c84d6',
        primaryColorHover: '#b19fe6',
        primaryColorPressed: '#8a72c4',
        primaryColorSuppl: '#9c84d680',
        bodyColor: '#1f2430',
        cardColor: '#262c39',
        modalColor: '#2e3545',
        popoverColor: '#2e3545',
        tableColor: '#262c39',
        inputColor: '#2e3545',
        inputColorDisabled: '#323a4a',
        actionColor: '#323a4a',
        hoverColor: '#323a4a',
        borderColor: 'rgba(180, 195, 215, 0.14)',
        dividerColor: 'rgba(180, 195, 215, 0.10)',
        textColorBase: '#e8edf5',
        textColor1: '#e8edf5',
        textColor2: '#b6c0cf',
        textColor3: '#8e98a9',
        placeholderColor: '#8e98a9',
        iconColor: '#b6c0cf',
      },
      Button: {
        colorPrimary: '#9c84d6',
        colorHoverPrimary: '#b19fe6',
        colorPressedPrimary: '#8a72c4',
        textColorPrimary: '#ffffff',
        borderRadiusMedium: '6px',
      },
      Input: {
        color: '#2e3545',
        colorFocus: '#323a4a',
        borderColor: 'rgba(180, 195, 215, 0.14)',
        borderColorHover: 'rgba(180, 195, 215, 0.25)',
        borderColorFocus: '#9c84d6',
        textColor: '#e8edf5',
        placeholderColor: '#8e98a9',
      },
      Card: {
        color: '#262c39',
        colorModal: '#2e3545',
        borderColor: 'rgba(180, 195, 215, 0.10)',
      },
      Drawer: {
        color: '#262c39',
        headerBorderBottom: 'rgba(180, 195, 215, 0.10)',
        bodyBorderBottom: 'rgba(180, 195, 215, 0.10)',
      },
      Tooltip: {
        color: '#2e3545',
        textColor: '#e8edf5',
        borderRadius: '6px',
      },
      Message: {
        colorSuccess: '#262c39',
        colorError: '#262c39',
        colorWarning: '#262c39',
        colorInfo: '#262c39',
        textColorSuccess: '#e8edf5',
        textColorError: '#e8edf5',
        textColorWarning: '#e8edf5',
        textColorInfo: '#e8edf5',
      },
      Empty: {
        textColor: '#8e98a9',
      },
      DataTable: {
        'th-color': '#262c39',
        'td-color': '#1f2430',
        'border-color': 'rgba(180, 195, 215, 0.10)',
        'th-text-color': '#b6c0cf',
        'td-text-color': '#e8edf5',
      },
      Tag: {
        borderRadius: '4px',
      },
    }
  } else {
    return {
      common: {
        primaryColor: '#7f6ab3',
        primaryColorHover: '#8d78bf',
        primaryColorPressed: '#6f5aa3',
        primaryColorSuppl: '#7f6ab380',
        bodyColor: '#f7f8fb',
        cardColor: '#ffffff',
        modalColor: '#ffffff',
        popoverColor: '#ffffff',
        tableColor: '#ffffff',
        inputColor: '#ffffff',
         inputColorDisabled: '#f5f5f7',
         actionColor: '#f5f5f7',
         hoverColor: '#f0f0f2',
         borderColor: 'rgba(45, 55, 72, 0.13)',
         dividerColor: 'rgba(45, 55, 72, 0.08)',
        textColorBase: '#1a1d26',
        textColor1: '#1a1d26',
        textColor2: '#5a6478',
        textColor3: '#8e97a8',
        placeholderColor: '#8e97a8',
        iconColor: '#5a6478',
      },
      Button: {
        colorPrimary: '#7f6ab3',
        colorHoverPrimary: '#8d78bf',
        colorPressedPrimary: '#6f5aa3',
        textColorPrimary: '#ffffff',
        borderRadiusMedium: '6px',
      },
      Input: {
        color: '#ffffff',
        colorFocus: '#ffffff',
         borderColor: 'rgba(45, 55, 72, 0.13)',
         borderColorHover: 'rgba(45, 55, 72, 0.20)',
        borderColorFocus: '#7f6ab3',
        textColor: '#1a1d26',
        placeholderColor: '#8e97a8',
      },
      Card: {
        color: '#ffffff',
        colorModal: '#ffffff',
         borderColor: 'rgba(45, 55, 72, 0.08)',
      },
      Drawer: {
        color: '#ffffff',
         headerBorderBottom: 'rgba(45, 55, 72, 0.08)',
         bodyBorderBottom: 'rgba(45, 55, 72, 0.08)',
      },
      Tooltip: {
        color: '#ffffff',
        textColor: '#1a1d26',
        borderRadius: '6px',
      },
      Message: {
        colorSuccess: '#ffffff',
        colorError: '#ffffff',
        colorWarning: '#ffffff',
        colorInfo: '#ffffff',
        textColorSuccess: '#1a1d26',
        textColorError: '#1a1d26',
        textColorWarning: '#1a1d26',
        textColorInfo: '#1a1d26',
      },
      Empty: {
        textColor: '#8e97a8',
      },
      DataTable: {
        'th-color': '#eef1f6',
        'td-color': '#ffffff',
         'border-color': 'rgba(45, 55, 72, 0.08)',
        'th-text-color': '#5a6478',
        'td-text-color': '#1a1d26',
      },
      Tag: {
        borderRadius: '4px',
      },
    }
  }
})
</script>

<template>
  <n-config-provider :theme="theme" :theme-overrides="themeOverrides">
    <n-message-provider>
      <n-dialog-provider>
        <main-layout />
      </n-dialog-provider>
    </n-message-provider>
  </n-config-provider>
</template>
