import { queryOptions } from '@tanstack/react-query';
import { getNotificationState } from './api';

export const notificationKeys = {
  all: ['notifications'] as const,
  state: (userId: string) => [...notificationKeys.all, 'state', userId] as const,
};

export function notificationStateQuery(userId: string) {
  return queryOptions({
    queryKey: notificationKeys.state(userId),
    queryFn: () => getNotificationState(userId),
  });
}
