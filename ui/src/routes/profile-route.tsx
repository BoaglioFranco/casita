import { createRoute, notFound } from '@tanstack/react-router';
import { ProfilePage } from '../features/profile/ProfilePage';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const profileRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/profile',
  validateSearch: (search: Record<string, unknown>): { householdId?: string } => ({
    householdId: typeof search.householdId === 'string' ? search.householdId : undefined,
  }),
  beforeLoad: async ({ context, search }) => {
    const authenticated = await requireAuthentication(context.queryClient);
    if (
      search.householdId &&
      !authenticated.households.some((household) => household.id === search.householdId)
    ) {
      throw notFound();
    }
    return { ...authenticated, householdId: search.householdId ?? authenticated.householdId };
  },
  component: () => {
    const { householdId } = profileRoute.useRouteContext();
    return <ProfilePage householdId={householdId} />;
  },
});
