import {
  ActionIcon,
  Button,
  Group,
  SegmentedControl,
  Text,
  Tooltip,
  useMantineTheme,
} from '@mantine/core';
import { ChevronLeft, ChevronRight, RefreshCcw } from 'lucide-react';
import { addDays, calendarViewRange, dateKey, parseDateKey } from './calendar-utils';
import { CalendarDateInput } from './CalendarDateInput';
import type { CalendarView } from './types';

export function CalendarToolbar({
  date,
  view,
  onDateChange,
  onToday,
  onViewChange,
  refreshing,
  onRefresh,
  multiDayView,
}: {
  date: Date;
  view: CalendarView;
  onDateChange: (date: Date) => void;
  onToday: () => void;
  onViewChange: (view: CalendarView) => void;
  refreshing: boolean;
  onRefresh: () => void;
  multiDayView: 'three-days' | 'week';
}) {
  const theme = useMantineTheme();
  const { dayCount } = calendarViewRange(date, view);
  const navigationLabel = view === 'three-days' ? '3 days' : view;
  return (
    <Group justify="space-between" gap="sm" align="center">
      <Group gap={4} wrap="nowrap">
        <Button
          variant="light"
          radius="md"
          h={44}
          px="sm"
          mr={4}
          onClick={onToday}
        >
          Today
        </Button>
        <ActionIcon
          variant="subtle"
          size={44}
          aria-label={`Previous ${navigationLabel}`}
          onClick={() => onDateChange(addDays(date, -dayCount))}
        >
          <ChevronLeft size={20} />
        </ActionIcon>
        <ActionIcon
          variant="subtle"
          size={44}
          aria-label={`Next ${navigationLabel}`}
          onClick={() => onDateChange(addDays(date, dayCount))}
        >
          <ChevronRight size={20} />
        </ActionIcon>
        <Text fw={600} fz={{ base: 'sm', sm: 'md' }}>
          {date.toLocaleDateString(undefined, {
            month: 'long',
            year: 'numeric',
          })}
        </Text>
      </Group>
      <Group gap={6} wrap="nowrap" ml="auto" w={{ base: '100%', sm: 'auto' }}>
        <CalendarDateInput
          aria-label="Jump to date"
          value={dateKey(date)}
          w={{ base: 128, sm: 154 }}
          valueFormat="MMM D, YY"
          size="sm"
          styles={{ input: { minHeight: 44 } }}
          onChange={(value) => {
            if (value) onDateChange(parseDateKey(value));
          }}
        />
        <SegmentedControl
          aria-label="Calendar view"
          value={view}
          size="sm"
          radius="md"
          color={theme.primaryColor}
          autoContrast
          ml="auto"
          w={112}
          h={44}
          flex="0 0 auto"
          withItemsBorders={false}
          styles={{
            label: {
              height: 36,
              padding: 0,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            },
          }}
          onChange={(value) => onViewChange(value as CalendarView)}
          data={[
            { value: 'day', label: 'Day' },
            { value: multiDayView, label: multiDayView === 'three-days' ? '3 days' : 'Week' },
          ]}
        />
        <Tooltip label="Refresh calendar">
          <ActionIcon
            aria-label="Refresh calendar"
            variant="subtle"
            size={44}
            loading={refreshing}
            onClick={onRefresh}
            flex="0 0 auto"
          >
            <RefreshCcw size={18} />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Group>
  );
}
