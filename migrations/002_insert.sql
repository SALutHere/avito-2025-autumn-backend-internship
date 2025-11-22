INSERT INTO teams (name) VALUES
('backend'),
('frontend'),
('devops'),
('qa'),
('analytics')
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, team_name, is_active) VALUES
('u1001', 'alice',     'backend',  TRUE),
('u1002', 'bob',       'backend',  TRUE),
('u1003', 'carol',     'backend',  TRUE),
('u1004', 'dave',      'backend',  TRUE),
('u1005', 'erin',      'backend',  FALSE),
('u1006', 'frank',     'backend',  TRUE),
('u1007', 'grace',     'backend',  TRUE),
('u1008', 'heidi',     'backend',  TRUE),
('u1009', 'ivan',      'backend',  FALSE),
('u1010', 'judy',      'backend',  TRUE),
('u1011', 'mallory',   'backend',  TRUE),
('u1012', 'oscar',     'backend',  TRUE),
('u1013', 'peggy',     'backend',  TRUE),
('u1014', 'trent',     'backend',  FALSE),
('u1015', 'victor',    'backend',  TRUE)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, team_name, is_active) VALUES
('u2001', 'amy',        'frontend', TRUE),
('u2002', 'ben',        'frontend', TRUE),
('u2003', 'chloe',      'frontend', TRUE),
('u2004', 'daniel',     'frontend', TRUE),
('u2005', 'ella',       'frontend', FALSE),
('u2006', 'felix',      'frontend', TRUE),
('u2007', 'george',     'frontend', TRUE),
('u2008', 'hannah',     'frontend', TRUE),
('u2009', 'isabel',     'frontend', TRUE),
('u2010', 'jack',       'frontend', FALSE),
('u2011', 'kevin',      'frontend', TRUE),
('u2012', 'lucy',       'frontend', TRUE)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, team_name, is_active) VALUES
('u3001', 'alex',       'devops', TRUE),
('u3002', 'brian',      'devops', TRUE),
('u3003', 'charlie',    'devops', TRUE),
('u3004', 'dennis',     'devops', TRUE),
('u3005', 'edward',     'devops', TRUE),
('u3006', 'fabian',     'devops', FALSE),
('u3007', 'george_d',   'devops', TRUE),
('u3008', 'harold',     'devops', TRUE),
('u3009', 'ivan_d',     'devops', TRUE),
('u3010', 'josh',       'devops', TRUE)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, team_name, is_active) VALUES
('u4001', 'adam',      'qa', TRUE),
('u4002', 'bella',     'qa', TRUE),
('u4003', 'claire',    'qa', TRUE),
('u4004', 'daria',     'qa', TRUE),
('u4005', 'elena',     'qa', FALSE),
('u4006', 'fiona',     'qa', TRUE),
('u4007', 'gianna',    'qa', TRUE),
('u4008', 'harper',    'qa', TRUE)
ON CONFLICT DO NOTHING;

INSERT INTO users (id, username, team_name, is_active) VALUES
('u5001', 'alexis',      'analytics', TRUE),
('u5002', 'blake',       'analytics', TRUE),
('u5003', 'casey',       'analytics', TRUE),
('u5004', 'dylan',       'analytics', FALSE),
('u5005', 'eva',         'analytics', TRUE),
('u5006', 'finley',      'analytics', TRUE),
('u5007', 'greta',       'analytics', TRUE),
('u5008', 'harry',       'analytics', TRUE),
('u5009', 'isabella',    'analytics', FALSE),
('u5010', 'jason',       'analytics', TRUE)
ON CONFLICT DO NOTHING;
