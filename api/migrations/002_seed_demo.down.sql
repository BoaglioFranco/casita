DELETE FROM lists
WHERE id = '019535d9-3df7-79fb-b466-fa907fa17f9f';

-- Household ownership is restrictive: explicitly remove accounts before
-- removing the demo household during a destructive rollback.
DELETE FROM users
WHERE household_id = '019535d9-3df7-79fb-b466-fa907fa17f9e';

DELETE FROM households
WHERE id = '019535d9-3df7-79fb-b466-fa907fa17f9e';
