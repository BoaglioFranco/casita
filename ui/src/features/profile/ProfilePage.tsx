import { Box, Container, Stack, Text, Title } from '@mantine/core';
import { NotificationSettings } from '../notifications/NotificationSettings';

export function ProfilePage({ householdId }: { householdId?: string }) {
  return (
    <Box component="main" flex={1} w="100%">
      <Container size="md" px="md" py="md">
        <Stack gap="lg">
          <div>
            <Title order={1} fz={{ base: 32, sm: 48 }} lh={1} lts={0}>
              Profile settings
            </Title>
            <Text c="dimmed" size="sm" mt="xs">
              Manage settings for your account and this device.
            </Text>
          </div>

          {householdId ? (
            <NotificationSettings householdId={householdId} />
          ) : (
            <Text c="dimmed">Join a household to manage its notifications.</Text>
          )}
        </Stack>
      </Container>
    </Box>
  );
}
