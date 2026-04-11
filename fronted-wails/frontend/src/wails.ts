/**
 * Wails runtime bindings for CATendar
 * These replace the Electron IPC calls (setAlwaysOnTop, openExternal)
 */
import { WindowSetAlwaysOnTop, BrowserOpenURL, EventsOn } from '../wailsjs/runtime/runtime'

export function setAlwaysOnTop(on: boolean): void {
  WindowSetAlwaysOnTop(on)
}

export function openExternal(url: string): void {
  BrowserOpenURL(url)
}

export { EventsOn }
