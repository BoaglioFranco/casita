import { queryOptions } from '@tanstack/react-query';
import { getLists } from './api';

export const listKeys = {
  household: (householdId: string) => ['lists', householdId] as const,
};

export function listsQuery(householdId: string) {
  return queryOptions({
    queryKey: listKeys.household(householdId),
    queryFn: () => getLists(householdId),
  });
}
