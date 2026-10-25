import { apiRequest } from '../../lib/api-client';
import type { CreateGroceryItemInput, GroceryItem } from './types';

type GroceryItemsResponse = { items?: unknown[] };
type RawGroceryItem = {
  id?: number;
  ID?: number;
  listId?: number;
  ListID?: number;
  name?: string;
  Name?: string;
  description?: string | null;
  Description?: string | null;
  createdAt?: string;
  CreatedAt?: string;
  updatedAt?: string;
  UpdatedAt?: string;
};

function itemsPath(householdId: number, listId: number): string {
  return `/households/${householdId}/lists/${listId}/items`;
}

export async function getGroceryItems(householdId: number, listId: number): Promise<GroceryItem[]> {
  const response = await apiRequest<GroceryItemsResponse>(itemsPath(householdId, listId));
  return normalizeItems(response.items);
}

export function createGroceryItem(
  householdId: number,
  listId: number,
  input: CreateGroceryItemInput,
): Promise<void> {
  return apiRequest(itemsPath(householdId, listId), {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export function deleteGroceryItem(
  householdId: number,
  listId: number,
  itemId: number,
): Promise<void> {
  return apiRequest(`${itemsPath(householdId, listId)}/${itemId}`, { method: 'DELETE' });
}

function normalizeItems(data: unknown): GroceryItem[] {
  if (!Array.isArray(data)) return [];

  return data.reduce<GroceryItem[]>((items, value) => {
    const item = value as RawGroceryItem;
    const id = item.id ?? item.ID;
    const name = item.name ?? item.Name;
    if (typeof id !== 'number' || typeof name !== 'string' || name.length === 0) return items;

    items.push({
      id,
      name,
      listId: item.listId ?? item.ListID,
      description: item.description ?? item.Description ?? null,
      createdAt: item.createdAt ?? item.CreatedAt,
      updatedAt: item.updatedAt ?? item.UpdatedAt,
    });
    return items;
  }, []);
}
