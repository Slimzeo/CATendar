export interface CalendarEvent {
  id: number
  title: string
  start: string
  end?: string
  color: string
  allDay: boolean
  description: string
  bold?: boolean
  createdAt: any
  updatedAt: any
}

export interface EventInput {
  title: string
  start: string
  end?: string
  color: string
  allDay: boolean
  description: string
  bold?: boolean
}

export interface ColorPalette {
  id: string
  name: string
  colors: string[]
}

export interface DayCell {
  date: Date
  isCurrentMonth: boolean
  isToday: boolean
  events: CalendarEvent[]
}