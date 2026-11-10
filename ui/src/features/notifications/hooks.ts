import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useCurrentUser } from '../auth/queries';
import { enableNotifications } from './api';
import { notificationKeys } from './queries';

export function useEnableNotifications() {
  const { id: userId } = useCurrentUser();
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => enableNotifications(userId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });
}
