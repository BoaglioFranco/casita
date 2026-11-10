import { apiRequest } from '../../lib/api-client';
import type { HouseholdList } from './types';

export async function getLists(): Promise<HouseholdList[]> {
  const response = await apiRequest<{ lists: HouseholdList[] }>(
    '/lists',
  );
  return response.lists;
}
