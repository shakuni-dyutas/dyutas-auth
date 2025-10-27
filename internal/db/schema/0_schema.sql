-- Basic PG schema initialization
-- named 0_schema.sql to be executed first.
-- official PG image when first started, it all executes 
-- *.sql and *.sh files in specific directory.
-- and the exact order follows en_US.utf8 locale.
CREATE SCHEMA IF NOT EXISTS auth;