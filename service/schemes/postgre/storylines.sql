CREATE TABLE storylines (
    id VARCHAR(64) PRIMARY KEY,
    account_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    cover_photo_uri TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN storylines.id IS '主键，故事线ID';
COMMENT ON COLUMN storylines.account_id IS '所属账号ID';
COMMENT ON COLUMN storylines.title IS '故事线标题';
COMMENT ON COLUMN storylines.description IS '故事线描述';
COMMENT ON COLUMN storylines.cover_photo_uri IS '封面图片URI';
COMMENT ON COLUMN storylines.created_at IS '创建时间';
COMMENT ON COLUMN storylines.updated_at IS '更新时间';

CREATE INDEX idx_storylines_account_id ON storylines(account_id);

CREATE TRIGGER trg_set_updated_at_storylines
BEFORE UPDATE ON storylines
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE storyline_nodes (
    id VARCHAR(64) PRIMARY KEY,
    storyline_id VARCHAR(64) NOT NULL REFERENCES storylines(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    date TIMESTAMPTZ NOT NULL,
    note TEXT NOT NULL DEFAULT '',
    location TEXT NOT NULL DEFAULT '',
    photo_uri TEXT NULL,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT storyline_nodes_storyline_sort_order_key UNIQUE (storyline_id, sort_order)
);

COMMENT ON COLUMN storyline_nodes.id IS '主键，节点ID';
COMMENT ON COLUMN storyline_nodes.storyline_id IS '所属故事线ID';
COMMENT ON COLUMN storyline_nodes.title IS '节点标题';
COMMENT ON COLUMN storyline_nodes.date IS '节点日期';
COMMENT ON COLUMN storyline_nodes.note IS '节点备注';
COMMENT ON COLUMN storyline_nodes.location IS '节点地点';
COMMENT ON COLUMN storyline_nodes.photo_uri IS '节点图片URI';
COMMENT ON COLUMN storyline_nodes.sort_order IS '节点排序';
COMMENT ON COLUMN storyline_nodes.created_at IS '创建时间';
COMMENT ON COLUMN storyline_nodes.updated_at IS '更新时间';

CREATE INDEX idx_storyline_nodes_storyline_id ON storyline_nodes(storyline_id);
CREATE INDEX idx_storyline_nodes_storyline_date ON storyline_nodes(storyline_id, date DESC);

CREATE TRIGGER trg_set_updated_at_storyline_nodes
BEFORE UPDATE ON storyline_nodes
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
