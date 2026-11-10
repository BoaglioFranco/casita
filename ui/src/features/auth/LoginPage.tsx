import { FormEvent, useState } from 'react';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useNavigate } from '@tanstack/react-router';
import {
  Alert,
  Button,
  Center,
  Container,
  Paper,
  PasswordInput,
  Stack,
  TextInput,
  Title,
} from '@mantine/core';
import { LockKeyhole, LogIn, Mail } from 'lucide-react';
import { login } from './api';
import { resetSubscriptionSyncs } from '../notifications/api';
import { getErrorMessage } from '../../lib/api-client';

export function LoginPage() {
  const navigate = useNavigate();
  const client = useQueryClient();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const loginMutation = useMutation({
    mutationFn: login,
    onSuccess: async () => {
      setPassword('');
      client.clear();
      resetSubscriptionSyncs();
      await navigate({ to: '/' });
    },
  });

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    loginMutation.mutate({ email, password });
  }

  return (
    <Center component="main" flex={1} w="100%">
      <Container size={420} w="100%" px="md" py={{ base: 'xl', sm: 80 }}>
        <Stack gap="lg">
          <Title order={1} fz={{ base: 32, sm: 48 }} lh={1} lts={0}>
            Sign in
          </Title>

          <Paper withBorder p="lg" radius="sm">
            <form onSubmit={handleSubmit}>
              <Stack gap="md">
                <TextInput
                  label="Email"
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.currentTarget.value)}
                  leftSection={<Mail size={16} />}
                  required
                  autoComplete="email"
                  autoFocus
                />
                <PasswordInput
                  label="Password"
                  value={password}
                  onChange={(event) => setPassword(event.currentTarget.value)}
                  leftSection={<LockKeyhole size={16} />}
                  required
                  autoComplete="current-password"
                />

                {loginMutation.isError && (
                  <Alert color="red" variant="light" radius="sm">
                    {getErrorMessage(loginMutation.error, 'Unable to sign in.')}
                  </Alert>
                )}

                <Button
                  type="submit"
                  leftSection={<LogIn size={16} />}
                  loading={loginMutation.isPending}
                >
                  Sign in
                </Button>
              </Stack>
            </form>
          </Paper>
        </Stack>
      </Container>
    </Center>
  );
}
