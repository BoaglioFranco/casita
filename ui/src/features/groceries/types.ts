export type GroceryItem = {
  id: number;
  listId?: number;
  name: string;
  description: string | null;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateGroceryItemInput = { name: string; description?: string };
