-- Create Forge DB and user
CREATE DATABASE forge WITH ENCODING 'UTF8';
-- Password should be changed before use in environments deployed in Cloud
CREATE USER forge WITH PASSWORD 'forge';
GRANT ALL PRIVILEGES ON DATABASE forge TO forge;
