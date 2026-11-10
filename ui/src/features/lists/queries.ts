import { queryOptions } from '@tanstack/react-query';
import { getLists } from './api';

export const listKeys = {
  user: (userId: string) => ['lists', userId] as const,
};

export function listsQuery(userId: string) {
  return queryOptions({
    queryKey: listKeys.user(userId),
    queryFn: () => getLists(),
  });
}
