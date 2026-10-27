import { useSuspenseQuery } from '@tanstack/react-query';
import {
  ActionIcon,
  Box,
  Container,
  Divider,
  Group,
  Paper,
  Stack,
  Text,
  Title,
  Tooltip,
} from '@mantine/core';
import { RefreshCcw } from 'lucide-react';
import { appConfig } from '../../app/config';
import { NotificationPrompt } from '../notifications/NotificationPrompt';
import { GroceryItemForm } from './GroceryItemForm';
import { GroceryItemList } from './GroceryItemList';
import { groceryItemsQuery } from './queries';

export function GroceryPage() {
  const { householdId, listId } = appConfig;
  const itemsQuery = useSuspenseQuery(groceryItemsQuery(householdId, listId));

  return (
    <Box component="main" flex={1} w="100%">
      <NotificationPrompt householdId={householdId} />
      <Container size="md" px="md" py="md">
        <Stack gap="md">
          <Group justify="space-between" align="center" gap="md" wrap="nowrap">
            <Title order={1} fz={{ base: 32, sm: 48 }} lh={1} lts={0}>
              Groceries
            </Title>

            <Tooltip label="Refresh list">
              <ActionIcon
                aria-label="Refresh list"
                variant="light"
                size="lg"
                loading={itemsQuery.isFetching}
                onClick={() => itemsQuery.refetch()}
              >
                <RefreshCcw size={18} />
              </ActionIcon>
            </Tooltip>
          </Group>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="md">
              <GroceryItemForm householdId={householdId} listId={listId} />
              <Divider />
              <GroceryItemList householdId={householdId} listId={listId} items={itemsQuery.data} />
            </Stack>
          </Paper>

          <Text role="status" c="dimmed" size="sm">
            {itemsQuery.isFetching
              ? 'Refreshing groceries...'
              : `${itemsQuery.data.length} item${itemsQuery.data.length === 1 ? '' : 's'} loaded.`}
          </Text>
        </Stack>
      </Container>
    </Box>
  );
}
