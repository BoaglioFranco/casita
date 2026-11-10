import { Container, Group, Paper, SimpleGrid, Stack, Text, Title } from '@mantine/core';
import { useSuspenseQuery } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';
import { ChevronRight, ListChecks } from 'lucide-react';
import { useCurrentUser } from '../auth/queries';
import { ListIcon } from '../list-icons/ListIcon';
import { listIconCatalog } from '../list-icons/icon-catalog';
import { listsQuery } from './queries';

export function ListsPage() {
  const { id: userId } = useCurrentUser();
  const { data: lists } = useSuspenseQuery(listsQuery(userId));
  return (
    <Container component="main" size="md" w="100%" py="xl" px="md">
      <Stack gap="lg">
        <div>
          <Title order={1}>Lists</Title>
          <Text c="dimmed" mt={4}>
            Your household’s shared lists.
          </Text>
        </div>
        {lists.length === 0 ? (
          <Text c="dimmed">This household has no lists yet.</Text>
        ) : (
          <SimpleGrid cols={{ base: 1, sm: 2 }}>
            {lists.map((list) => {
              const icon = listIconCatalog.find((candidate) => candidate.id === list.icon);
              return (
                <Paper
                  key={list.id}
                  renderRoot={(props) => (
                    <Link
                      {...props}
                      to="/lists/$listId"
                      params={{ listId: list.id }}
                    />
                  )}
                  withBorder
                  p="lg"
                  c="inherit"
                  td="none"
                >
                  <Group wrap="nowrap">
                    {icon ? <ListIcon iconId={icon.id} size={28} /> : <ListChecks size={28} />}
                    <Text fw={600} flex={1} style={{ overflowWrap: 'anywhere' }}>
                      {list.name}
                    </Text>
                    <ChevronRight size={18} aria-hidden />
                  </Group>
                </Paper>
              );
            })}
          </SimpleGrid>
        )}
      </Stack>
    </Container>
  );
}
