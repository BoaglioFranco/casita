import type { QueryClient } from '@tanstack/react-query';
import { createRootRouteWithContext, Outlet, useRouterState } from '@tanstack/react-router';
import {
  Alert,
  AppShell,
  Button,
  Center,
  Container,
  Loader,
  Stack,
  Text,
  Title,
} from '@mantine/core';
import { getErrorMessage } from '../lib/api-client';
import { AppHeader } from '../components/layout/AppHeader';

export type RouterContext = { queryClient: QueryClient; householdId?: string };

function RootLayout() {
  const showAccountMenu = useRouterState({
    select: (state) => state.location.pathname !== '/login',
  });
  const householdId = useRouterState({
    select: (state) => state.matches[state.matches.length - 1]?.context.householdId,
  });

  return (
    <AppShell
      header={{
        height: {
          base: householdId && showAccountMenu ? 120 : 72,
          sm: householdId && showAccountMenu ? 128 : 80,
        },
      }}
    >
      <AppHeader showAccountMenu={showAccountMenu} householdId={householdId} />
      <AppShell.Main display="flex">
        <Outlet />
      </AppShell.Main>
    </AppShell>
  );
}

export const rootRoute = createRootRouteWithContext<RouterContext>()({
  component: RootLayout,
  pendingComponent: () => (
    <Center component="main" flex={1} w="100%">
      <Loader />
    </Center>
  ),
  errorComponent: ({ error, reset }) => (
    <Center component="main" flex={1} w="100%">
      <Container size="sm" w="100%" px="md">
        <Stack gap="md">
          <Title order={1}>Something went wrong</Title>
          <Alert color="red" variant="light">
            {getErrorMessage(error, 'The page could not be loaded.')}
          </Alert>
          <Button onClick={reset} w="fit-content">
            Try again
          </Button>
          <Text c="dimmed" size="sm">
            If the problem continues, check that the API is reachable.
          </Text>
        </Stack>
      </Container>
    </Center>
  ),
  notFoundComponent: () => (
    <Center component="main" flex={1} w="100%">
      <Stack align="center">
        <Title order={1}>Page not found</Title>
        <Text c="dimmed">That page does not exist.</Text>
      </Stack>
    </Center>
  ),
});
