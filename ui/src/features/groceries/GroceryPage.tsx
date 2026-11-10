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
import { useCurrentUser } from '../auth/queries';
import { ListIcon } from '../list-icons/ListIcon';
import { listIconCatalog } from '../list-icons/icon-catalog';
import type { HouseholdList } from '../lists/types';
import { NotificationPrompt } from '../notifications/NotificationPrompt';
import { GroceryItemForm } from './GroceryItemForm';
import { GroceryItemList } from './GroceryItemList';
import { groceryItemsQuery } from './queries';

type GroceryPageProps = {
  listId: string;
  list: HouseholdList;
};

export function GroceryPage({ listId, list }: GroceryPageProps) {
  const { id: userId } = useCurrentUser();
  const itemsQuery = useSuspenseQuery(groceryItemsQuery(userId, listId));
  const icon = listIconCatalog.find((candidate) => candidate.id === list.icon);

  return (
    <Box component="main" flex={1} w="100%">
      <NotificationPrompt />
      <Container size="md" px="md" py="md">
        <Stack gap="md">
          <Group justify="space-between" align="center" gap="md" wrap="nowrap">
            <Group gap="sm" wrap="nowrap" style={{ minWidth: 0 }}>
              {icon && <ListIcon iconId={icon.id} size={32} style={{ flexShrink: 0 }} />}
              <Title
                order={1}
                fz={{ base: 32, sm: 48 }}
                lh={1}
                lts={0}
                style={{ overflowWrap: 'anywhere' }}
              >
                {list.name}
              </Title>
            </Group>

            <Tooltip label="Refresh list">
              <ActionIcon
                aria-label="Refresh list"
                variant="light"
                size="lg"
                style={{ flexShrink: 0 }}
                loading={itemsQuery.isFetching}
                onClick={() => itemsQuery.refetch()}
              >
                <RefreshCcw size={18} />
              </ActionIcon>
            </Tooltip>
          </Group>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="md">
              <GroceryItemForm listId={listId} />
              <Divider />
              <GroceryItemList listId={listId} items={itemsQuery.data} />
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
