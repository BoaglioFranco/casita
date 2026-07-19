import { FormEvent, useEffect, useState } from 'react';
import {
  ActionIcon,
  Alert,
  Badge,
  Button,
  Container,
  Divider,
  Group,
  Loader,
  Paper,
  PasswordInput,
  Stack,
  Text,
  TextInput,
  Title,
  Tooltip,
} from '@mantine/core';
import {
  Bell,
  BellRing,
  Check,
  LockKeyhole,
  LogIn,
  Mail,
  Plus,
  RefreshCcw,
  ShoppingBasket,
  Trash2,
} from 'lucide-react';

const DEMO_HOUSEHOLD_ID = 1;
const DEMO_LIST_ID = 1;
const API_URL = '/api';

type GroceryItem = {
  id: number;
  listId?: number;
  name: string;
  description?: string | null;
  createdAt?: string;
  updatedAt?: string;
};

type RawGroceryItem = {
  id?: number;
  ID?: number;
  listId?: number;
  ListID?: number;
  name?: string;
  Name?: string;
  description?: string | null;
  Description?: string | null;
  createdAt?: string;
  CreatedAt?: string;
  updatedAt?: string;
  UpdatedAt?: string;
};

function normalizeItems(data: unknown): GroceryItem[] {
  const rawItems = Array.isArray(data)
    ? data
    : Array.isArray((data as { items?: unknown[] })?.items)
      ? (data as { items: unknown[] }).items
      : [];

  return rawItems.reduce<GroceryItem[]>((items, item) => {
    const raw = item as RawGroceryItem;
    const id = raw.id ?? raw.ID;
    const name = raw.name ?? raw.Name;

    if (id == null || !name) {
      return items;
    }

    items.push({
      id,
      listId: raw.listId ?? raw.ListID,
      name,
      description: raw.description ?? raw.Description ?? null,
      createdAt: raw.createdAt ?? raw.CreatedAt,
      updatedAt: raw.updatedAt ?? raw.UpdatedAt,
    });

    return items;
  }, []);
}

function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);

  for (let i = 0; i < rawData.length; i += 1) {
    outputArray[i] = rawData.charCodeAt(i);
  }

  return outputArray;
}

export function App() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isAuthenticated, setIsAuthenticated] = useState<boolean | null>(null);
  const [isLoggingIn, setIsLoggingIn] = useState(false);
  const [items, setItems] = useState<GroceryItem[]>([]);
  const [newItemName, setNewItemName] = useState('');
  const [status, setStatus] = useState('Ready.');
  const [error, setError] = useState<string | null>(null);
  const [isLoadingItems, setIsLoadingItems] = useState(false);
  const [isAddingItem, setIsAddingItem] = useState(false);
  const [deletingItemId, setDeletingItemId] = useState<number | null>(null);
  const [isEnablingNotifications, setIsEnablingNotifications] = useState(false);
  const [notificationsEnabled, setNotificationsEnabled] = useState(false);

  const itemsUrl = `${API_URL}/households/${DEMO_HOUSEHOLD_ID}/lists/${DEMO_LIST_ID}/items`;

  async function loadItems() {
    setIsLoadingItems(true);
    setError(null);
    setStatus('Loading groceries...');

    try {
      const response = await fetch(itemsUrl, {
        credentials: 'include',
      });

      if (response.status === 401) {
        setItems([]);
        setIsAuthenticated(false);
        setStatus('Sign in to continue.');
        return;
      }

      if (!response.ok) {
        throw new Error(`List request failed with ${response.status}`);
      }

      const data = await response.json();
      const nextItems = normalizeItems(data);
      setItems(nextItems);
      setIsAuthenticated(true);
      setStatus(`${nextItems.length} item${nextItems.length === 1 ? '' : 's'} loaded.`);
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to load groceries.';
      setItems([]);
      setError(message);
      setStatus(message);
    } finally {
      setIsLoadingItems(false);
    }
  }

  async function login(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    setIsLoggingIn(true);
    setError(null);
    setStatus('Signing in...');

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email: email.trim(), password }),
      });

      if (response.status === 401) {
        throw new Error('Invalid email or password.');
      }
      if (!response.ok) {
        throw new Error(`Login failed with ${response.status}`);
      }

      setPassword('');
      setIsAuthenticated(true);
      await loadItems();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to sign in.';
      setIsAuthenticated(false);
      setError(message);
      setStatus(message);
    } finally {
      setIsLoggingIn(false);
    }
  }

  async function createItem(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const name = newItemName.trim();
    if (!name) {
      return;
    }

    setIsAddingItem(true);
    setError(null);
    setStatus('Adding item...');

    try {
      const response = await fetch(itemsUrl, {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ name }),
      });

      if (!response.ok) {
        throw new Error(`Add item failed with ${response.status}`);
      }

      setNewItemName('');
      await loadItems();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to add item.';
      setError(message);
      setStatus(message);
    } finally {
      setIsAddingItem(false);
    }
  }

  async function deleteItem(item: GroceryItem) {
    setDeletingItemId(item.id);
    setError(null);
    setStatus(`Deleting ${item.name}...`);

    try {
      const response = await fetch(`${itemsUrl}/${item.id}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error(`Delete item failed with ${response.status}`);
      }

      await loadItems();
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to delete item.';
      setError(message);
      setStatus(message);
    } finally {
      setDeletingItemId(null);
    }
  }

  async function getVapidPublicKey() {
    const response = await fetch(`${API_URL}/vapid-public-key`);
    if (!response.ok) {
      throw new Error(`VAPID key request failed with ${response.status}`);
    }

    const data = (await response.json()) as { publicKey?: string };
    if (!data.publicKey) {
      throw new Error('Backend response is missing publicKey.');
    }

    return data.publicKey;
  }

  async function enableNotifications() {
    if (!('serviceWorker' in navigator)) {
      throw new Error('This browser does not support service workers.');
    }

    if (!('PushManager' in window)) {
      throw new Error('This browser does not support Web Push.');
    }

    const permission = await Notification.requestPermission();
    if (permission !== 'granted') {
      throw new Error(`Notification permission is ${permission}.`);
    }

    const registration = await navigator.serviceWorker.register('/sw.js');
    const publicKey = await getVapidPublicKey();
    const existingSubscription = await registration.pushManager.getSubscription();

    if (existingSubscription) {
      await existingSubscription.unsubscribe();
    }

    const subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(publicKey),
    });

    const response = await fetch(
      `${API_URL}/households/${DEMO_HOUSEHOLD_ID}/notifications/subscribe`,
      {
        method: 'POST',
        credentials: 'include',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(subscription),
      },
    );

    if (!response.ok) {
      throw new Error(`Subscription save failed with ${response.status}`);
    }
  }

  async function handleEnableNotifications() {
    setIsEnablingNotifications(true);
    setError(null);
    setStatus('Enabling notifications...');

    try {
      await enableNotifications();
      setNotificationsEnabled(true);
      setStatus('Notifications enabled.');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to enable notifications.';
      setError(message);
      setStatus(message);
    } finally {
      setIsEnablingNotifications(false);
    }
  }

  useEffect(() => {
    void loadItems();
  }, [itemsUrl]);

  if (isAuthenticated === null) {
    return (
      <main className="app-shell auth-shell">
        <Loader color="green" />
      </main>
    );
  }

  if (!isAuthenticated) {
    return (
      <main className="app-shell">
        <Container size={420} px="md" py={{ base: 'xl', sm: 80 }}>
          <Stack gap="lg">
            <div>
              <Badge variant="light" color="green" radius="sm" mb="xs">
                Home
              </Badge>
              <Title order={1} className="page-title">
                Groceries
              </Title>
            </div>

            <Paper withBorder p="lg" radius="sm">
              <form onSubmit={login}>
                <Stack gap="md">
                  <TextInput
                    label="Email"
                    type="email"
                    value={email}
                    onChange={(event) => setEmail(event.currentTarget.value)}
                    leftSection={<Mail size={16} />}
                    required
                    autoComplete="email"
                  />
                  <PasswordInput
                    label="Password"
                    value={password}
                    onChange={(event) => setPassword(event.currentTarget.value)}
                    leftSection={<LockKeyhole size={16} />}
                    required
                    autoComplete="current-password"
                  />

                  {error && (
                    <Alert color="red" variant="light" radius="sm">
                      {error}
                    </Alert>
                  )}

                  <Button type="submit" leftSection={<LogIn size={16} />} loading={isLoggingIn}>
                    Sign in
                  </Button>
                </Stack>
              </form>
            </Paper>

            <Text role="status" c="dimmed" size="sm">
              {status}
            </Text>
          </Stack>
        </Container>
      </main>
    );
  }

  return (
    <main className="app-shell">
      <Container size="md" px="md" py={{ base: 'md', sm: 'xl' }}>
        <Stack gap="md">
          <Group justify="space-between" align="flex-start" gap="md" wrap="nowrap">
            <div>
              <Group gap="xs" mb={4}>
                <Badge variant="light" color="green" radius="sm">
                  Home
                </Badge>
                <Badge variant="outline" color="gray" radius="sm">
                  List #{DEMO_LIST_ID}
                </Badge>
              </Group>
              <Title order={1} className="page-title">
                Groceries
              </Title>
            </div>

            <Tooltip label="Refresh list">
              <ActionIcon
                aria-label="Refresh list"
                variant="light"
                color="green"
                size="lg"
                loading={isLoadingItems}
                onClick={loadItems}
              >
                <RefreshCcw size={18} />
              </ActionIcon>
            </Tooltip>
          </Group>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="sm">
              <Group justify="flex-end" align="end" gap="sm">
                <Button
                  variant={notificationsEnabled ? 'light' : 'outline'}
                  color={notificationsEnabled ? 'green' : 'gray'}
                  leftSection={notificationsEnabled ? <BellRing size={16} /> : <Bell size={16} />}
                  loading={isEnablingNotifications}
                  onClick={handleEnableNotifications}
                >
                  {notificationsEnabled ? 'Enabled' : 'Enable'}
                </Button>
              </Group>
            </Stack>
          </Paper>

          <Paper withBorder p="md" radius="sm">
            <Stack gap="md">
              <form onSubmit={createItem}>
                <Group gap="sm" align="end" wrap="nowrap">
                  <TextInput
                    className="new-item-input"
                    label="New item"
                    placeholder="Milk, eggs, detergent"
                    value={newItemName}
                    onChange={(event) => setNewItemName(event.currentTarget.value)}
                    leftSection={<ShoppingBasket size={16} />}
                  />
                  <Button type="submit" leftSection={<Plus size={16} />} loading={isAddingItem}>
                    Add
                  </Button>
                </Group>
              </form>

              <Divider />

              {error && (
                <Alert color="red" variant="light" radius="sm">
                  {error}
                </Alert>
              )}

              {isLoadingItems && items.length === 0 ? (
                <Group justify="center" py="xl">
                  <Loader color="green" />
                </Group>
              ) : items.length === 0 ? (
                <Paper className="empty-state" p="xl" radius="sm">
                  <ShoppingBasket size={28} />
                  <Text fw={650}>No items yet</Text>
                  <Text c="dimmed" size="sm">
                    Add the first grocery item above.
                  </Text>
                </Paper>
              ) : (
                <Stack gap="xs">
                  {items.map((item) => (
                    <Paper key={item.id} className="item-row" withBorder p="sm" radius="sm">
                      <Group justify="space-between" align="center" gap="sm" wrap="nowrap">
                        <Group gap="sm" wrap="nowrap" className="item-content">
                          <Check size={17} className="item-check" />
                          <div className="item-text">
                            <Text fw={700} className="item-name">
                              {item.name}
                            </Text>
                            {item.description && (
                              <Text c="dimmed" size="sm" className="item-description">
                                {item.description}
                              </Text>
                            )}
                          </div>
                        </Group>

                        <Tooltip label={`Delete ${item.name}`}>
                          <ActionIcon
                            aria-label={`Delete ${item.name}`}
                            variant="subtle"
                            color="red"
                            loading={deletingItemId === item.id}
                            onClick={() => deleteItem(item)}
                          >
                            <Trash2 size={18} />
                          </ActionIcon>
                        </Tooltip>
                      </Group>
                    </Paper>
                  ))}
                </Stack>
              )}
            </Stack>
          </Paper>

          <Text role="status" c="dimmed" size="sm">
            {status}
          </Text>
        </Stack>
      </Container>
    </main>
  );
}
