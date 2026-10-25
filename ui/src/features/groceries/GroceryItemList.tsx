import { ActionIcon, Alert, Box, Group, Paper, Stack, Text, ThemeIcon, Tooltip } from '@mantine/core';
import { ShoppingBasket, Trash2 } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { useDeleteGroceryItem } from './hooks';
import type { GroceryItem } from './types';

type GroceryItemListProps = {
  householdId: number;
  listId: number;
  items: GroceryItem[];
};

export function GroceryItemList({ householdId, listId, items }: GroceryItemListProps) {
  const deleteItem = useDeleteGroceryItem(householdId, listId);

  if (items.length === 0) {
    return (
      <Paper p="xl" radius="sm" bg="gray.0">
        <Stack align="center" gap={8}>
          <ThemeIcon variant="transparent" color="gray" size="lg">
            <ShoppingBasket size={28} />
          </ThemeIcon>
          <Text fw={650}>No items yet</Text>
          <Text c="dimmed" size="sm">
            Add the first grocery item above.
          </Text>
        </Stack>
      </Paper>
    );
  }

  return (
    <Stack gap="xs">
      {deleteItem.isError && (
        <Alert color="red" variant="light" radius="sm">
          {getErrorMessage(deleteItem.error, 'Unable to delete item.')}
        </Alert>
      )}

      {items.map((item) => (
        <Paper key={item.id} withBorder p="sm" radius="sm">
          <Group justify="space-between" align="center" gap="sm" wrap="nowrap">
            <Box miw={0}>
              <Text fw={700} style={{ overflowWrap: 'anywhere' }}>
                {item.name}
              </Text>
              {item.description && (
                <Text c="dimmed" size="sm" style={{ overflowWrap: 'anywhere' }}>
                  {item.description}
                </Text>
              )}
            </Box>

            <Tooltip label={`Delete ${item.name}`}>
              <ActionIcon
                aria-label={`Delete ${item.name}`}
                variant="subtle"
                color="red"
                loading={deleteItem.isPending && deleteItem.variables === item.id}
                disabled={deleteItem.isPending}
                onClick={() => deleteItem.mutate(item.id)}
              >
                <Trash2 size={18} />
              </ActionIcon>
            </Tooltip>
          </Group>
        </Paper>
      ))}
    </Stack>
  );
}
