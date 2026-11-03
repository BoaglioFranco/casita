import { createRoute, notFound } from '@tanstack/react-router';
import { ListsPage } from '../features/lists/ListsPage';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const listsRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/$householdId/lists',
  beforeLoad: async ({ context, params }) => {
    const { households } = await requireAuthentication(context.queryClient);
    if (!households.some((household) => household.id === params.householdId)) throw notFound();
    return { householdId: params.householdId };
  },
  loader: ({ context, params }) =>
    context.queryClient.ensureQueryData(listsQuery(params.householdId)),
  component: () => {
    const { householdId } = listsRoute.useParams();
    return <ListsPage householdId={householdId} />;
  },
});
