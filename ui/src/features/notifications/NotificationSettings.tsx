import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Group, Paper, Stack, Text, Title } from '@mantine/core';
import { Bell, BellRing } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import type { NotificationState } from './api';
import { useEnableNotifications } from './hooks';
import { notificationStateQuery } from './queries';

export function NotificationSettings({ householdId }: { householdId: number }) {
  const state = useQuery(notificationStateQuery);
  const enable = useEnableNotifications(householdId);
  const isEnabled = state.data === 'enabled';
  const canEnable = state.data === 'disabled';

  return (
    <Paper withBorder p={{ base: 'md', sm: 'lg' }} radius="sm">
      <Stack gap="md">
        <Group justify="space-between" align="flex-start" gap="md">
          <div>
            <Title order={2} fz="lg">
              Notifications
            </Title>
            <Text c="dimmed" size="sm" mt={4}>
              Get updates from your household on this device.
            </Text>
          </div>

          <Button
            variant={isEnabled ? 'light' : 'outline'}
            leftSection={isEnabled ? <BellRing size={16} /> : <Bell size={16} />}
            loading={enable.isPending || state.isPending}
            disabled={!canEnable}
            onClick={() => enable.mutate()}
          >
            {notificationButtonLabel(state.data)}
          </Button>
        </Group>

        {state.data === 'install-required' && (
          <Alert color="myColor" variant="light" radius="sm">
            On iPhone and iPad, add Casita to your Home Screen and open it from there to enable
            notifications.
          </Alert>
        )}

        {state.data === 'denied' && (
          <Alert color="orange" variant="light" radius="sm">
            Notifications are blocked in your browser settings. Allow them there, then reload
            Casita.
          </Alert>
        )}

        {state.data === 'unsupported' && (
          <Alert color="gray" variant="light" radius="sm">
            Notifications are not supported by this browser.
          </Alert>
        )}

        {state.isError && (
          <Alert color="red" variant="light" radius="sm">
            {getErrorMessage(state.error, 'Unable to check notification settings.')}
          </Alert>
        )}

        {enable.isError && (
          <Alert color="red" variant="light" radius="sm">
            {getErrorMessage(enable.error, 'Unable to enable notifications.')}
          </Alert>
        )}
      </Stack>
    </Paper>
  );
}

function notificationButtonLabel(state: NotificationState | undefined) {
  switch (state) {
    case 'enabled':
      return 'Notifications enabled';
    case 'denied':
      return 'Notifications blocked';
    case 'install-required':
      return 'Install Casita first';
    case 'unsupported':
      return 'Not supported';
    default:
      return 'Enable notifications';
  }
}
