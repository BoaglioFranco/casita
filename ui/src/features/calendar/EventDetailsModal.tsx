import { Group, Modal, Stack, Text, ThemeIcon } from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { Bell, Clock3, Repeat2 } from 'lucide-react';
import { formatDuration, formatTime } from './calendar-utils';
import type { CalendarOccurrence } from './types';

export function EventDetailsModal({
  occurrence,
  opened,
  onClose,
  onExited,
}: {
  occurrence: CalendarOccurrence | null;
  opened: boolean;
  onClose: () => void;
  onExited: () => void;
}) {
  const mobile = useMediaQuery('(max-width: 47.99em)');
  return (
    <Modal
      opened={opened}
      onClose={onClose}
      onExitTransitionEnd={onExited}
      title="Event details"
      size="sm"
      padding="lg"
      closeButtonProps={{ 'aria-label': 'Close event details' }}
      transitionProps={{
        transition: mobile ? 'slide-up' : 'pop',
        duration: 220,
        exitDuration: 160,
      }}
    >
      {occurrence && <EventDetailsContent occurrence={occurrence} />}
    </Modal>
  );
}

function EventDetailsContent({ occurrence }: { occurrence: CalendarOccurrence }) {
  const { event } = occurrence;
  const start = new Date(occurrence.startsAt);
  const end = new Date(occurrence.endsAt);
  return (
    <Stack gap="md">
      <Text fz="xl" fw={700} style={{ overflowWrap: 'anywhere' }}>
        {event.name}
      </Text>
      <Group wrap="nowrap" align="flex-start">
        <ThemeIcon variant="light">
          <Clock3 size={18} />
        </ThemeIcon>
        <div>
          <Text fw={500}>
            {start.toLocaleDateString(undefined, {
              weekday: 'long',
              month: 'long',
              day: 'numeric',
              year: 'numeric',
            })}
          </Text>
          <Text size="sm" c="dimmed">
            {formatTime(start)} –{' '}
            {start.toDateString() !== end.toDateString() && `${end.toLocaleDateString()} `}
            {formatTime(end)}
          </Text>
          {event.durationMinutes != null && (
            <Text size="sm" c="dimmed">
              {formatDuration(event.durationMinutes)}
            </Text>
          )}
        </div>
      </Group>
      <Group wrap="nowrap">
        <ThemeIcon variant="light">
          <Bell size={18} />
        </ThemeIcon>
        <Text size="sm">
          Reminder:{' '}
          {event.notificationOffsetMinutes === 0
            ? 'at the start'
            : `${formatDuration(event.notificationOffsetMinutes)} before`}
        </Text>
      </Group>
      {event.recurrenceRule && (
        <Group wrap="nowrap" align="flex-start">
          <ThemeIcon variant="light">
            <Repeat2 size={18} />
          </ThemeIcon>
          <div>
            <Text size="sm" fw={500}>
              Repeating event
            </Text>
            <Text size="xs" c="dimmed" style={{ overflowWrap: 'anywhere' }}>
              {event.recurrenceRule}
            </Text>
          </div>
        </Group>
      )}
      {event.description && (
        <Text size="sm" style={{ whiteSpace: 'pre-wrap', overflowWrap: 'anywhere' }}>
          {event.description}
        </Text>
      )}
    </Stack>
  );
}
