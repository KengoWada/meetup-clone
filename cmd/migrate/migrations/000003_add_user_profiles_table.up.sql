CREATE TABLE IF NOT EXISTS user_profiles (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    profile_pic VARCHAR(255) NULL,
    date_of_birth DATE NOT NULL,
    user_id BIGINT NOT NULL UNIQUE,
    version BIGINT DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,
    deleted_at TIMESTAMP(0) WITH TIME ZONE DEFAULT NULL,

    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TRIGGER update_user_profiles_updated_at BEFORE UPDATE
ON user_profiles FOR EACH ROW EXECUTE PROCEDURE 
update_updated_at_column();