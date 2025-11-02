CREATE TABLE s3kv (
    key VARCHAR(64) NOT NULL UNIQUE,
    value TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN s3kv.key IS '键，唯一';
COMMENT ON COLUMN s3kv.value IS '值';
COMMENT ON COLUMN s3kv.created_at IS '创建时间';
COMMENT ON COLUMN s3kv.updated_at IS '更新时间';

CREATE TRIGGER trg_set_updated_at_s3kv
BEFORE UPDATE ON s3kv
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();