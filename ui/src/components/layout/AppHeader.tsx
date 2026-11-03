import {
  AppShell,
  Avatar,
  Button,
  Container,
  Group,
  Menu,
  Text,
  UnstyledButton,
  useMantineTheme,
} from '@mantine/core';
import { useQuery } from '@tanstack/react-query';
import { Link } from '@tanstack/react-router';
import { ChevronDown, House, Settings, UsersRound } from 'lucide-react';
import { currentMemberQuery } from '../../features/households/queries';
import { FeatureNavigation } from './FeatureNavigation';

export function AppHeader({
  showAccountMenu = true,
  householdId,
}: {
  showAccountMenu?: boolean;
  householdId?: string;
}) {
  const theme = useMantineTheme();
  const primaryColor = theme.colors[theme.primaryColor][6];
  const memberQuery = useQuery({
    ...currentMemberQuery(householdId ?? ''),
    enabled: showAccountMenu && Boolean(householdId),
  });
  const memberName = memberQuery.data?.name ?? 'Account';
  const memberInitial = memberName.charAt(0).toUpperCase();

  return (
    <AppShell.Header bg={primaryColor} withBorder={false}>
      <Container size="lg" h={{ base: 72, sm: 80 }} px="md">
        <Group h="100%" justify="space-between" wrap="nowrap">
          <UnstyledButton component={Link} to="/" c="white" aria-label="Go to Casita home">
            <Group gap="xs" align="center" wrap="nowrap">
              <House size={32} strokeWidth={1.8} aria-hidden />
              <Text component="span" fz={32} fw={800} lh={1}>
                Casita
              </Text>
            </Group>
          </UnstyledButton>

          {showAccountMenu && (
            <Menu position="bottom-end" shadow="md" width={220}>
              <Menu.Target>
                <Button
                  variant="subtle"
                  color="white"
                  h={48}
                  px="xs"
                  aria-label={`Open menu for ${memberName}`}
                  leftSection={
                    <Avatar
                      name={memberInitial}
                      color="initials"
                      allowedInitialsColors={['orange', 'teal', 'lime', 'pink', 'yellow']}
                      variant="filled"
                      autoContrast
                      radius="xl"
                      size="md"
                    />
                  }
                  rightSection={<ChevronDown size={16} />}
                >
                  {memberName}
                </Button>
              </Menu.Target>

              <Menu.Dropdown>
                <Menu.Label>Account</Menu.Label>
                <Menu.Item
                  renderRoot={(props) => <Link {...props} to="/profile" search={{ householdId }} />}
                  leftSection={<Settings size={16} />}
                >
                  Profile settings
                </Menu.Item>
                <Menu.Item leftSection={<UsersRound size={16} />} disabled>
                  Manage household
                </Menu.Item>
              </Menu.Dropdown>
            </Menu>
          )}
        </Group>
      </Container>
      {showAccountMenu && householdId && (
        <div
          style={{
            background: 'var(--mantine-color-body)',
            borderBottom: '1px solid var(--mantine-color-default-border)',
          }}
        >
          <FeatureNavigation householdId={householdId} />
        </div>
      )}
    </AppShell.Header>
  );
}
