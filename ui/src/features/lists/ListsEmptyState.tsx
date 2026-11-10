import { Center, Container, Stack, Text, Title } from '@mantine/core';

export function ListsEmptyState() {
  return (
    <Center component="main" flex={1} w="100%">
      <Container size="sm" px="md" py="xl">
        <Stack gap="sm">
          <Title order={1}>No lists yet</Title>
          <Text c="dimmed">
            This household has no lists yet.
          </Text>
        </Stack>
      </Container>
    </Center>
  );
}
