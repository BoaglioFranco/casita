import { queryOptions } from '@tanstack/react-query';
import { getNotificationState } from './api';

export const notificationKeys = {
  all: ['notifications'] as const,
  state: () => [...notificationKeys.all, 'state'] as const,
};

export const notificationStateQuery = queryOptions({
  queryKey: notificationKeys.state(),
  queryFn: getNotificationState,
});
