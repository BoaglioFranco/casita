INSERT INTO households (id, name)
VALUES ('019535d9-3df7-79fb-b466-fa907fa17f9e', 'Casita')
ON CONFLICT (id) DO NOTHING;

INSERT INTO lists (id, household_id, name)
VALUES (
    '019535d9-3df7-79fb-b466-fa907fa17f9f',
    '019535d9-3df7-79fb-b466-fa907fa17f9e',
    'Compras'
)
ON CONFLICT (id) DO NOTHING;
