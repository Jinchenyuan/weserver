CREATE TABLE event_records (
    id VARCHAR(64) PRIMARY KEY,
    account_id BIGINT NOT NULL,
    title TEXT NOT NULL,
    occurred_at TIMESTAMPTZ NOT NULL,
    location TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    cover_photo_uri TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN event_records.id IS '主键，事件ID';
COMMENT ON COLUMN event_records.account_id IS '所属账号ID';
COMMENT ON COLUMN event_records.title IS '事件标题';
COMMENT ON COLUMN event_records.occurred_at IS '事件发生时间';
COMMENT ON COLUMN event_records.location IS '事件地点';
COMMENT ON COLUMN event_records.description IS '事件描述';
COMMENT ON COLUMN event_records.cover_photo_uri IS '事件封面图片URI';
COMMENT ON COLUMN event_records.created_at IS '创建时间';
COMMENT ON COLUMN event_records.updated_at IS '更新时间';

CREATE INDEX idx_event_records_account_id ON event_records(account_id);
CREATE INDEX idx_event_records_account_occurred_at ON event_records(account_id, occurred_at DESC);

CREATE TRIGGER trg_set_updated_at_event_records
BEFORE UPDATE ON event_records
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE event_notes (
    id VARCHAR(64) PRIMARY KEY,
    event_id VARCHAR(64) NOT NULL REFERENCES event_records(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    photo_uri TEXT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

COMMENT ON COLUMN event_notes.id IS '主键，事件补充笔记ID';
COMMENT ON COLUMN event_notes.event_id IS '所属事件ID';
COMMENT ON COLUMN event_notes.title IS '补充笔记标题';
COMMENT ON COLUMN event_notes.content IS '补充笔记内容';
COMMENT ON COLUMN event_notes.photo_uri IS '补充笔记图片URI';
COMMENT ON COLUMN event_notes.created_at IS '补充笔记创建时间';
COMMENT ON COLUMN event_notes.updated_at IS '补充笔记更新时间';

CREATE INDEX idx_event_notes_event_id ON event_notes(event_id);
CREATE INDEX idx_event_notes_event_created_at ON event_notes(event_id, created_at DESC);

CREATE TRIGGER trg_set_updated_at_event_notes
BEFORE UPDATE ON event_notes
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
