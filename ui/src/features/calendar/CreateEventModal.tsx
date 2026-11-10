import { useEffect, useState, type FormEvent } from 'react';
import {
  Accordion,
  Alert,
  Button,
  Fieldset,
  Group,
  Modal,
  NativeSelect,
  Stack,
  Text,
  Textarea,
  TextInput,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { AlignLeft, Bell, Clock3 } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { dateKey, formatDuration, formatTime, timeKey } from './calendar-utils';
import { EventScheduleFields } from './EventScheduleFields';
import { reminderOptions, repeatOptions } from './event-options';
import { useCreateEvent } from './queries';

export function CreateEventModal({
  opened,
  initialStart,
  onClose,
  onExited,
  onCreated,
}: {
  opened: boolean;
  initialStart: Date | null;
  onClose: () => void;
  onExited: () => void;
  onCreated: (start: Date) => void;
}) {
  const mobile = useMediaQuery('(max-width: 47.99em)');
  const create = useCreateEvent();
  const [section, setSection] = useState<string | null>(null);
  const [pickerOpened, setPickerOpened] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [date, setDate] = useState('');
  const [time, setTime] = useState('');
  const [duration, setDuration] = useState(60);
  const [offset, setOffset] = useState(10);
  const [repeat, setRepeat] = useState('');
  const [customRule, setCustomRule] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    if (opened && initialStart) {
      setDate(dateKey(initialStart));
      setTime(timeKey(initialStart));
    }
  }, [opened, initialStart]);

  function resetAfterClose() {
    setName('');
    setDescription('');
    setDuration(60);
    setOffset(10);
    setRepeat('');
    setCustomRule('');
    setErrors({});
    setSection(null);
    setPickerOpened(false);
    create.reset();
    onExited();
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (create.isPending) return;
    const start = new Date(`${date}T${time}`);
    const nextErrors: Record<string, string> = {};
    if (!name.trim()) nextErrors.name = 'Give your event a name.';
    if (
      !date ||
      !time ||
      !Number.isFinite(start.getTime()) ||
      dateKey(start) !== date ||
      timeKey(start) !== time
    ) {
      nextErrors.start = 'Choose a valid local date and time.';
    }
    if (typeof duration !== 'number' || !Number.isSafeInteger(duration) || duration <= 0) {
      nextErrors.duration = 'Enter a positive whole number of minutes.';
    }
    if (typeof offset !== 'number' || !Number.isSafeInteger(offset) || offset < 0) {
      nextErrors.offset = 'Enter zero or a positive whole number of minutes.';
    }
    if (repeat === 'custom' && !customRule.trim()) nextErrors.rule = 'Enter a recurrence rule.';
    setErrors(nextErrors);
    if (nextErrors.start || nextErrors.duration || nextErrors.rule) setSection('schedule');
    else if (nextErrors.offset) setSection('reminder');
    if (Object.keys(nextErrors).length > 0) return;

    try {
      await create.mutateAsync({
        name: name.trim(),
        description: description.trim() || null,
        startsAt: start.toISOString(),
        durationMinutes: Number(duration),
        notificationOffsetMinutes: Number(offset),
        recurrenceRule: repeat === 'custom' ? customRule.trim() : repeat || null,
      });
      onCreated(start);
    } catch (error) {
      if (getErrorMessage(error, '') === 'Invalid recurrence rule') setSection('schedule');
    }
  }

  const start = new Date(`${date}T${time}`);
  const validStart = Number.isFinite(start.getTime());
  const end = new Date(start.getTime() + duration * 60_000);
  const scheduleDate = validStart
    ? start.toLocaleDateString(undefined, {
        weekday: 'short',
        month: 'short',
        day: 'numeric',
        ...(start.getFullYear() !== new Date().getFullYear() ? { year: 'numeric' as const } : {}),
      })
    : 'Choose a date and time';
  const endDate =
    validStart && dateKey(start) !== dateKey(end)
      ? `${end.toLocaleDateString(undefined, { month: 'short', day: 'numeric' })} `
      : '';
  const scheduleTime = validStart ? `${formatTime(start)} – ${endDate}${formatTime(end)}` : '';
  const scheduleDetail = `${formatDuration(duration)} · ${repeatOptions.find((option) => option.value === repeat)?.label ?? 'Does not repeat'}`;
  const reminderSummary =
    typeof offset !== 'number'
      ? 'Set reminder'
      : offset === 0
        ? 'At start of event'
        : `${formatDuration(offset)} before`;

  return (
    <Modal
      opened={opened}
      onClose={() => {
        if (!create.isPending) onClose();
      }}
      title="New event"
      onExitTransitionEnd={resetAfterClose}
      closeOnEscape={!pickerOpened}
      fullScreen={mobile}
      size="md"
      padding="lg"
      closeButtonProps={{ 'aria-label': 'Close new event' }}
      withCloseButton={!create.isPending}
      transitionProps={{
        transition: mobile ? 'slide-up' : 'pop',
        duration: 220,
        exitDuration: 160,
      }}
    >
      <form onSubmit={submit} noValidate>
        <Stack gap="md">
          <Fieldset disabled={create.isPending} variant="unstyled">
            <Stack gap="md">
              <TextInput
                aria-label="Event name"
                placeholder="Event name"
                required
                withAsterisk={false}
                size="md"
                radius="md"
                data-autofocus
                value={name}
                onChange={(e) => setName(e.currentTarget.value)}
                error={errors.name}
              />
              <Accordion value={section} onChange={setSection} transitionDuration={160}>
                <Accordion.Item value="schedule">
                  <Accordion.Control
                    px="xs"
                    mih={48}
                    icon={<Clock3 size={19} aria-hidden />}
                    aria-label={`Edit schedule: ${scheduleDate}, ${scheduleTime}, ${scheduleDetail}`}
                  >
                    <Group component="span" gap={8}>
                      <Text component="span" size="sm" fw={500} textWrap="nowrap">
                        {scheduleDate}
                      </Text>
                      {scheduleTime && (
                        <Text component="span" size="sm" textWrap="nowrap">
                          {scheduleTime}
                        </Text>
                      )}
                    </Group>
                    <Text component="span" display="block" size="xs" c="dimmed" mt={2}>
                      {scheduleDetail}
                    </Text>
                  </Accordion.Control>
                  <Accordion.Panel>
                    <EventScheduleFields
                      date={date}
                      time={time}
                      duration={duration}
                      repeat={repeat}
                      customRule={customRule}
                      errors={errors}
                      onDateChange={setDate}
                      onTimeChange={setTime}
                      onDurationChange={setDuration}
                      onRepeatChange={setRepeat}
                      onCustomRuleChange={setCustomRule}
                      active={opened && section === 'schedule'}
                      onPickerOpenChange={setPickerOpened}
                    />
                  </Accordion.Panel>
                </Accordion.Item>
                <Accordion.Item value="description">
                  <Accordion.Control
                    px="xs"
                    mih={48}
                    icon={<AlignLeft size={19} aria-hidden />}
                    aria-label={
                      description.trim()
                        ? `Edit description: ${description.trim()}`
                        : 'Add description'
                    }
                  >
                    <Text component="span" size="sm" lineClamp={2}>
                      {description.trim() || 'Add description'}
                    </Text>
                  </Accordion.Control>
                  <Accordion.Panel>
                    <Textarea
                      aria-label="Description"
                      placeholder="Description"
                      autosize
                      minRows={2}
                      maxRows={5}
                      size="md"
                      radius="md"
                      value={description}
                      onChange={(event) => setDescription(event.currentTarget.value)}
                    />
                  </Accordion.Panel>
                </Accordion.Item>
                <Accordion.Item value="reminder">
                  <Accordion.Control
                    px="xs"
                    mih={48}
                    icon={<Bell size={19} aria-hidden />}
                    aria-label={`Edit reminder: ${reminderSummary}`}
                  >
                    <Text component="span" size="sm">
                      {reminderSummary}
                    </Text>
                  </Accordion.Control>
                  <Accordion.Panel>
                    <NativeSelect
                      label="Remind before"
                      size="md"
                      radius="md"
                      data={reminderOptions}
                      value={String(offset)}
                      onChange={(event) => setOffset(Number(event.currentTarget.value))}
                      error={errors.offset}
                    />
                  </Accordion.Panel>
                </Accordion.Item>
              </Accordion>
            </Stack>
          </Fieldset>
          {create.isError && (
            <Alert color="red" title="Could not create event" role="alert">
              {getErrorMessage(create.error, 'Please try again. Your details have been kept.')}
            </Alert>
          )}
          <Group justify="flex-end" mt="xs">
            <Button variant="default" size="md" onClick={onClose} disabled={create.isPending}>
              Cancel
            </Button>
            <Button type="submit" size="md" loading={create.isPending}>
              Create event
            </Button>
          </Group>
        </Stack>
      </form>
    </Modal>
  );
}
