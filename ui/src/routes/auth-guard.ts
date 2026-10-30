import type { QueryClient } from '@tanstack/react-query';
import { redirect } from '@tanstack/react-router';
import { householdsQuery } from '../features/households/queries';
import { ApiError } from '../lib/api-client';

export async function requireAuthentication(queryClient: QueryClient) {
  try {
    const households = await queryClient.fetchQuery(householdsQuery());
    return { households, householdId: households[0]?.id };
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      throw redirect({ to: '/login' });
    }

    throw error;
  }
}
