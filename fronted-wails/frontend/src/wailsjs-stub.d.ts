// This file provides type stubs for Wails-generated bindings.
// The actual wailsjs/go/calendar/CalendarApp.ts is generated at build time.

declare module '../../../wailsjs/go/calendar/CalendarApp' {
  import type { CalendarEvent, EventInput } from '@/types'

  export class CalendarApp {
    GetEvents(start: string, end: string): Promise<CalendarEvent[]>
    CreateEvent(input: EventInput): Promise<CalendarEvent>
    UpdateEvent(id: number, input: EventInput): Promise<CalendarEvent>
    DeleteEvent(id: number): Promise<void>
  }
}
