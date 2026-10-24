INSERT INTO households (id, name)
VALUES (1, 'Home')
ON CONFLICT (id) DO NOTHING;

INSERT INTO lists (id, household_id, name)
VALUES (1, 1, 'Groceries')
ON CONFLICT (id) DO NOTHING;

SELECT setval(
    pg_get_serial_sequence('households', 'id'),
    (SELECT MAX(id) FROM households)
);

SELECT setval(
    pg_get_serial_sequence('lists', 'id'),
    (SELECT MAX(id) FROM lists)
);
