import type { CalendarOccurrence, CalendarView, CalendarWindow } from './types';

export const WINDOW_DAYS = 28;

export function startOfDay(date: Date): Date {
  return new Date(date.getFullYear(), date.getMonth(), date.getDate());
}

export function addDays(date: Date, count: number): Date {
  const result = new Date(date);
  result.setDate(result.getDate() + count);
  return result;
}

export function startOfWeek(date: Date): Date {
  return addDays(startOfDay(date), -((date.getDay() + 6) % 7));
}

export function calendarViewRange(date: Date, view: CalendarView) {
  const start = view === 'week' ? startOfWeek(date) : startOfDay(date);
  const dayCount = view === 'week' ? 7 : view === 'three-days' ? 3 : 1;
  return { start, end: addDays(start, dayCount), dayCount };
}

export function dateKey(date: Date): string {
  return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`;
}

export function parseDateKey(value: string): Date {
  const [year, month, day] = value.split('-').map(Number);
  return new Date(year, month - 1, day);
}

export function initialWindow(today: Date): CalendarWindow {
  return {
    from: dateKey(addDays(today, -14)),
    to: dateKey(addDays(today, 14)),
  };
}

export function adjacentWindow(window: CalendarWindow, direction: -1 | 1): CalendarWindow {
  const boundary = parseDateKey(direction === 1 ? window.to : window.from);
  const other = addDays(boundary, direction * WINDOW_DAYS);
  return direction === 1
    ? { from: dateKey(boundary), to: dateKey(other) }
    : { from: dateKey(other), to: dateKey(boundary) };
}

export function occurrenceKey(occurrence: CalendarOccurrence): string {
  return `${occurrence.event.id}:${new Date(occurrence.startsAt).getTime()}`;
}

export function mergeOccurrences(batches: CalendarOccurrence[][]): CalendarOccurrence[] {
  const occurrences = new Map<string, CalendarOccurrence>();
  for (const batch of batches) {
    for (const occurrence of batch) occurrences.set(occurrenceKey(occurrence), occurrence);
  }
  return [...occurrences.values()].sort(
    (a, b) => new Date(a.startsAt).getTime() - new Date(b.startsAt).getTime(),
  );
}

export function formatTime(value: string | Date): string {
  return new Date(value).toLocaleTimeString(undefined, {
    hour: '2-digit',
    minute: '2-digit',
    hourCycle: 'h23',
  });
}

// Native time inputs require ASCII HH:mm, regardless of the display locale.
export function timeKey(date: Date): string {
  return `${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
}

export function formatDuration(minutes: number): string {
  const hours = Math.floor(minutes / 60);
  const rest = minutes % 60;
  return [hours ? `${hours}h` : '', rest ? `${rest}m` : ''].filter(Boolean).join(' ') || '0m';
}

export type PositionedOccurrence = {
  occurrence: CalendarOccurrence;
  startMinute: number;
  endMinute: number;
  column: number;
  columns: number;
  continuesBefore: boolean;
  continuesAfter: boolean;
};

// Split at local midnight and give each connected overlap group its own lanes.
// Use the minimum rendered height when allocating lanes so short events remain clickable.
export function layoutDay(occurrences: CalendarOccurrence[], day: Date): PositionedOccurrence[] {
  const dayStart = startOfDay(day);
  const dayEnd = addDays(dayStart, 1);
  const minutes = (date: Date) => date.getHours() * 60 + date.getMinutes() + date.getSeconds() / 60;
  const segments = occurrences
    .flatMap((occurrence) => {
      const start = new Date(occurrence.startsAt);
      const end = new Date(occurrence.endsAt);
      if (!(start < dayEnd && end > dayStart && end > start)) return [];
      return [
        {
          occurrence,
          startMinute: start < dayStart ? 0 : minutes(start),
          endMinute: end >= dayEnd ? 1440 : minutes(end),
          column: 0,
          columns: 1,
          continuesBefore: start < dayStart,
          continuesAfter: end > dayEnd,
        },
      ];
    })
    .sort((a, b) => a.startMinute - b.startMinute || b.endMinute - a.endMinute);

  let group: PositionedOccurrence[] = [];
  let laneEnds: number[] = [];
  let groupEnd = -1;
  const finishGroup = () => {
    for (const segment of group) segment.columns = laneEnds.length;
  };
  for (const segment of segments) {
    if (segment.startMinute >= groupEnd) {
      finishGroup();
      group = [];
      laneEnds = [];
    }
    const renderedEnd = Math.max(segment.endMinute, segment.startMinute + 24);
    const availableLane = laneEnds.findIndex((end) => end <= segment.startMinute);
    segment.column = availableLane === -1 ? laneEnds.length : availableLane;
    laneEnds[segment.column] = renderedEnd;
    groupEnd = Math.max(...laneEnds);
    group.push(segment);
  }
  finishGroup();
  return segments;
}
