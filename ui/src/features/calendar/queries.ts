import { infiniteQueryOptions, useMutation, useQueryClient } from '@tanstack/react-query';
import { useCurrentUser } from '../auth/queries';
import { createEvent, getCalendar } from './api';
import { adjacentWindow } from './calendar-utils';
import type { CalendarWindow, CreateEventInput } from './types';

export const calendarKeys = {
  user: (userId: string) => ['calendar', userId] as const,
};

export function calendarQuery(userId: string, window: CalendarWindow) {
  return infiniteQueryOptions({
    queryKey: [...calendarKeys.user(userId), window.from, window.to],
    initialPageParam: window,
    queryFn: async ({ pageParam, signal }) => ({
      ...(await getCalendar(pageParam, signal)),
      window: pageParam,
    }),
    getNextPageParam: (lastPage) => adjacentWindow(lastPage.window, 1),
    getPreviousPageParam: (firstPage) => adjacentWindow(firstPage.window, -1),
    // Navigating back must retain earlier occurrences, including multi-day events.
    gcTime: 30 * 60 * 1000,
  });
}

export function useCreateEvent() {
  const { id: userId } = useCurrentUser();
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateEventInput) => createEvent(input),
    onSuccess: () => {
      // A successful POST has no response body. Refetch the server-expanded timeline.
      void client.invalidateQueries({
        queryKey: calendarKeys.user(userId),
      });
    },
  });
}
