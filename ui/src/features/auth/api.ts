import { apiRequest } from '../../lib/api-client';

export type LoginCredentials = { email: string; password: string };

export function login(credentials: LoginCredentials): Promise<void> {
  return apiRequest('/login', {
    method: 'POST',
    body: JSON.stringify({ email: credentials.email.trim(), password: credentials.password }),
  });
}
