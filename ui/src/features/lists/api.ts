import { apiRequest } from '../../lib/api-client';
import type { HouseholdList } from './types';

export async function getLists(householdId: string): Promise<HouseholdList[]> {
  const response = await apiRequest<{ lists: HouseholdList[] }>(
    `/households/${householdId}/lists`,
  );
  return response.lists;
}
