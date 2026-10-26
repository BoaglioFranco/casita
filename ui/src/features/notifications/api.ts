import { apiRequest } from '../../lib/api-client';

type VapidKeyResponse = { publicKey?: string };

export type NotificationState =
  | 'disabled'
  | 'enabled'
  | 'denied'
  | 'install-required'
  | 'unsupported';

export function isStandaloneApp(): boolean {
  const navigatorWithStandalone = navigator as Navigator & { standalone?: boolean };

  return (
    window.matchMedia('(display-mode: standalone)').matches ||
    navigatorWithStandalone.standalone === true
  );
}

export function isMobileDevice(): boolean {
  return (
    /Android|iPad|iPhone|iPod/.test(navigator.userAgent) ||
    (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
  );
}

function isIosDevice(): boolean {
  return (
    /iPad|iPhone|iPod/.test(navigator.userAgent) ||
    (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
  );
}

export async function getNotificationState(): Promise<NotificationState> {
  if (isIosDevice() && !isStandaloneApp()) return 'install-required';

  if (
    !('serviceWorker' in navigator) ||
    !('PushManager' in window) ||
    !('Notification' in window)
  ) {
    return 'unsupported';
  }

  if (Notification.permission === 'denied') return 'denied';
  if (Notification.permission !== 'granted') return 'disabled';

  const registration = await navigator.serviceWorker.getRegistration();
  return (await registration?.pushManager.getSubscription()) == null ? 'disabled' : 'enabled';
}

export async function enableNotifications(householdId: string): Promise<void> {
  if (!('serviceWorker' in navigator)) {
    throw new Error('This browser does not support service workers.');
  }
  if (!('PushManager' in window) || !('Notification' in window)) {
    throw new Error('This browser does not support Web Push.');
  }

  const permission = await Notification.requestPermission();
  if (permission !== 'granted') {
    throw new Error(`Notification permission is ${permission}.`);
  }

  const registration = await navigator.serviceWorker.register('/sw.js');
  const { publicKey } = await apiRequest<VapidKeyResponse>('/vapid-public-key');
  if (!publicKey) throw new Error('Backend response is missing publicKey.');

  const existingSubscription = await registration.pushManager.getSubscription();
  const subscription =
    existingSubscription ??
    (await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(publicKey),
    }));

  await apiRequest(`/households/${householdId}/notifications/subscribe`, {
    method: 'POST',
    body: JSON.stringify(subscription),
  });
}

function urlBase64ToUint8Array(base64String: string): Uint8Array<ArrayBuffer> {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
  const rawData = window.atob(base64);
  return Uint8Array.from(rawData, (character) => character.charCodeAt(0));
}
