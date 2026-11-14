import assert from 'node:assert/strict';
import { after, before, test } from 'node:test';
import { build } from 'esbuild';
import { mkdtemp, writeFile, rm } from 'node:fs/promises';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import { createRequire } from 'node:module';
import { fileURLToPath } from 'node:url';

let bundleDirectory;
let auth, lists, groceries, calendar, notifications, listQueries, groceryQueries, calendarQueries, notificationQueries, authGuard, QueryClient;
const originalFetch = globalThis.fetch;
const originalNavigator = Object.getOwnPropertyDescriptor(globalThis, 'navigator');

before(async () => {
  // Bundle the actual TypeScript modules; no HTTP server or browser is required.
  const modules = {
    auth: 'auth/api.ts', lists: 'lists/api.ts', groceries: 'groceries/api.ts',
    calendar: 'calendar/api.ts', notifications: 'notifications/api.ts',
    listQueries: 'lists/queries.ts', groceryQueries: 'groceries/queries.ts',
    calendarQueries: 'calendar/queries.ts', notificationQueries: 'notifications/queries.ts',
  };
  const result = await build({
    stdin: {
      contents: Object.entries(modules).map(([name, path]) => `export * as ${name} from './src/features/${path}';`).join('\n') + "\nexport * as authGuard from './src/routes/auth-guard.ts';\nexport { QueryClient } from '@tanstack/react-query';",
      resolveDir: fileURLToPath(new URL('..', import.meta.url)),
    },
    bundle: true, platform: 'node', format: 'cjs', write: false,
  });
  bundleDirectory = await mkdtemp(join(tmpdir(), 'casita-identity-tests-'));
  const bundlePath = join(bundleDirectory, 'modules.cjs');
  await writeFile(bundlePath, result.outputFiles[0].contents);
  ({ auth, lists, groceries, calendar, notifications, listQueries, groceryQueries, calendarQueries, notificationQueries, authGuard, QueryClient } = createRequire(import.meta.url)(bundlePath));
});

after(async () => {
  globalThis.fetch = originalFetch;
  if (originalNavigator) Object.defineProperty(globalThis, 'navigator', originalNavigator);
  else delete globalThis.navigator;
  delete globalThis.window;
  delete globalThis.Notification;
  delete globalThis.PushManager;
  if (bundleDirectory) await rm(bundleDirectory, { recursive: true, force: true });
});

const json = (body) => new Response(JSON.stringify(body), { headers: { 'Content-Type': 'application/json' } });

test('resource API calls use implicit household routes and session cookies', async () => {
  const requests = [];
  globalThis.fetch = async (path, init) => {
    requests.push({ path, ...init });
    if (path === '/api/me') return json({ id: 'user-a', name: 'Alice', householdId: 'home-a' });
    if (path === '/api/lists') return json({ lists: [] });
    if (path === '/api/lists/list-a/items' && !init.method) return json({ items: [] });
    return json({});
  };
  assert.deepEqual(await auth.getCurrentUser(), { id: 'user-a', name: 'Alice', householdId: 'home-a' });
  await lists.getLists();
  await groceries.getGroceryItems('list-a');
  await groceries.createGroceryItem('list-a', { name: 'Milk' });
  await groceries.deleteGroceryItem('list-a', 'item-a');
  await calendar.getCalendar({ from: '2026-11-10', to: '2026-11-12' });
  await calendar.createEvent({ name: 'Lunch', startsAt: '2026-11-11T12:00:00Z', durationMinutes: 60 });
  assert.deepEqual(requests.map(({ path }) => path), [
    '/api/me', '/api/lists', '/api/lists/list-a/items', '/api/lists/list-a/items',
    '/api/lists/list-a/items/item-a', '/api/calendar?from=2026-11-10&to=2026-11-12', '/api/calendar/events',
  ]);
  assert.ok(requests.every(({ credentials }) => credentials === 'include'));
  assert.deepEqual(JSON.parse(requests[3].body), { name: 'Milk' });
});

test('resource and notification cache keys are partitioned by user identity', () => {
  const window = { from: '2026-11-10', to: '2026-11-12' };
  for (const query of [
    (user) => listQueries.listsQuery(user),
    (user) => groceryQueries.groceryItemsQuery(user, 'same-list'),
    (user) => calendarQueries.calendarQuery(user, window),
    (user) => notificationQueries.notificationStateQuery(user),
  ]) {
    assert.notDeepEqual(query('user-a').queryKey, query('user-b').queryKey);
    assert.deepEqual(query('user-a').queryKey, query('user-a').queryKey);
  }
});

function allowPush() {
  const subscription = { endpoint: 'https://push.example.test/shared-device', keys: { p256dh: 'key', auth: 'auth' } };
  Object.defineProperty(globalThis, 'navigator', { configurable: true, value: {
    userAgent: 'Desktop', platform: 'Linux', maxTouchPoints: 0,
    serviceWorker: { getRegistration: async () => ({ pushManager: { getSubscription: async () => subscription } }) },
  } });
  globalThis.window = { matchMedia: () => ({ matches: false }), PushManager: class {}, Notification: {} };
  globalThis.Notification = { permission: 'granted' };
  notifications.resetSubscriptionSyncs();
  return subscription;
}

test('existing subscriptions resynchronize when switching A to B and back to A', async () => {
  const subscription = allowPush();
  let signedInUser = 'user-a';
  const registrations = [];
  globalThis.fetch = async (path, init) => {
    assert.equal(path, '/api/notifications/subscribe');
    assert.equal(init.method, 'POST');
    assert.deepEqual(JSON.parse(init.body), subscription);
    registrations.push(signedInUser);
    return json({});
  };
  assert.equal(await notifications.getNotificationState('user-a'), 'enabled');
  await notifications.getNotificationState('user-a');
  signedInUser = 'user-b';
  notifications.resetSubscriptionSyncs(); // Called after successful login.
  await notifications.getNotificationState('user-b');
  signedInUser = 'user-a';
  notifications.resetSubscriptionSyncs();
  await notifications.getNotificationState('user-a');
  assert.deepEqual(registrations, ['user-a', 'user-b', 'user-a']);
});

test('a failed subscription synchronization can be retried', async () => {
  allowPush();
  let attempts = 0;
  globalThis.fetch = async () => ++attempts === 1 ? new Response('', { status: 502 }) : json({});
  await assert.rejects(notifications.getNotificationState('user-a'));
  assert.equal(await notifications.getNotificationState('user-a'), 'enabled');
  assert.equal(attempts, 2);
});

test('authentication guard reads /me and redirects an invalid session to login', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, gcTime: 0 } } });
  const user = { id: 'user-a', name: 'Alice', householdId: 'home-a' };
  globalThis.fetch = async (path) => { assert.equal(path, '/api/me'); return json(user); };
  assert.deepEqual(await authGuard.requireAuthentication(client), { user });
  client.clear();
  globalThis.fetch = async () => new Response('', { status: 401 });
  await assert.rejects(authGuard.requireAuthentication(client), (error) => error.options.to === '/login');
  client.clear();
});

test('switching accounts cannot reuse the prior account resource cache', async () => {
  const client = new QueryClient({ defaultOptions: { queries: { retry: false, staleTime: Infinity, gcTime: 0 } } });
  let account = 'user-a';
  globalThis.fetch = async () => json({ lists: [{ id: `${account}-list`, name: account }] });
  assert.equal((await client.fetchQuery(listQueries.listsQuery('user-a')))[0].id, 'user-a-list');
  account = 'user-b';
  assert.equal((await client.fetchQuery(listQueries.listsQuery('user-b')))[0].id, 'user-b-list');
  assert.equal(client.getQueryData(listQueries.listsQuery('user-a').queryKey)[0].id, 'user-a-list');
  client.clear();
});
