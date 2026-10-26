import { apiRequest } from '../../lib/api-client';
import type { CreateGroceryItemInput, GroceryItem, ItemCreator } from './types';

type GroceryItemsResponse = { items?: unknown[] };
type RawGroceryItem = {
  id?: string;
  ID?: string;
  listId?: string;
  ListID?: string;
  name?: string;
  Name?: string;
  description?: string | null;
  Description?: string | null;
  createdBy?: unknown;
  CreatedBy?: unknown;
  createdAt?: string;
  CreatedAt?: string;
  updatedAt?: string;
  UpdatedAt?: string;
};

function itemsPath(householdId: string, listId: string): string {
  return `/households/${householdId}/lists/${listId}/items`;
}

export async function getGroceryItems(householdId: string, listId: string): Promise<GroceryItem[]> {
  const response = await apiRequest<GroceryItemsResponse>(itemsPath(householdId, listId));
  return normalizeItems(response.items);
}

export function createGroceryItem(
  householdId: string,
  listId: string,
  input: CreateGroceryItemInput,
): Promise<void> {
  return apiRequest(itemsPath(householdId, listId), {
    method: 'POST',
    body: JSON.stringify(input),
  });
}

export function deleteGroceryItem(
  householdId: string,
  listId: string,
  itemId: string,
): Promise<void> {
  return apiRequest(`${itemsPath(householdId, listId)}/${itemId}`, { method: 'DELETE' });
}

function normalizeItems(data: unknown): GroceryItem[] {
  if (!Array.isArray(data)) return [];

  return data.reduce<GroceryItem[]>((items, value) => {
    const item = value as RawGroceryItem;
    const id = item.id ?? item.ID;
    const name = item.name ?? item.Name;
    if (typeof id !== 'string' || typeof name !== 'string' || name.length === 0) return items;

    items.push({
      id,
      name,
      listId: item.listId ?? item.ListID,
      description: item.description ?? item.Description ?? null,
      createdBy: normalizeCreator(item.createdBy ?? item.CreatedBy),
      createdAt: item.createdAt ?? item.CreatedAt,
      updatedAt: item.updatedAt ?? item.UpdatedAt,
    });
    return items;
  }, []);
}

function normalizeCreator(value: unknown): ItemCreator | null {
  if (value == null || typeof value !== 'object') return null;

  const creator = value as { id?: unknown; name?: unknown };
  if (typeof creator.id !== 'string' || typeof creator.name !== 'string') return null;

  return { id: creator.id, name: creator.name };
}
