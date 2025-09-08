CREATE Table peoples (
    id VARCHAR(64) PRIMARY KEY,
    name TEXT NOT NULL,
    gender VARCHAR(10) NOT NULL,
    birthdate DATE NOT NULL,
    address TEXT NOT NULL,
    phone BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN peoples.id IS '主键，身份证件号';
COMMENT ON COLUMN peoples.name IS '姓名';
COMMENT ON COLUMN peoples.gender IS '性别';
COMMENT ON COLUMN peoples.birthdate IS '出生日期';
COMMENT ON COLUMN peoples.address IS '地址';
COMMENT ON COLUMN peoples.phone IS '电话, 没有则填0';
COMMENT ON COLUMN peoples.created_at IS '创建时间';
COMMENT ON COLUMN peoples.updated_at IS '更新时间';

CREATE TRIGGER trg_set_updated_at_peoples
BEFORE UPDATE ON peoples
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();