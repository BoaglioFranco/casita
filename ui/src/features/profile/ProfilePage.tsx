import { Box, Container, Stack, Text, Title } from '@mantine/core';
import { appConfig } from '../../app/config';
import { NotificationSettings } from '../notifications/NotificationSettings';

export function ProfilePage() {
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

          <NotificationSettings householdId={appConfig.householdId} />
        </Stack>
      </Container>
    </Box>
  );
}
