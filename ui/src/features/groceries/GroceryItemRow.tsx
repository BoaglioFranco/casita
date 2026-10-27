import { useEffect, useRef, useState } from 'react';
import { ActionIcon, Box, Group, Paper, Popover, Stack, Text, Tooltip, UnstyledButton } from '@mantine/core';
import { useReducedMotion } from '@mantine/hooks';
import { Trash2 } from 'lucide-react';
import { formatAddedTime } from './format-added-time';
import type { GroceryItem } from './types';
import { useSwipeToDelete } from './use-swipe-to-delete';
import classes from './GroceryItemRow.module.css';

type GroceryItemRowProps = {
  item: GroceryItem;
  opened: boolean;
  onOpen: () => void;
  onToggle: () => void;
  onClose: () => void;
  onDelete: () => void;
  deleting: boolean;
  deleteDisabled: boolean;
};

export function GroceryItemRow({
  item, opened, onOpen, onToggle, onClose, onDelete, deleting, deleteDisabled,
}: GroceryItemRowProps) {
  const [now, setNow] = useState(Date.now);
  const closeTimer = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const pointerType = useRef('mouse');
  const reduceMotion = useReducedMotion();
  const swipe = useSwipeToDelete({
    disabled: deleteDisabled,
    onSwipeStart: () => {
      clearTimeout(closeTimer.current);
      onClose();
    },
    onDelete,
  });

  function openDetails() {
    clearTimeout(closeTimer.current);
    if (swipe.dragging) return;
    onOpen();
  }

  function scheduleClose() {
    clearTimeout(closeTimer.current);
    closeTimer.current = setTimeout(onClose, 150);
  }

  useEffect(() => () => clearTimeout(closeTimer.current), []);

  useEffect(() => {
    if (!opened) return;
    setNow(Date.now());
    const interval = window.setInterval(() => setNow(Date.now()), 60_000);
    return () => window.clearInterval(interval);
  }, [opened]);

  return (
    <Paper withBorder radius="sm" pos="relative" style={{ overflow: 'hidden' }}>
      <Box
        aria-hidden="true"
        pos="absolute"
        inset={0}
        bg={swipe.armed ? 'red.7' : 'red.1'}
        c={swipe.armed ? 'white' : 'red.7'}
        style={{ pointerEvents: 'none' }}
      >
        <Group h="100%" w="40%" ml="auto" justify="flex-end" gap="xs" px="xs" wrap="nowrap">
          {swipe.armed && <Text size="xs" fw={600} miw={0} ta="right">Release to delete</Text>}
          <Trash2 size={20} style={{ flexShrink: 0 }} />
        </Group>
      </Box>
      <Group
        {...swipe.handlers}
        pos="relative"
        bg="var(--mantine-color-body)"
        align="center" gap="sm" wrap="nowrap" pr="sm"
        style={{
          touchAction: 'pan-y',
          userSelect: 'none',
          transform: `translateX(-${swipe.offset}px)`,
          transition: swipe.dragging || reduceMotion ? 'none' : 'transform 180ms ease-out',
        }}
      >
        <Popover
          opened={opened}
          onDismiss={onClose}
          position="bottom-start"
          withArrow
          shadow="md"
          radius="sm"
          width={260}
          withinPortal
          closeOnClickOutside
          closeOnEscape
          middlewares={{ flip: true, shift: { padding: 12 } }}
        >
          <Popover.Target>
            <UnstyledButton
              type="button"
              flex={1}
              miw={0}
              p="sm"
              style={{ alignSelf: 'stretch', overflowWrap: 'anywhere' }}
              aria-label={`Details for ${item.name}`}
              onPointerEnter={(event) => {
                if (event.pointerType === 'mouse') openDetails();
              }}
              onPointerLeave={(event) => {
                if (event.pointerType === 'mouse') scheduleClose();
              }}
              onPointerDown={(event) => { pointerType.current = event.pointerType; }}
              onFocus={(event) => {
                if (event.currentTarget.matches(':focus-visible')) openDetails();
              }}
              onBlur={onClose}
              onClick={(event) => {
                if (event.detail === 0 || pointerType.current === 'mouse') openDetails();
                else onToggle();
              }}
            >
              <Text component="span" display="block" fw={700}>{item.name}</Text>
              {item.description && (
                <Text component="span" display="block" c="dimmed" size="sm">
                  {item.description}
                </Text>
              )}
            </UnstyledButton>
          </Popover.Target>
          <Popover.Dropdown
            style={{ maxWidth: 'calc(100vw - 24px)', overflowWrap: 'anywhere' }}
            onPointerEnter={(event) => {
              if (event.pointerType === 'mouse') openDetails();
            }}
            onPointerLeave={(event) => {
              if (event.pointerType === 'mouse') scheduleClose();
            }}
          >
            <Stack gap={4}>
              <Text size="sm" fw={600}>Added by {item.createdBy?.name || 'unknown member'}</Text>
              <Text size="sm" c="dimmed">{formatAddedTime(item.createdAt, now)}</Text>
            </Stack>
          </Popover.Dropdown>
        </Popover>
        <Tooltip label={`Delete ${item.name}`}>
          <ActionIcon
            className={classes.deleteButton}
            data-swipe-exclude
            aria-label={`Delete ${item.name}`}
            variant="subtle"
            color="red"
            loading={deleting}
            disabled={deleteDisabled}
            onClick={onDelete}
          >
            <Trash2 size={18} />
          </ActionIcon>
        </Tooltip>
      </Group>
    </Paper>
  );
}
