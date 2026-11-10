export type CalendarView = 'day' | 'three-days' | 'week';

export type CalendarEvent = {
  id: string;
  householdId: string | null;
  createdByUserId: string | null;
  name: string;
  description: string | null;
  startsAt: string;
  durationMinutes: number | null;
  notificationOffsetMinutes: number;
  recurrenceRule: string | null;
  createdAt: string;
  updatedAt: string;
};

export type CalendarOccurrence = {
  event: CalendarEvent;
  startsAt: string;
  endsAt: string;
};

export type CalendarResponse = {
  from: string;
  to: string;
  occurrences: CalendarOccurrence[];
};

export type CreateEventInput = {
  name: string;
  description: string | null;
  startsAt: string;
  durationMinutes: number;
  notificationOffsetMinutes: number;
  recurrenceRule: string | null;
};

export type CalendarWindow = { from: string; to: string };
