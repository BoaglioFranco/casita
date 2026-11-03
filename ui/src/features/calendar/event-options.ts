import { formatDuration } from './calendar-utils';

export const durationOptions = [
  ...Array.from({ length: 96 }, (_, index) => {
    const minutes = (index + 1) * 15;
    return { value: String(minutes), label: formatDuration(minutes) };
  }),
  { value: '2880', label: '2 days' },
  { value: '4320', label: '3 days' },
  { value: '10080', label: '1 week' },
];

export const repeatOptions = [
  { value: '', label: 'Does not repeat' },
  { value: 'FREQ=DAILY', label: 'Every day' },
  { value: 'FREQ=WEEKLY', label: 'Every week' },
  { value: 'FREQ=WEEKLY;INTERVAL=2', label: 'Every two weeks' },
  { value: 'FREQ=MONTHLY', label: 'Every month on this date' },
  { value: 'custom', label: 'Custom repeat' },
];

export const reminderOptions = [
  { value: '0', label: 'At start' },
  ...[5, 10, 20, 30, 45].map((minutes) => ({
    value: String(minutes),
    label: `${minutes} minutes`,
  })),
  ...Array.from({ length: 8 }, (_, index) => {
    const hours = index + 1;
    return { value: String(hours * 60), label: `${hours} hour${hours === 1 ? '' : 's'}` };
  }),
];
