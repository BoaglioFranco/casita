import { createRoute } from '@tanstack/react-router';
import { GroceryPage } from '../features/groceries/GroceryPage';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  component: GroceryPage,
});
