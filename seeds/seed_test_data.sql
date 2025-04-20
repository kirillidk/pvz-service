TRUNCATE products CASCADE;
TRUNCATE receptions CASCADE;
TRUNCATE pvz CASCADE;
TRUNCATE users CASCADE;

INSERT INTO users (id, email, password_hash, role) VALUES
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'employee1@example.com', '$2a$10$JQIPygpojfjEmm9l.TZ0G.l4r1EyrrCdJ4LfKuFtXFN2g3sUs8isK', 'employee'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'employee2@example.com', '$2a$10$JQIPygpojfjEmm9l.TZ0G.l4r1EyrrCdJ4LfKuFtXFN2g3sUs8isK', 'employee'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'employee3@example.com', '$2a$10$JQIPygpojfjEmm9l.TZ0G.l4r1EyrrCdJ4LfKuFtXFN2g3sUs8isK', 'employee'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'moderator1@example.com', '$2a$10$JQIPygpojfjEmm9l.TZ0G.l4r1EyrrCdJ4LfKuFtXFN2g3sUs8isK', 'moderator'),
  ('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'moderator2@example.com', '$2a$10$JQIPygpojfjEmm9l.TZ0G.l4r1EyrrCdJ4LfKuFtXFN2g3sUs8isK', 'moderator');

INSERT INTO pvz (id, registration_date, city) VALUES
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2023-01-01 10:00:00', 'Москва'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '2023-01-02 11:00:00', 'Москва'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '2023-01-03 12:00:00', 'Санкт-Петербург'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '2023-01-04 13:00:00', 'Санкт-Петербург'),
  ('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', '2023-01-05 14:00:00', 'Казань');

INSERT INTO receptions (id, date_time, pvz_id, status) VALUES
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2023-01-10 10:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'in_progress'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '2023-01-11 11:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'close'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '2023-01-12 12:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', 'in_progress'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '2023-01-13 13:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', 'in_progress'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', '2023-01-14 14:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', 'close'),
  ('c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', '2023-01-15 15:00:00', 'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', 'in_progress');

INSERT INTO products (id, date_time, type, reception_id) VALUES
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', '2023-01-10 10:30:00', 'электроника', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12', '2023-01-10 10:45:00', 'одежда', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13', '2023-01-11 11:30:00', 'обувь', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14', '2023-01-11 11:45:00', 'электроника', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a12'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15', '2023-01-12 12:30:00', 'одежда', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16', '2023-01-13 13:30:00', 'обувь', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a14'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a17', '2023-01-14 14:30:00', 'электроника', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a15'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a18', '2023-01-15 15:30:00', 'одежда', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16'),
  ('d0eebc99-9c0b-4ef8-bb6d-6bb9bd380a19', '2023-01-15 15:45:00', 'обувь', 'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a16');