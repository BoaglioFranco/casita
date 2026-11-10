import { useEffect, useMemo, useState } from 'react';
import {
  Alert,
  Button,
  Center,
  Container,
  Group,
  Loader,
  Stack,
  Text,
  Title,
  VisuallyHidden,
} from '@mantine/core';
import { useMediaQuery } from '@mantine/hooks';
import { useInfiniteQuery } from '@tanstack/react-query';
import { notifications } from '@mantine/notifications';
import { Plus } from 'lucide-react';
import { useCurrentUser } from '../auth/queries';
import { getErrorMessage } from '../../lib/api-client';
import {
  addDays,
  calendarViewRange,
  dateKey,
  initialWindow,
  mergeOccurrences,
  startOfDay,
  WINDOW_DAYS,
} from './calendar-utils';
import { calendarQuery } from './queries';
import { CalendarToolbar } from './CalendarToolbar';
import { CalendarTimeGrid } from './CalendarTimeGrid';
import { CreateEventModal } from './CreateEventModal';
import { EventDetailsModal } from './EventDetailsModal';
import type { CalendarOccurrence, CalendarView } from './types';

export function CalendarPage() {
  const { id: userId } = useCurrentUser();
  const desktop = useMediaQuery('(min-width: 48em)');
  const [date, setDate] = useState(() => startOfDay(new Date()));
  const [currentTimeRequest, setCurrentTimeRequest] = useState(0);
  const [desktopView, setDesktopView] = useState<'day' | 'week'>('week');
  const [mobileView, setMobileView] = useState<'day' | 'three-days'>('day');
  const view = desktop ? desktopView : mobileView;
  const [window, setWindow] = useState(() => initialWindow(new Date()));
  const calendar = useInfiniteQuery(calendarQuery(userId, window));
  const [creatingAt, setCreatingAt] = useState<Date | null>(null);
  const [createOpened, setCreateOpened] = useState(false);
  const [selected, setSelected] = useState<CalendarOccurrence | null>(null);
  const [detailsOpened, setDetailsOpened] = useState(false);
  const occurrences = useMemo(
    () => mergeOccurrences(calendar.data?.pages.map((page) => page.occurrences) ?? []),
    [calendar.data],
  );
  const { start: visibleStart, end: visibleEnd } = calendarViewRange(date, view);
  const pages = calendar.data?.pages;
  const firstPage = pages?.[0];
  const lastPage = pages?.[pages.length - 1];
  const covered = Boolean(
    firstPage &&
    lastPage &&
    new Date(firstPage.from) <= visibleStart &&
    new Date(lastPage.to) >= visibleEnd,
  );

  useEffect(() => {
    // Prefetch three days before the loaded edge, keeping every previous batch in cache.
    // Response boundaries are timestamps; request boundaries are server-local date strings.
    if (!firstPage || !lastPage || calendar.isFetching || calendar.isError) return;
    if (addDays(visibleStart, -3) < new Date(firstPage.from)) {
      void calendar.fetchPreviousPage();
    } else if (addDays(visibleEnd, 3) > new Date(lastPage.to)) {
      void calendar.fetchNextPage();
    }
  }, [
    firstPage,
    lastPage,
    calendar.isFetching,
    calendar.isError,
    calendar.fetchPreviousPage,
    calendar.fetchNextPage,
    visibleStart.getTime(),
    visibleEnd.getTime(),
  ]);

  function changeDate(nextDate: Date) {
    const next = startOfDay(nextDate);
    // A distant date jump starts a new four-week window instead of downloading
    // every intervening month. The earlier window remains in the query cache.
    if (
      firstPage &&
      lastPage &&
      (next < addDays(new Date(firstPage.from), -WINDOW_DAYS) ||
        next > addDays(new Date(lastPage.to), WINDOW_DAYS))
    )
      setWindow(initialWindow(next));
    setDate(next);
  }

  function changeView(nextView: CalendarView) {
    if (desktop && nextView !== 'three-days') setDesktopView(nextView);
    else if (!desktop && nextView !== 'week') setMobileView(nextView);
  }

  function openCreate() {
    const start = new Date(date);
    if (dateKey(date) === dateKey(new Date())) {
      const now = new Date();
      start.setHours(now.getHours(), Math.ceil((now.getMinutes() + 1) / 30) * 30, 0, 0);
    } else start.setHours(9, 0, 0, 0);
    startCreate(start);
  }

  function startCreate(start: Date) {
    setCreatingAt(start);
    setCreateOpened(true);
  }

  return (
    <Container
      component="main"
      size="lg"
      w="100%"
      px={{ base: 'sm', sm: 'md' }}
      py="lg"
      style={{ minWidth: 0 }}
    >
      <Stack gap="md">
        <Group justify="space-between" gap="sm" wrap="nowrap">
          <div>
            <Title order={1} fz={{ base: 30, sm: 36 }}>
              Calendar
            </Title>
          </div>
          <Button
            leftSection={<Plus size={18} />}
            size="md"
            onClick={openCreate}
            style={{ flexShrink: 0 }}
          >
            New event
          </Button>
        </Group>
        <CalendarToolbar
          date={date}
          view={view}
          onDateChange={changeDate}
          onToday={() => {
            changeDate(new Date());
            setCurrentTimeRequest((request) => request + 1);
          }}
          onViewChange={changeView}
          multiDayView={desktop ? 'week' : 'three-days'}
          refreshing={calendar.isFetching}
          onRefresh={() => {
            void calendar.refetch();
          }}
        />
        {calendar.isError && (
          <Alert color="red" title="Could not load calendar" role="alert">
            <Stack gap="xs">
              <Text size="sm">{getErrorMessage(calendar.error, 'Please try again.')}</Text>
              <Button
                size="xs"
                variant="light"
                color="red"
                w="fit-content"
                onClick={() => {
                  if (calendar.isFetchPreviousPageError) void calendar.fetchPreviousPage();
                  else if (calendar.isFetchNextPageError) void calendar.fetchNextPage();
                  else void calendar.refetch();
                }}
              >
                Try again
              </Button>
            </Stack>
          </Alert>
        )}
        {calendar.isPending ? (
          <Center mih={360}>
            <Loader aria-label="Loading calendar" />
          </Center>
        ) : (
          calendar.data && (
            <>
              <VisuallyHidden role="status" aria-live="polite">
                {!covered || calendar.isFetching ? 'Loading calendar' : 'Calendar updated'}
              </VisuallyHidden>
              <CalendarTimeGrid
                date={date}
                view={view}
                currentTimeRequest={currentTimeRequest}
                occurrences={occurrences}
                onSelectDate={(day) => {
                  setDate(day);
                  changeView('day');
                }}
                onCreate={startCreate}
                onSelectEvent={(occurrence) => {
                  setSelected(occurrence);
                  setDetailsOpened(true);
                }}
              />
            </>
          )
        )}
      </Stack>
      <CreateEventModal
        opened={createOpened}
        initialStart={creatingAt}
        onClose={() => setCreateOpened(false)}
        onExited={() => setCreatingAt(null)}
        onCreated={(start) => {
          setCreateOpened(false);
          changeDate(start);
          notifications.show({
            message: 'Event created',
            autoClose: 2500,
            withCloseButton: false,
            role: 'status',
          });
        }}
      />
      <EventDetailsModal
        occurrence={selected}
        opened={detailsOpened}
        onClose={() => setDetailsOpened(false)}
        onExited={() => setSelected(null)}
      />
    </Container>
  );
}
