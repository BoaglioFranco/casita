import { infiniteQueryOptions, useMutation, useQueryClient } from '@tanstack/react-query';
import { createEvent, getCalendar } from './api';
import { adjacentWindow } from './calendar-utils';
import type { CalendarWindow, CreateEventInput } from './types';

export const calendarKeys = {
  household: (householdId: string) => ['calendar', householdId] as const,
};

export function calendarQuery(householdId: string, window: CalendarWindow) {
  return infiniteQueryOptions({
    queryKey: [...calendarKeys.household(householdId), window.from, window.to],
    initialPageParam: window,
    queryFn: async ({ pageParam, signal }) => ({
      ...(await getCalendar(householdId, pageParam, signal)),
      window: pageParam,
    }),
    getNextPageParam: (lastPage) => adjacentWindow(lastPage.window, 1),
    getPreviousPageParam: (firstPage) => adjacentWindow(firstPage.window, -1),
    // Navigating back must retain earlier occurrences, including multi-day events.
    gcTime: 30 * 60 * 1000,
  });
}

export function useCreateEvent(householdId: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateEventInput) => createEvent(householdId, input),
    onSuccess: () => {
      // A successful POST has no response body. Refetch the server-expanded timeline.
      void client.invalidateQueries({
        queryKey: calendarKeys.household(householdId),
      });
    },
  });
}
