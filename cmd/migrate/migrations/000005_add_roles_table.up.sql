CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    org_id BIGINT NOT NULL,
    permissions VARCHAR(100) [],
    version BIGINT DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,

    CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES organizations (id)
);

CREATE TRIGGER update_roles_updated_at BEFORE UPDATE
ON roles FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();