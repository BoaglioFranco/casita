import { Center, Container, Stack, Text, Title } from '@mantine/core';

export function HouseholdEmptyState({ householdId }: { householdId?: string }) {
  return (
    <Center component="main" flex={1} w="100%">
      <Container size="sm" px="md" py="xl">
        <Stack gap="sm">
          <Title order={1}>{householdId ? 'No lists yet' : 'No households yet'}</Title>
          <Text c="dimmed">
            {householdId
              ? 'This household has no lists yet.'
              : 'Your account does not belong to a household yet.'}
          </Text>
        </Stack>
      </Container>
    </Center>
  );
}
