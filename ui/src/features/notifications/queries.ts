import { queryOptions } from '@tanstack/react-query';
import { getNotificationState } from './api';

export const notificationKeys = {
  all: ['notifications'] as const,
  state: (householdId: string) => [...notificationKeys.all, 'state', householdId] as const,
};

export function notificationStateQuery(householdId: string) {
  return queryOptions({
    queryKey: notificationKeys.state(householdId),
    queryFn: () => getNotificationState(householdId),
  });
}
