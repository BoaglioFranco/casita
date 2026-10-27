import { queryOptions } from '@tanstack/react-query';
import { getCurrentMember } from './api';

export const householdKeys = {
  all: ['households'] as const,
  currentMember: (householdId: string) =>
    [...householdKeys.all, householdId, 'current-member'] as const,
};

export function currentMemberQuery(householdId: string) {
  return queryOptions({
    queryKey: householdKeys.currentMember(householdId),
    queryFn: () => getCurrentMember(householdId),
  });
}
