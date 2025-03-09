CREATE TABLE IF NOT EXISTS organizations (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT NOT NULL,
    profile_pic VARCHAR(255) NULL,
    is_active BOOLEAN DEFAULT TRUE,
    version BIGINT DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL
);

CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE
ON organizations FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();