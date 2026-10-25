import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createGroceryItem, deleteGroceryItem } from './api';
import { groceryKeys } from './queries';
import type { CreateGroceryItemInput } from './types';

export function useCreateGroceryItem(householdId: number, listId: number) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CreateGroceryItemInput) => createGroceryItem(householdId, listId, input),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(householdId, listId) }),
  });
}

export function useDeleteGroceryItem(householdId: number, listId: number) {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (itemId: number) => deleteGroceryItem(householdId, listId, itemId),
    onSuccess: () => client.invalidateQueries({ queryKey: groceryKeys.list(householdId, listId) }),
  });
}
