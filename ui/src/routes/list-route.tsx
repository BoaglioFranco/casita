import { createRoute, notFound } from '@tanstack/react-router';
import { GroceryPage } from '../features/groceries/GroceryPage';
import { groceryItemsQuery } from '../features/groceries/queries';
import { listsQuery } from '../features/lists/queries';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const listRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/lists/$listId',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  loader: async ({ context, params }) => {
    const lists = await context.queryClient.fetchQuery(listsQuery(context.user.id));
    const list = lists.find((candidate) => candidate.id === params.listId);
    if (!list) throw notFound();
    await context.queryClient.ensureQueryData(groceryItemsQuery(context.user.id, params.listId));
    return { list };
  },
  component: () => {
    const { listId } = listRoute.useParams();
    const { list } = listRoute.useLoaderData();
    return <GroceryPage listId={listId} list={list} />;
  },
});
