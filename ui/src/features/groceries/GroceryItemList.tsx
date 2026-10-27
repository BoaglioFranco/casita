import { useEffect, useState } from 'react';
import { Alert, Paper, Stack, Text, ThemeIcon } from '@mantine/core';
import { ShoppingBasket } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { useDeleteGroceryItem } from './hooks';
import type { GroceryItem } from './types';
import { GroceryItemRow } from './GroceryItemRow';

type GroceryItemListProps = {
  householdId: string;
  listId: string;
  items: GroceryItem[];
};

export function GroceryItemList({ householdId, listId, items }: GroceryItemListProps) {
  const deleteItem = useDeleteGroceryItem(householdId, listId);
  const [openedItemId, setOpenedItemId] = useState<string | null>(null);

  useEffect(() => {
    if (openedItemId && !items.some((item) => item.id === openedItemId)) {
      setOpenedItemId(null);
    }
  }, [items, openedItemId]);

  if (items.length === 0) {
    return (
      <Paper p="xl" radius="sm" bg="gray.0">
        <Stack align="center" gap={8}>
          <ThemeIcon variant="transparent" color="gray" size="lg">
            <ShoppingBasket size={28} />
          </ThemeIcon>
          <Text fw={650}>No items yet</Text>
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
        <GroceryItemRow
          key={item.id}
          item={item}
          opened={openedItemId === item.id}
          onOpen={() => setOpenedItemId(item.id)}
          onToggle={() => setOpenedItemId((current) => current === item.id ? null : item.id)}
          onClose={() => setOpenedItemId((current) => current === item.id ? null : current)}
          onDelete={() => deleteItem.mutate(item.id)}
          deleting={deleteItem.isPending && deleteItem.variables === item.id}
          deleteDisabled={deleteItem.isPending}
        />
      ))}
    </Stack>
  );
}
