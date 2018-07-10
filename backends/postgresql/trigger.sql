-- Provides a function for updating lastmodified timestamps that can be used in triggers
-- https://www.revsys.com/tidbits/automatically-updating-a-timestamp-column-in-postgresql/
CREATE OR REPLACE FUNCTION update_lastupdated_column() 
  RETURNS trigger 
AS
$BODY$
DECLARE
    depatureDate DATE;
BEGIN
    NEW.lastupdated = now();
    RETURN NEW;
END;
$BODY$
LANGUAGE plpgsql;