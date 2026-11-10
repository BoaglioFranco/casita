import { createRoute, redirect } from '@tanstack/react-router';
import { ListsEmptyState } from '../features/lists/ListsEmptyState';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  loader: async ({ context }) => {
    const { user, queryClient } = context;

    const lists = await queryClient.fetchQuery(listsQuery(user.id));
    if (lists[0]) {
      throw redirect({
        to: '/lists/$listId',
        params: { listId: lists[0].id },
        replace: true,
      });
    }
  },
  component: ListsEmptyState,
});
