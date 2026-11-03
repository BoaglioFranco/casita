import { apiRequest } from '../../lib/api-client';
import type { CalendarResponse, CalendarWindow, CreateEventInput } from './types';

export function getCalendar(householdId: string, window: CalendarWindow, signal?: AbortSignal) {
  const search = new URLSearchParams(window);
  return apiRequest<CalendarResponse>(`/households/${householdId}/calendar?${search}`, { signal });
}

export function createEvent(householdId: string, input: CreateEventInput): Promise<void> {
  return apiRequest(`/households/${householdId}/calendar/events`, {
    method: 'POST',
    body: JSON.stringify(input),
  });
}
