CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    body JSONB NOT NULL
);

CREATE OR REPLACE FUNCTION notify_event()
RETURNS TRIGGER LANGUAGE plpgsql AS
$$
BEGIN
    PERFORM pg_notify('events', row_to_json(NEW)::text);
    RETURN null;
END;
$$;


CREATE OR REPLACE TRIGGER notify_on_new_event
    AFTER INSERT ON events
    FOR EACH ROW
    EXECUTE PROCEDURE notify_event();
