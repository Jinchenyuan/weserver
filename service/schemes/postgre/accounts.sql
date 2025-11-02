CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    account VARCHAR(64) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    phone BIGINT NOT NULL,
    position TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN accounts.id IS '主键，账号ID';
COMMENT ON COLUMN accounts.account IS '账号，唯一';
COMMENT ON COLUMN accounts.password IS '密码，存储base64加密后的密码';
COMMENT ON COLUMN accounts.name IS '姓名';
COMMENT ON COLUMN accounts.email IS '邮箱';
COMMENT ON COLUMN accounts.phone IS '电话';
COMMENT ON COLUMN accounts.position IS '职位';
COMMENT ON COLUMN accounts.created_at IS '创建时间';
COMMENT ON COLUMN accounts.updated_at IS '更新时间';

CREATE TRIGGER trg_set_updated_at_accounts
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();