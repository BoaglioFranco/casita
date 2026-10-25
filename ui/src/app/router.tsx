import { createRouter } from '@tanstack/react-router';
import { queryClient } from './query-client';
import { indexRoute } from '../routes/index-route';
import { loginRoute } from '../routes/login-route';
import { profileRoute } from '../routes/profile-route';
import { rootRoute } from '../routes/root-route';

const routeTree = rootRoute.addChildren([indexRoute, loginRoute, profileRoute]);

export const router = createRouter({
  routeTree,
  context: { queryClient },
  defaultPreload: 'intent',
  defaultPreloadStaleTime: 0,
  scrollRestoration: true,
});

declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
