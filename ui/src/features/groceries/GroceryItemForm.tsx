import { FormEvent, useState } from 'react';
import { Alert, Button, Group, Stack, TextInput } from '@mantine/core';
import { Plus, ShoppingBasket } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { useCreateGroceryItem } from './hooks';

type GroceryItemFormProps = { householdId: number; listId: number };

export function GroceryItemForm({ householdId, listId }: GroceryItemFormProps) {
  const [name, setName] = useState('');
  const createItem = useCreateGroceryItem(householdId, listId);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName) return;

    createItem.mutate({ name: trimmedName }, { onSuccess: () => setName('') });
  }

  return (
    <Stack gap="sm">
      <form onSubmit={handleSubmit}>
        <Group gap="sm" align="end" wrap="nowrap">
          <TextInput
            flex={1}
            miw={0}
            label="New item"
            placeholder="Type something"
            value={name}
            onChange={(event) => setName(event.currentTarget.value)}
            leftSection={<ShoppingBasket size={16} />}
            disabled={createItem.isPending}
          />
          <Button type="submit" leftSection={<Plus size={16} />} loading={createItem.isPending}>
            Add
          </Button>
        </Group>
      </form>

      {createItem.isError && (
        <Alert color="red" variant="light" radius="sm">
          {getErrorMessage(createItem.error, 'Unable to add item.')}
        </Alert>
      )}
    </Stack>
  );
}
