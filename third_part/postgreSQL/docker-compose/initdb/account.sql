CREATE TABLE accounts (
    id BIGSERIAL PRIMARY KEY,
    owner_id BIGINT NOT NULL,
    balance NUMERIC(18,2) NOT NULL DEFAULT 0, -- 可改为 double precision 视精度要求
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 按 owner 查询常用时建索引
CREATE INDEX idx_accounts_owner_id ON accounts(owner_id);

-- 自动维护 updated_at 的触发器
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;

CREATE TRIGGER trg_set_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();