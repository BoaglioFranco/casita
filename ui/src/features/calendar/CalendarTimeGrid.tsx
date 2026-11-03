import { useEffect, useLayoutEffect, useMemo, useRef, useState, type CSSProperties } from 'react';
import { useReducedMotion } from '@mantine/hooks';
import { Repeat2 } from 'lucide-react';
import {
  addDays,
  calendarViewRange,
  dateKey,
  formatDuration,
  formatTime,
  layoutDay,
  occurrenceKey,
} from './calendar-utils';
import type { CalendarOccurrence, CalendarView } from './types';
import classes from './CalendarTimeGrid.module.css';

const HOUR_HEIGHT = 64;
const hours = Array.from({ length: 24 }, (_, hour) => hour);

export function CalendarTimeGrid({
  date,
  view,
  currentTimeRequest,
  occurrences,
  onSelectDate,
  onCreate,
  onSelectEvent,
}: {
  date: Date;
  view: CalendarView;
  currentTimeRequest: number;
  occurrences: CalendarOccurrence[];
  onSelectDate: (date: Date) => void;
  onCreate: (date: Date) => void;
  onSelectEvent: (occurrence: CalendarOccurrence) => void;
}) {
  const scroll = useRef<HTMLDivElement>(null);
  const content = useRef<HTMLDivElement>(null);
  const headers = useRef<HTMLDivElement>(null);
  const reducedMotion = useReducedMotion();
  const rangeStart = calendarViewRange(date, view).start.getTime();
  const previousRange = useRef({ start: rangeStart, view });
  const [now, setNow] = useState(() => new Date());
  const days = useMemo(() => {
    const { start, dayCount } = calendarViewRange(date, view);
    return Array.from({ length: dayCount }, (_, index) => addDays(start, index));
  }, [date, view]);

  useLayoutEffect(() => {
    const previous = previousRange.current;
    previousRange.current = { start: rangeStart, view };
    if (reducedMotion || previous.view !== view || previous.start === rangeStart) return;

    const direction = rangeStart > previous.start ? 1 : -1;
    const days = content.current?.querySelectorAll<HTMLElement>(
      `.${classes.columnHeader}, .${classes.dayColumn}`,
    );
    // Animate the existing columns so navigation preserves scroll position and
    // immediately uses the latest date, even when previous animations are interrupted.
    const animations = Array.from(days ?? []).map((day) =>
      day.animate(
        [
          { transform: `translateX(${direction * 24}px)`, opacity: 0.3 },
          { transform: 'translateX(0)', opacity: 1 },
        ],
        { duration: 220, easing: 'cubic-bezier(0.2, 0, 0, 1)' },
      ),
    );
    return () => animations.forEach((animation) => animation.cancel());
  }, [rangeStart, view, reducedMotion]);

  useLayoutEffect(() => {
    const viewport = scroll.current;
    if (!viewport) return;
    const current = new Date();
    setNow(current);
    const currentTimeOffset = (current.getHours() + current.getMinutes() / 60) * HOUR_HEIGHT;
    const visibleTimelineHeight = viewport.clientHeight - (headers.current?.offsetHeight ?? 0);
    // Center the current time below the sticky day headers. The browser clamps
    // the offset near midnight; ordinary navigation keeps the user's scroll position.
    viewport.scrollTop = Math.max(0, currentTimeOffset - visibleTimelineHeight / 2);
  }, [currentTimeRequest]);

  useEffect(() => {
    const timer = window.setInterval(() => setNow(new Date()), 60_000);
    return () => window.clearInterval(timer);
  }, []);

  return (
    <div className={classes.calendar}>
      <div
        className={classes.scroll}
        ref={scroll}
        tabIndex={0}
        role="region"
        aria-label="Calendar time grid"
      >
        <div
          className={classes.inner}
          ref={content}
          data-week={view === 'week' || undefined}
          style={
            {
              '--day-count': days.length,
              '--hour-height': `${HOUR_HEIGHT}px`,
            } as CSSProperties
          }
        >
          <div className={classes.headers} ref={headers}>
            <div className={classes.timeGutter} aria-hidden />
            {days.map((day) => (
              <button
                key={dateKey(day)}
                type="button"
                className={classes.columnHeader}
                data-today={dateKey(day) === dateKey(now) || undefined}
                onClick={() => onSelectDate(day)}
                aria-label={`Show ${day.toLocaleDateString(undefined, { weekday: 'long', month: 'long', day: 'numeric' })}`}
              >
                <span>{day.toLocaleDateString(undefined, { weekday: 'short' })}</span>
                <strong>
                  {view === 'day'
                    ? day.toLocaleDateString(undefined, {
                        month: 'short',
                        day: 'numeric',
                      })
                    : day.getDate()}
                </strong>
              </button>
            ))}
          </div>
          <div className={classes.body}>
            <div className={classes.timeGutter} aria-hidden>
              {hours.map((hour) => (
                <span key={hour} style={{ top: hour * HOUR_HEIGHT }}>
                  {String(hour).padStart(2, '0')}:00
                </span>
              ))}
            </div>
            {days.map((day) => {
              const segments = layoutDay(occurrences, day);
              const isToday = dateKey(day) === dateKey(now);
              return (
                <div
                  className={classes.dayColumn}
                  key={dateKey(day)}
                  data-today={isToday || undefined}
                >
                  {hours.map((hour) => (
                    <button
                      key={hour}
                      type="button"
                      className={classes.timeSlot}
                      style={{ top: hour * HOUR_HEIGHT }}
                      tabIndex={hour === 9 ? 0 : -1}
                      aria-label={`Create event on ${day.toLocaleDateString()} at ${String(hour).padStart(2, '0')}:00`}
                      onClick={(event) => {
                        const start = new Date(day);
                        const rect = event.currentTarget.getBoundingClientRect();
                        const minute =
                          event.detail > 0 && event.clientY - rect.top >= HOUR_HEIGHT / 2 ? 30 : 0;
                        start.setHours(hour, minute, 0, 0);
                        onCreate(start);
                      }}
                    />
                  ))}
                  {segments.map((segment) => {
                    const { occurrence, startMinute, endMinute, column, columns } = segment;
                    const height = Math.max(((endMinute - startMinute) / 60) * HOUR_HEIGHT, 24);
                    const time = `${segment.continuesBefore ? '← ' : ''}${formatTime(occurrence.startsAt)} – ${formatTime(occurrence.endsAt)}${segment.continuesAfter ? ' →' : ''}`;
                    const duration =
                      occurrence.event.durationMinutes == null
                        ? ''
                        : formatDuration(occurrence.event.durationMinutes);
                    return (
                      <button
                        type="button"
                        key={occurrenceKey(occurrence)}
                        className={classes.event}
                        style={{
                          top: (startMinute / 60) * HOUR_HEIGHT,
                          height: Math.min(
                            height,
                            24 * HOUR_HEIGHT - (startMinute / 60) * HOUR_HEIGHT,
                          ),
                          left: `calc(${(column / columns) * 100}% + 3px)`,
                          width: `calc(${100 / columns}% - 6px)`,
                        }}
                        onClick={() => onSelectEvent(occurrence)}
                        title={`${occurrence.event.name}\n${time} · ${duration}`}
                        aria-label={`${occurrence.event.name}, ${time}, ${duration}${occurrence.event.recurrenceRule ? ', repeating' : ''}`}
                      >
                        <span className={classes.eventTitle}>
                          {occurrence.event.recurrenceRule && <Repeat2 size={12} aria-hidden />}
                          {occurrence.event.name}
                        </span>
                        {height >= 40 && <span className={classes.eventTime}>{time}</span>}
                        {height >= 62 && <span className={classes.eventDuration}>{duration}</span>}
                      </button>
                    );
                  })}
                  {isToday && (
                    <div
                      className={classes.now}
                      aria-label={`Current time ${formatTime(now)}`}
                      style={{
                        top: (now.getHours() + now.getMinutes() / 60) * HOUR_HEIGHT,
                      }}
                    />
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
