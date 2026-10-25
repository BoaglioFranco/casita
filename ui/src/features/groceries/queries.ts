import { queryOptions } from '@tanstack/react-query';
import { getGroceryItems } from './api';

export const groceryKeys = {
  all: ['groceries'] as const,
  list: (householdId: number, listId: number) =>
    [...groceryKeys.all, 'list', householdId, listId] as const,
};

export function groceryItemsQuery(householdId: number, listId: number) {
  return queryOptions({
    queryKey: groceryKeys.list(householdId, listId),
    queryFn: () => getGroceryItems(householdId, listId),
  });
}
