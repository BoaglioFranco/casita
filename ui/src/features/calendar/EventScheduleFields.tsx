import { useEffect, useState } from 'react';
import { ActionIcon, NativeSelect, SimpleGrid, Stack, TextInput } from '@mantine/core';
import { TimePicker } from '@mantine/dates';
import { Clock3 } from 'lucide-react';
import { CalendarDateInput } from './CalendarDateInput';
import { durationOptions, repeatOptions } from './event-options';

export function EventScheduleFields({
  date,
  time,
  duration,
  repeat,
  customRule,
  errors,
  onDateChange,
  onTimeChange,
  onDurationChange,
  onRepeatChange,
  onCustomRuleChange,
  active,
  onPickerOpenChange,
}: {
  date: string;
  time: string;
  duration: number;
  repeat: string;
  customRule: string;
  errors: Record<string, string>;
  onDateChange: (date: string) => void;
  onTimeChange: (time: string) => void;
  onDurationChange: (duration: number) => void;
  onRepeatChange: (repeat: string) => void;
  onCustomRuleChange: (rule: string) => void;
  active: boolean;
  onPickerOpenChange: (opened: boolean) => void;
}) {
  const [timeOpened, setTimeOpened] = useState(false);
  const [dateOpened, setDateOpened] = useState(false);

  useEffect(() => {
    if (!active) {
      setDateOpened(false);
      setTimeOpened(false);
    }
    onPickerOpenChange(active && (dateOpened || timeOpened));
    return () => onPickerOpenChange(false);
  }, [active, dateOpened, timeOpened, onPickerOpenChange]);

  return (
    <Stack gap="sm">
      <CalendarDateInput
        label="Date"
        required
        withAsterisk={false}
        value={date || null}
        onChange={(value) => {
          if (value) onDateChange(value);
        }}
        error={errors.start}
        onPickerOpenChange={setDateOpened}
      />
      <SimpleGrid cols={2} spacing="sm">
        <TimePicker
          label="Start time"
          required
          withAsterisk={false}
          size="md"
          radius="md"
          value={time}
          onChange={onTimeChange}
          format="24h"
          withDropdown
          minutesStep={15}
          styles={{ scrollarea: { flex: 1 }, control: { width: '100%', minHeight: 40 } }}
          hoursInputLabel="Start time hours"
          minutesInputLabel="Start time minutes"
          onFocus={() => setTimeOpened(true)}
          onKeyDownCapture={(event) => {
            if (event.key === 'Escape' && timeOpened) {
              event.stopPropagation();
              setTimeOpened(false);
            }
          }}
          rightSectionWidth={36}
          rightSectionPointerEvents="all"
          rightSection={
            <ActionIcon
              size={32}
              variant="transparent"
              c="dimmed"
              aria-label="Choose start time"
              aria-expanded={timeOpened}
              onMouseDown={(event) => event.preventDefault()}
              onClick={(event) => {
                event.stopPropagation();
                event.currentTarget.focus({ preventScroll: true });
                setTimeOpened((opened) => !opened);
              }}
            >
              <Clock3 size={18} aria-hidden />
            </ActionIcon>
          }
          popoverProps={{
            opened: timeOpened,
            onChange: setTimeOpened,
            withinPortal: false,
            width: 176,
            shadow: 'md',
            radius: 'md',
            trapFocus: false,
            returnFocus: true,
            transitionProps: { transition: 'pop', duration: 150 },
          }}
        />
        <NativeSelect
          label="Duration"
          required
          withAsterisk={false}
          size="md"
          radius="md"
          data={durationOptions}
          value={String(duration)}
          onChange={(event) => onDurationChange(Number(event.currentTarget.value))}
          error={errors.duration}
        />
      </SimpleGrid>
      <NativeSelect
        label="Repeat"
        size="md"
        radius="md"
        data={repeatOptions}
        value={repeat}
        onChange={(event) => onRepeatChange(event.currentTarget.value)}
      />
      {repeat === 'custom' && (
        <TextInput
          label="Recurrence rule"
          placeholder="FREQ=WEEKLY;BYDAY=MO,WE;COUNT=10"
          size="md"
          radius="md"
          value={customRule}
          onChange={(event) => onCustomRuleChange(event.currentTarget.value)}
          error={errors.rule}
        />
      )}
    </Stack>
  );
}
