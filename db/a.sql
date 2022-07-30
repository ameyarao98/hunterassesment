CREATE TABLE IF NOT EXISTS "user"(username VARCHAR(100) PRIMARY KEY);

CREATE TABLE IF NOT EXISTS "resource"(resource_name VARCHAR(6) PRIMARY KEY);

INSERT INTO
    "resource" (resource_name)
VALUES
    ('iron'),
    ('copper'),
    ('gold') ON CONFLICT (resource_name) DO NOTHING;

CREATE TABLE IF NOT EXISTS "factory"(
    id SERIAL PRIMARY KEY,
    resource_name VARCHAR(6) NOT NULL REFERENCES "resource" ON DELETE CASCADE,
    factory_level INTEGER NOT NULL,
    production_per_second INTEGER NOT NULL,
    next_upgrade_duration INTEGER NOT NULL,
    upgrade_cost JSON
);

INSERT INTO
    "factory" (
        id,
        resource_name,
        factory_level,
        production_per_second,
        next_upgrade_duration,
        upgrade_cost
    )
VALUES
    (
        1,
        'iron',
        1,
        10,
        15,
        '{ "iron": 300, "copper": 100, "gold": 1 }'
    ),
    (
        2,
        'iron',
        2,
        20,
        30,
        '{ "iron": 800, "copper": 250, "gold": 2 }'
    ),
    (
        3,
        'iron',
        3,
        40,
        60,
        '{ "iron": 1600, "copper": 500, "gold": 4 }'
    ),
    (
        4,
        'iron',
        4,
        80,
        90,
        '{ "iron": 3000, "copper": 1000, "gold": 8 }'
    ),
    (5, 'iron', 5, 150, 120, '{}'),
    (
        6,
        'copper',
        1,
        3,
        15,
        '{ "iron": 200, "copper": 70}'
    ),
    (
        7,
        'copper',
        2,
        7,
        30,
        '{ "iron": 400, "copper": 150}'
    ),
    (
        8,
        'copper',
        3,
        14,
        60,
        '{ "iron": 800, "copper": 300}'
    ),
    (
        9,
        'copper',
        4,
        30,
        90,
        '{ "iron": 1600, "copper": 600}'
    ),
    (10, 'copper', 5, 60, 120, '{}'),
    (
        11,
        'gold',
        1,
        2,
        15,
        '{ "copper": 100, "gold": 2}'
    ),
    (
        12,
        'gold',
        2,
        3,
        30,
        '{ "copper": 200, "gold": 4}'
    ),
    (
        13,
        'gold',
        3,
        4,
        60,
        '{ "copper": 400, "gold": 8}'
    ),
    (
        14,
        'gold',
        4,
        6,
        90,
        '{ "copper": 800, "gold": 16}'
    ),
    (15, 'gold', 5, 8, 120, '{}') ON CONFLICT (id) DO NOTHING;

CREATE TABLE IF NOT EXISTS "user_resource"(
    resource_name VARCHAR(6) NOT NULL REFERENCES "resource" ON DELETE CASCADE,
    username VARCHAR(100) NOT NULL REFERENCES "user" ON DELETE CASCADE,
    factory_level INTEGER DEFAULT 1 NOT NULL,
    amount INTEGER DEFAULT 0 NOT NULL,
    time_until_upgrade_complete INTEGER,
    PRIMARY KEY (resource_name, username)
);