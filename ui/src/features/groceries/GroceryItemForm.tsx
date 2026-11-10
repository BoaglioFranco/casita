import { FormEvent, useId, useState } from 'react';
import { Alert, Button, Collapse, Group, Stack, Textarea, TextInput } from '@mantine/core';
import { ChevronDown, ChevronUp, Plus, ShoppingBasket } from 'lucide-react';
import { getErrorMessage } from '../../lib/api-client';
import { useCreateGroceryItem } from './hooks';

type GroceryItemFormProps = { listId: string };

export function GroceryItemForm({ listId }: GroceryItemFormProps) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [descriptionExpanded, setDescriptionExpanded] = useState(false);
  const descriptionPanelId = useId();
  const createItem = useCreateGroceryItem(listId);

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName) return;

    createItem.mutate(
      { name: trimmedName, description: description.trim() || undefined },
      {
        onSuccess: () => {
          setName('');
          setDescription('');
          setDescriptionExpanded(false);
        },
      },
    );
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
        <Button
          type="button"
          variant="subtle"
          size="compact-sm"
          mt="xs"
          px="xs"
          ml="calc(-1 * var(--mantine-spacing-xs))"
          aria-expanded={descriptionExpanded}
          aria-controls={descriptionPanelId}
          rightSection={descriptionExpanded ? <ChevronUp size={16} /> : <ChevronDown size={16} />}
          disabled={createItem.isPending}
          onClick={() => setDescriptionExpanded((expanded) => !expanded)}
        >
          {descriptionExpanded ? 'Hide description' : 'Show description'}
        </Button>
        <Collapse expanded={descriptionExpanded} id={descriptionPanelId}>
          <Textarea
            mt="sm"
            aria-label="Description"
            placeholder="Description"
            autosize
            minRows={2}
            value={description}
            onChange={(event) => setDescription(event.currentTarget.value)}
            disabled={createItem.isPending}
          />
        </Collapse>
      </form>

      {createItem.isError && (
        <Alert color="red" variant="light" radius="sm">
          {getErrorMessage(createItem.error, 'Unable to add item.')}
        </Alert>
      )}
    </Stack>
  );
}
