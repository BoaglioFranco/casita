import { Container, Group, UnstyledButton } from '@mantine/core';
import { Link, useRouterState } from '@tanstack/react-router';
import { CalendarDays, ListChecks } from 'lucide-react';
import classes from './FeatureNavigation.module.css';

export function FeatureNavigation({ householdId }: { householdId: string }) {
  const path = useRouterState({ select: (state) => state.location.pathname });
  const calendarActive = path === `/${householdId}/calendar`;
  const listsActive = path.startsWith(`/${householdId}/lists`);

  return (
    <Container
      size="lg"
      px="md"
      component="nav"
      aria-label="Main features"
      className={classes.container}
    >
      <Group gap="lg" h="100%">
        <UnstyledButton
          renderRoot={(props) => (
            <Link {...props} to="/$householdId/lists" params={{ householdId }} />
          )}
          className={classes.link}
          data-active={listsActive || undefined}
          aria-current={listsActive ? 'page' : undefined}
        >
          <ListChecks size={18} aria-hidden /> Lists
        </UnstyledButton>
        <UnstyledButton
          renderRoot={(props) => (
            <Link {...props} to="/$householdId/calendar" params={{ householdId }} />
          )}
          className={classes.link}
          data-active={calendarActive || undefined}
          aria-current={calendarActive ? 'page' : undefined}
        >
          <CalendarDays size={18} aria-hidden /> Calendar
        </UnstyledButton>
      </Group>
    </Container>
  );
}
