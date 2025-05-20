-- Insert sample UserRelationship data with current timestamps
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('mandy@example.com',   'trendy@example.com',  'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('trendy@example.com',  'mandy@example.com',  'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('trendy@example.com',  'alameda@example.com', 'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('alameda@example.com', 'trendy@example.com',  'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('alameda@example.com', 'bingo@example.com',   'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('bingo@example.com',   'alameda@example.com', 'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('bingo@example.com',   'trendy@example.com',  'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('trendy@example.com',  'bingo@example.com',   'FRIEND', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('leo@example.com',     'trendy@example.com',  'BLOCK',  CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('adison@example.com',  'trendy@example.com',  'SUBSCRIBER', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
INSERT INTO user_relationships
    (requestor_email, target_email, type, created_at, updated_at)
VALUES
    ('lucas@example.com',   'trendy@example.com',  'SUBSCRIBER', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
