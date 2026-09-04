import { WindowSetAlwaysOnTop, BrowserOpenURL } from '../wailsjs/runtime/runtime'

function hasWailsRuntime(): boolean {
  return 'runtime' in window
}

export function setAlwaysOnTop(on: boolean): boolean {
  if (!hasWailsRuntime()) return false
  WindowSetAlwaysOnTop(on)
  return true
}

export function openExternal(url: string): void {
  if (hasWailsRuntime()) {
    BrowserOpenURL(url)
    return
  }
  window.open(url, '_blank', 'noopener,noreferrer')
}
