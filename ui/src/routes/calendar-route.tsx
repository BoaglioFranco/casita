import { createRoute, notFound } from '@tanstack/react-router';
import { CalendarPage } from '../features/calendar/CalendarPage';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const calendarRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/$householdId/calendar',
  beforeLoad: async ({ context, params }) => {
    const { households } = await requireAuthentication(context.queryClient);
    if (!households.some((household) => household.id === params.householdId)) throw notFound();
    return { householdId: params.householdId };
  },
  component: () => {
    const { householdId } = calendarRoute.useParams();
    return <CalendarPage key={householdId} householdId={householdId} />;
  },
});
