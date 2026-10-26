export type ItemCreator = {
  id: string;
  name: string;
};

export type GroceryItem = {
  id: string;
  listId?: string;
  name: string;
  description: string | null;
  createdBy: ItemCreator | null;
  createdAt?: string;
  updatedAt?: string;
};

export type CreateGroceryItemInput = { name: string; description?: string };
