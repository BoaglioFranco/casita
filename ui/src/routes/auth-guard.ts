import type { QueryClient } from '@tanstack/react-query';
import { redirect } from '@tanstack/react-router';
import { appConfig } from '../app/config';
import { groceryItemsQuery } from '../features/groceries/queries';
import { ApiError } from '../lib/api-client';

export async function requireAuthentication(queryClient: QueryClient): Promise<void> {
  try {
    await queryClient.ensureQueryData(
      groceryItemsQuery(appConfig.householdId, appConfig.listId),
    );
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      throw redirect({ to: '/login' });
    }

    throw error;
  }
}
