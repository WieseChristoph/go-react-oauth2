CREATE EXTENSION citext;

CREATE TABLE user_data (
  id SERIAL PRIMARY KEY,
  name TEXT NOT NULL,
  display_name TEXT,
  email CITEXT UNIQUE,
  avatar TEXT,
  role CITEXT NOT NULL,

  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE TABLE account (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  provider TEXT NOT NULL,
  provider_account_id TEXT NOT NULL,
  access_token TEXT,
  refresh_token TEXT,
  expires_at TIMESTAMP(0) WITH TIME ZONE,
  scope TEXT,

  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

  FOREIGN KEY (user_id) REFERENCES user_data(id)
);

CREATE TABLE session (
  id SERIAL PRIMARY KEY,
  token TEXT UNIQUE NOT NULL,
  user_id INTEGER NOT NULL,
  ip_address INET,
  expires_at TIMESTAMP(0) WITH TIME ZONE NOT NULL,

  created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),

  FOREIGN KEY (user_id) REFERENCES user_data(id)
);

-- Create a function to set the updated_at column to the current timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a function to create a trigger for each table that has a updated_at column to set the this column to the current timestamp on update
CREATE OR REPLACE FUNCTION create_trigger_for_tables()
RETURNS VOID AS $$
DECLARE
  t TEXT;
BEGIN
  FOR t IN 
    SELECT table_name 
    FROM information_schema.columns
    WHERE (
      column_name = 'updated_at' 
      AND (
        SELECT 1 
        FROM information_schema.triggers
        WHERE trigger_name = 'set_timestamp'
        AND event_object_table = table_name
      ) IS NULL
    ) 
  LOOP
    EXECUTE format('
      CREATE TRIGGER set_timestamp
      BEFORE UPDATE ON %I
      FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();
    ', t);
  END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Call the function to create the triggers
SELECT create_trigger_for_tables();