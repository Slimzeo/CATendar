/// <reference types="vite/client" />

declare global {
  interface Window {
    electronAPI?: {
      setAlwaysOnTop: (flag: boolean, level?: string) => void
    }
  }
}

export {}
