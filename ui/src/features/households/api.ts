import { apiRequest } from '../../lib/api-client';
import type { CurrentMember, Household } from './types';

export async function getHouseholds(): Promise<Household[]> {
  const response = await apiRequest<{ households: Household[] }>('/households');
  return response.households;
}

export function getCurrentMember(householdId: string): Promise<CurrentMember> {
  return apiRequest(`/households/${householdId}/me`);
}
