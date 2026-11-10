import { createRoute } from '@tanstack/react-router';
import { ListsPage } from '../features/lists/ListsPage';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const listsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/lists',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  loader: ({ context }) => context.queryClient.ensureQueryData(listsQuery(context.user.id)),
  component: ListsPage,
});
