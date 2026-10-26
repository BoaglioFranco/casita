import { queryOptions } from '@tanstack/react-query';
import { getGroceryItems } from './api';

export const groceryKeys = {
  all: ['groceries'] as const,
  list: (householdId: string, listId: string) =>
    [...groceryKeys.all, 'list', householdId, listId] as const,
};

export function groceryItemsQuery(householdId: string, listId: string) {
  return queryOptions({
    queryKey: groceryKeys.list(householdId, listId),
    queryFn: () => getGroceryItems(householdId, listId),
  });
}
