import { queryOptions, useSuspenseQuery } from '@tanstack/react-query';
import { getCurrentUser } from './api';

export function currentUserQuery() {
  return queryOptions({ queryKey: ['current-user'] as const, queryFn: getCurrentUser });
}

export function useCurrentUser() {
  return useSuspenseQuery(currentUserQuery()).data;
}
