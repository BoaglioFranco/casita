import { apiRequest } from '../../lib/api-client';
import type { CalendarResponse, CalendarWindow, CreateEventInput } from './types';

export function getCalendar(window: CalendarWindow, signal?: AbortSignal) {
  const search = new URLSearchParams(window);
  return apiRequest<CalendarResponse>(`/calendar?${search}`, { signal });
}

export function createEvent(input: CreateEventInput): Promise<void> {
  return apiRequest(`/calendar/events`, {
    method: 'POST',
    body: JSON.stringify(input),
  });
}
