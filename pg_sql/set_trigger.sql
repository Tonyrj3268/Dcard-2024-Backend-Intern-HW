CREATE OR REPLACE FUNCTION update_active_trigger()
RETURNS TRIGGER AS $$
BEGIN
  IF NEW.end_at < NOW() THEN
    NEW.active = FALSE;
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER advertisement_update_active
BEFORE INSERT OR UPDATE ON Advertisements
FOR EACH ROW
EXECUTE FUNCTION update_active_trigger();
