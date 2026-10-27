import { apiRequest } from '../../lib/api-client';
import type { CurrentMember } from './types';

export function getCurrentMember(householdId: string): Promise<CurrentMember> {
  return apiRequest(`/households/${householdId}/me`);
}
