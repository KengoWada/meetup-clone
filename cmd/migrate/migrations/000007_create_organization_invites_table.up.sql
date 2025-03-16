CREATE TABLE IF NOT EXISTS organization_invites (
    id BIGSERIAL PRIMARY KEY,
    org_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role_id BIGINT NOT NULL,
    accepted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    declined_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    version BIGINT DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,

    CONSTRAINT fk_org FOREIGN KEY (org_id) REFERENCES organizations (id),
    CONSTRAINT fk_user_profile FOREIGN KEY (user_id) REFERENCES user_profiles (id),
    CONSTRAINT fk_role FOREIGN KEY (role_id) REFERENCES roles (id)
);

CREATE TRIGGER update_organization_invites_updated_at BEFORE UPDATE
ON organization_invites FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();