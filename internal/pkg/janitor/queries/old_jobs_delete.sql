
DELETE FROM jobs
WHERE updated_at < extract(epoch FROM now() - ($1 || ' days')::interval)::int;
