import { createRoute } from '@tanstack/react-router';
import { CalendarPage } from '../features/calendar/CalendarPage';
import { requireAuthentication } from './auth-guard';
import { rootRoute } from './root-route';

export const calendarRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/calendar',
  beforeLoad: ({ context }) => requireAuthentication(context.queryClient),
  component: () => {
    const { user } = calendarRoute.useRouteContext();
    return <CalendarPage key={user.id} />;
  },
});
