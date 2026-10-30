import { queryOptions } from '@tanstack/react-query';
import { getCurrentMember, getHouseholds } from './api';

export const householdKeys = {
  all: ['households'] as const,
  mine: ['households', 'mine'] as const,
  currentMember: (householdId: string) =>
    [...householdKeys.all, householdId, 'current-member'] as const,
};

export function householdsQuery() {
  return queryOptions({
    queryKey: householdKeys.mine,
    queryFn: getHouseholds,
  });
}

export function currentMemberQuery(householdId: string) {
  return queryOptions({
    queryKey: householdKeys.currentMember(householdId),
    queryFn: () => getCurrentMember(householdId),
  });
}
