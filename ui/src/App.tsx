import { FormEvent, useEffect, useMemo, useState } from 'react';
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
  Plus,
  RefreshCcw,
  Send,
  ShoppingBasket,
  Trash2,
  Wifi,
} from 'lucide-react';

const DEMO_HOUSEHOLD_ID = 1;
const DEMO_LIST_ID = 1;
const DEMO_MEMBER_ID = 1;
const API_URL_STORAGE_KEY = 'groceries-demo-api-url';

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

function authHeaders(extraHeaders: HeadersInit = {}): HeadersInit {
  return {
    ...extraHeaders,
    Authorization: `Bearer ${DEMO_MEMBER_ID}`,
  };
}

function normalizeApiUrl(value: string) {
  return value.trim().replace(/\/$/, '');
}

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
  const [apiUrl, setApiUrl] = useState(() => {
    return localStorage.getItem(API_URL_STORAGE_KEY) ?? 'http://localhost:2345';
  });
  const [items, setItems] = useState<GroceryItem[]>([]);
  const [newItemName, setNewItemName] = useState('');
  const [status, setStatus] = useState('Ready.');
  const [error, setError] = useState<string | null>(null);
  const [isLoadingItems, setIsLoadingItems] = useState(false);
  const [isAddingItem, setIsAddingItem] = useState(false);
  const [deletingItemId, setDeletingItemId] = useState<number | null>(null);
  const [isEnablingNotifications, setIsEnablingNotifications] = useState(false);
  const [isSendingTestPush, setIsSendingTestPush] = useState(false);
  const [notificationsEnabled, setNotificationsEnabled] = useState(false);

  const normalizedApiUrl = useMemo(() => normalizeApiUrl(apiUrl), [apiUrl]);
  const itemsUrl = `${normalizedApiUrl}/households/${DEMO_HOUSEHOLD_ID}/lists/${DEMO_LIST_ID}/items`;

  async function loadItems() {
    setIsLoadingItems(true);
    setError(null);
    setStatus('Loading groceries...');

    try {
      const response = await fetch(itemsUrl, {
        headers: authHeaders(),
      });

      if (!response.ok) {
        throw new Error(`List request failed with ${response.status}`);
      }

      const data = await response.json();
      const nextItems = normalizeItems(data);
      setItems(nextItems);
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
        headers: authHeaders({
          'Content-Type': 'application/json',
        }),
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
        headers: authHeaders(),
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
    const response = await fetch(`${normalizedApiUrl}/vapid-public-key`);
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

    const response = await fetch(`${normalizedApiUrl}/subscribe`, {
      method: 'POST',
      headers: authHeaders({
        'Content-Type': 'application/json',
      }),
      body: JSON.stringify(subscription),
    });

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

  async function sendTestPush() {
    setIsSendingTestPush(true);
    setError(null);
    setStatus('Sending test push...');

    try {
      const response = await fetch(`${normalizedApiUrl}/send-test`, {
        method: 'POST',
      });

      if (!response.ok) {
        throw new Error(`Test push failed with ${response.status}`);
      }

      setStatus('Test push requested.');
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Unable to send test push.';
      setError(message);
      setStatus(message);
    } finally {
      setIsSendingTestPush(false);
    }
  }

  useEffect(() => {
    localStorage.setItem(API_URL_STORAGE_KEY, normalizedApiUrl);
  }, [normalizedApiUrl]);

  useEffect(() => {
    void loadItems();
  }, [itemsUrl]);

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
              <Group justify="space-between" align="end" gap="sm">
                <TextInput
                  className="api-input"
                  label="API"
                  value={apiUrl}
                  onChange={(event) => setApiUrl(event.currentTarget.value)}
                  leftSection={<Wifi size={16} />}
                />

                <Group gap="xs" wrap="nowrap">
                  <Button
                    variant={notificationsEnabled ? 'light' : 'outline'}
                    color={notificationsEnabled ? 'green' : 'gray'}
                    leftSection={notificationsEnabled ? <BellRing size={16} /> : <Bell size={16} />}
                    loading={isEnablingNotifications}
                    onClick={handleEnableNotifications}
                  >
                    {notificationsEnabled ? 'Enabled' : 'Enable'}
                  </Button>

                  <Tooltip label="Send test push">
                    <ActionIcon
                      aria-label="Send test push"
                      variant="outline"
                      color="gray"
                      size="lg"
                      loading={isSendingTestPush}
                      disabled={!notificationsEnabled}
                      onClick={sendTestPush}
                    >
                      <Send size={18} />
                    </ActionIcon>
                  </Tooltip>
                </Group>
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
