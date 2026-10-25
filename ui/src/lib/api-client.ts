import { appConfig } from '../app/config';

type ApiErrorBody = { error?: string };

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

export async function apiRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  const headers = new Headers(init.headers);
  if (init.body && !headers.has('Content-Type')) headers.set('Content-Type', 'application/json');

  const response = await fetch(`${appConfig.apiUrl}${path}`, {
    ...init,
    headers,
    credentials: 'include',
  });

  if (!response.ok) {
    const body = await readJson<ApiErrorBody>(response);
    const fallback = response.statusText || `Request failed with status ${response.status}`;
    throw new ApiError(response.status, body?.error || fallback);
  }

  if (response.status === 204) return undefined as T;
  return (await readJson<T>(response)) as T;
}

async function readJson<T>(response: Response): Promise<T | undefined> {
  if (!response.headers.get('content-type')?.includes('application/json')) return undefined;
  return response.json() as Promise<T>;
}

export function getErrorMessage(error: unknown, fallback: string): string {
  return error instanceof Error ? error.message : fallback;
}
