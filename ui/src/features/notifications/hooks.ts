import { useMutation, useQueryClient } from '@tanstack/react-query';
import { enableNotifications } from './api';
import { notificationKeys } from './queries';

export function useEnableNotifications(householdId: number) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: () => enableNotifications(householdId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: notificationKeys.all }),
  });
}
