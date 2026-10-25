import { createRoute } from '@tanstack/react-router';
import { ProfilePage } from '../features/profile/ProfilePage';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  component: ProfilePage,
});
