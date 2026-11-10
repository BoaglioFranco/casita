import type { QueryClient } from '@tanstack/react-query';
import { redirect } from '@tanstack/react-router';
import { currentUserQuery } from '../features/auth/queries';
import { ApiError } from '../lib/api-client';

export async function requireAuthentication(queryClient: QueryClient) {
  try {
    const user = await queryClient.fetchQuery(currentUserQuery());
    return { user };
  } catch (error) {
    if (error instanceof ApiError && error.status === 401) {
      throw redirect({ to: '/login' });
    }

    throw error;
  }
}
