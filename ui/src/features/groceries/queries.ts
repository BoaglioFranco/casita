import { queryOptions } from '@tanstack/react-query';
import { getGroceryItems } from './api';

export const groceryKeys = {
  all: ['groceries'] as const,
  list: (userId: string, listId: string) =>
    [...groceryKeys.all, 'list', userId, listId] as const,
};

export function groceryItemsQuery(userId: string, listId: string) {
  return queryOptions({
    queryKey: groceryKeys.list(userId, listId),
    queryFn: () => getGroceryItems(listId),
  });
}
