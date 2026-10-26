import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { Alert, Button, Modal, Stack, Text } from '@mantine/core';
import { Bell } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { isMobileDevice, isStandaloneApp } from './api';
import { useEnableNotifications } from './hooks';
import { notificationStateQuery } from './queries';

function promptStorageKey(): string {
  return `casita:notification-prompt-dismissed:${isStandaloneApp() ? 'standalone' : 'browser'}`;
}

function wasPromptDismissed(): boolean {
  try {
    return window.localStorage.getItem(promptStorageKey()) === 'true';
  } catch {
    return false;
  }
}

export function NotificationPrompt({ householdId }: { householdId: string }) {
  const state = useQuery(notificationStateQuery);
  const enable = useEnableNotifications(householdId);
  const [dismissed, setDismissed] = useState(wasPromptDismissed);
  const isNonInstalledMobileBrowser = isMobileDevice() && !isStandaloneApp();
  const shouldPrompt = state.data === 'disabled' && !isNonInstalledMobileBrowser;

  function dismiss() {
    try {
      window.localStorage.setItem(promptStorageKey(), 'true');
    } catch {
      // Storage can be unavailable in private browsing; dismiss for this session anyway.
    }
    setDismissed(true);
  }

  function handleEnable() {
    enable.mutate(undefined, { onSuccess: dismiss });
  }

  return (
    <Modal
      opened={!dismissed && shouldPrompt}
      onClose={dismiss}
      title="Enable notifications?"
      centered
      size="sm"
    >
      <Stack gap="md">
        <Text c="dimmed" size="sm">
          Casita can let you know when someone in your household updates a list. You can change
          this later in Profile settings.
        </Text>

        {enable.isError && (
          <Alert color="red" variant="light" radius="sm">
            {getErrorMessage(enable.error, 'Unable to enable notifications.')}
          </Alert>
        )}

        <Stack gap="xs">
          <Button
            leftSection={<Bell size={16} />}
            loading={enable.isPending}
            onClick={handleEnable}
          >
            Enable notifications
          </Button>
          <Button variant="subtle" color="gray" onClick={dismiss} disabled={enable.isPending}>
            Not now
          </Button>
        </Stack>
      </Stack>
    </Modal>
  );
}
