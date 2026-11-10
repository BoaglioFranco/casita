import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useCurrentUser } from '../auth/queries';
import { createGroceryItem, deleteGroceryItem } from './api';
import { groceryKeys } from './queries';
import type { CreateGroceryItemInput } from './types';

export function useCreateGroceryItem(listId: string) {
  const { id: userId } = useCurrentUser();
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateGroceryItemInput) => createGroceryItem(listId, input),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(userId, listId) }),
  });
}

export function useDeleteGroceryItem(listId: string) {
  const { id: userId } = useCurrentUser();
  const client = useQueryClient();
  return useMutation({
    mutationFn: (itemId: string) => deleteGroceryItem(listId, itemId),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(userId, listId) }),
  });
}
