import { createRoute, notFound } from '@tanstack/react-router';
import { GroceryPage } from '../features/groceries/GroceryPage';
import { groceryItemsQuery } from '../features/groceries/queries';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const listRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/$householdId/lists/$listId',
  beforeLoad: async ({ context, params }) => {
    const { households } = await requireAuthentication(context.queryClient);
    if (!households.some((household) => household.id === params.householdId)) {
      throw notFound();
    }
    return { householdId: params.householdId };
  },
  loader: async ({ context, params }) => {
    const lists = await context.queryClient.fetchQuery(listsQuery(params.householdId));
    const list = lists.find((candidate) => candidate.id === params.listId);
    if (!list) throw notFound();

    await context.queryClient.ensureQueryData(
      groceryItemsQuery(params.householdId, params.listId),
    );
    return { list };
  },
  component: () => {
    const { householdId, listId } = listRoute.useParams();
    const { list } = listRoute.useLoaderData();
    return <GroceryPage householdId={householdId} listId={listId} list={list} />;
  },
});
