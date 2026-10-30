import { createRoute, redirect } from '@tanstack/react-router';
import { HouseholdEmptyState } from '../features/households/HouseholdEmptyState';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  loader: async ({ context }) => {
    const { householdId, queryClient } = context;
    if (!householdId) return;

    const lists = await queryClient.fetchQuery(listsQuery(householdId));
    if (lists[0]) {
      throw redirect({
        to: '/$householdId/lists/$listId',
        params: { householdId, listId: lists[0].id },
        replace: true,
      });
    }
  },
  component: () => {
    const { householdId } = indexRoute.useRouteContext();
    return <HouseholdEmptyState householdId={householdId} />;
  },
});
