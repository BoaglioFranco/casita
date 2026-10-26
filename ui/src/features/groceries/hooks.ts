import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createGroceryItem, deleteGroceryItem } from './api';
import { groceryKeys } from './queries';
import type { CreateGroceryItemInput } from './types';

export function useCreateGroceryItem(householdId: string, listId: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateGroceryItemInput) => createGroceryItem(householdId, listId, input),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(householdId, listId) }),
  });
}

export function useDeleteGroceryItem(householdId: string, listId: string) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (itemId: string) => deleteGroceryItem(householdId, listId, itemId),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(householdId, listId) }),
  });
}
