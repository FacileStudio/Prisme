-- +goose Up
-- baseline: the production schema as GORM AutoMigrate and schemas.Migrate built it up to v0.17.0, captured 2026-11-13 with pg_dump --schema-only --no-owner --no-privileges --no-comments --no-tablespaces --schema=public

CREATE TABLE accounts (
    user_id bigint NOT NULL,
    localpart text NOT NULL,
    matrix_user_id text NOT NULL,
    created_at timestamp with time zone,
    avatar_url text DEFAULT ''::text NOT NULL
);

CREATE TABLE agents (
    id bigint NOT NULL,
    owner_id bigint NOT NULL,
    slug text NOT NULL,
    name text NOT NULL,
    matrix_user_id text NOT NULL,
    device_id text,
    sealed_token text,
    created_at timestamp with time zone,
    deactivated_at timestamp with time zone,
    avatar_url text DEFAULT ''::text NOT NULL,
    avatar_source text DEFAULT ''::text NOT NULL
);

CREATE SEQUENCE agents_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE agents_id_seq OWNED BY agents.id;

CREATE TABLE calls (
    room_id text NOT NULL,
    call_id text NOT NULL,
    event_id text NOT NULL,
    started_by text NOT NULL,
    started_at timestamp with time zone
);

CREATE TABLE channels (
    room_id text NOT NULL,
    space_id bigint NOT NULL,
    name text NOT NULL,
    alias text NOT NULL,
    created_at timestamp with time zone
);

CREATE TABLE directory_deletion_state (
    id integer NOT NULL,
    state bytea NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE erased_subjects (
    subject_hash text NOT NULL,
    erased_at timestamp with time zone NOT NULL
);

CREATE TABLE matrix_sessions (
    porte_session_id bigint NOT NULL,
    user_id bigint NOT NULL,
    device_id text,
    sealed_token text NOT NULL,
    created_at timestamp with time zone
);

CREATE TABLE porte_api_key_usage (
    api_key_id bigint NOT NULL,
    day date DEFAULT CURRENT_DATE NOT NULL,
    count bigint DEFAULT 0 NOT NULL
);

CREATE TABLE porte_api_keys (
    id bigint NOT NULL,
    app text NOT NULL,
    kind text DEFAULT 'secret'::text NOT NULL,
    prefix text NOT NULL,
    token_hash text NOT NULL,
    allowed_origins jsonb DEFAULT '[]'::jsonb NOT NULL,
    daily_quota integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    revoked_at timestamp with time zone
);

CREATE SEQUENCE porte_api_keys_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE porte_api_keys_id_seq OWNED BY porte_api_keys.id;

CREATE TABLE porte_identities (
    user_id bigint NOT NULL,
    provider text NOT NULL,
    subject text NOT NULL,
    password_hash text DEFAULT ''::text NOT NULL,
    access_token text DEFAULT ''::text NOT NULL,
    refresh_token text DEFAULT ''::text NOT NULL,
    token_expiry timestamp with time zone,
    roles jsonb,
    roles_synced_at timestamp with time zone,
    synced_at timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);

CREATE TABLE porte_login_codes (
    code_hash text NOT NULL,
    user_id bigint NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    consumed_at timestamp with time zone,
    provider text DEFAULT ''::text NOT NULL,
    subject text DEFAULT ''::text NOT NULL
);

CREATE TABLE porte_magic_links (
    token_hash text NOT NULL,
    user_id bigint NOT NULL,
    ip text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    expires_at timestamp with time zone NOT NULL,
    consumed_at timestamp with time zone
);

CREATE TABLE porte_sessions (
    id bigint NOT NULL,
    token_hash text NOT NULL,
    user_id bigint NOT NULL,
    label text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    last_used_at timestamp with time zone DEFAULT now() NOT NULL,
    expires_at timestamp with time zone
);

CREATE SEQUENCE porte_sessions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE porte_sessions_id_seq OWNED BY porte_sessions.id;

CREATE TABLE porte_users (
    id bigint NOT NULL,
    facile_id text,
    email text NOT NULL,
    email_verified boolean DEFAULT false NOT NULL,
    name text DEFAULT ''::text NOT NULL,
    avatar_url text DEFAULT ''::text NOT NULL,
    avatar_source text DEFAULT ''::text NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    locale text DEFAULT 'en'::text NOT NULL
);

CREATE SEQUENCE porte_users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE porte_users_id_seq OWNED BY porte_users.id;

CREATE TABLE push_devices (
    token text NOT NULL,
    user_id bigint NOT NULL,
    porte_session_id bigint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);

CREATE TABLE spaces (
    id bigint NOT NULL,
    room_id text NOT NULL,
    name text NOT NULL,
    slug text NOT NULL,
    created_by bigint NOT NULL,
    created_at timestamp with time zone,
    registre_space_id text
);

CREATE SEQUENCE spaces_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;

ALTER SEQUENCE spaces_id_seq OWNED BY spaces.id;

ALTER TABLE ONLY agents ALTER COLUMN id SET DEFAULT nextval('agents_id_seq'::regclass);

ALTER TABLE ONLY porte_api_keys ALTER COLUMN id SET DEFAULT nextval('porte_api_keys_id_seq'::regclass);

ALTER TABLE ONLY porte_sessions ALTER COLUMN id SET DEFAULT nextval('porte_sessions_id_seq'::regclass);

ALTER TABLE ONLY porte_users ALTER COLUMN id SET DEFAULT nextval('porte_users_id_seq'::regclass);

ALTER TABLE ONLY spaces ALTER COLUMN id SET DEFAULT nextval('spaces_id_seq'::regclass);

ALTER TABLE ONLY accounts
    ADD CONSTRAINT accounts_pkey PRIMARY KEY (user_id);

ALTER TABLE ONLY agents
    ADD CONSTRAINT agents_pkey PRIMARY KEY (id);

ALTER TABLE ONLY calls
    ADD CONSTRAINT calls_pkey PRIMARY KEY (room_id);

ALTER TABLE ONLY channels
    ADD CONSTRAINT channels_pkey PRIMARY KEY (room_id);

ALTER TABLE ONLY directory_deletion_state
    ADD CONSTRAINT directory_deletion_state_pkey PRIMARY KEY (id);

ALTER TABLE ONLY erased_subjects
    ADD CONSTRAINT erased_subjects_pkey PRIMARY KEY (subject_hash);

ALTER TABLE ONLY matrix_sessions
    ADD CONSTRAINT matrix_sessions_pkey PRIMARY KEY (porte_session_id);

ALTER TABLE ONLY porte_api_key_usage
    ADD CONSTRAINT porte_api_key_usage_pkey PRIMARY KEY (api_key_id, day);

ALTER TABLE ONLY porte_api_keys
    ADD CONSTRAINT porte_api_keys_pkey PRIMARY KEY (id);

ALTER TABLE ONLY porte_api_keys
    ADD CONSTRAINT porte_api_keys_token_hash_key UNIQUE (token_hash);

ALTER TABLE ONLY porte_identities
    ADD CONSTRAINT porte_identities_pkey PRIMARY KEY (provider, subject);

ALTER TABLE ONLY porte_login_codes
    ADD CONSTRAINT porte_login_codes_pkey PRIMARY KEY (code_hash);

ALTER TABLE ONLY porte_magic_links
    ADD CONSTRAINT porte_magic_links_pkey PRIMARY KEY (token_hash);

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_pkey PRIMARY KEY (id);

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_token_hash_key UNIQUE (token_hash);

ALTER TABLE ONLY porte_users
    ADD CONSTRAINT porte_users_email_key UNIQUE (email);

ALTER TABLE ONLY porte_users
    ADD CONSTRAINT porte_users_facile_id_key UNIQUE (facile_id);

ALTER TABLE ONLY porte_users
    ADD CONSTRAINT porte_users_pkey PRIMARY KEY (id);

ALTER TABLE ONLY push_devices
    ADD CONSTRAINT push_devices_pkey PRIMARY KEY (token);

ALTER TABLE ONLY spaces
    ADD CONSTRAINT spaces_pkey PRIMARY KEY (id);

CREATE UNIQUE INDEX agents_owner_slug ON agents USING btree (owner_id, slug);

CREATE UNIQUE INDEX idx_accounts_localpart ON accounts USING btree (localpart);

CREATE UNIQUE INDEX idx_accounts_matrix_user_id ON accounts USING btree (matrix_user_id);

CREATE UNIQUE INDEX idx_agents_matrix_user_id ON agents USING btree (matrix_user_id);

CREATE INDEX idx_agents_owner_id ON agents USING btree (owner_id);

CREATE UNIQUE INDEX idx_agents_slug ON agents USING btree (slug);

CREATE UNIQUE INDEX idx_channels_alias ON channels USING btree (alias);

CREATE INDEX idx_channels_space_id ON channels USING btree (space_id);

CREATE INDEX idx_matrix_sessions_user_id ON matrix_sessions USING btree (user_id);

CREATE UNIQUE INDEX idx_spaces_registre_space_id ON spaces USING btree (registre_space_id);

CREATE UNIQUE INDEX idx_spaces_room_id ON spaces USING btree (room_id);

CREATE UNIQUE INDEX idx_spaces_slug ON spaces USING btree (slug);

CREATE INDEX porte_api_keys_app_idx ON porte_api_keys USING btree (app);

CREATE INDEX porte_api_keys_token_hash_idx ON porte_api_keys USING btree (token_hash) WHERE (revoked_at IS NULL);

CREATE INDEX porte_identities_user_idx ON porte_identities USING btree (user_id);

CREATE INDEX porte_magic_links_expiry_idx ON porte_magic_links USING btree (expires_at);

CREATE INDEX porte_magic_links_ip_idx ON porte_magic_links USING btree (ip, created_at) WHERE (ip <> ''::text);

CREATE INDEX porte_magic_links_user_idx ON porte_magic_links USING btree (user_id, created_at);

CREATE INDEX porte_sessions_expiry_idx ON porte_sessions USING btree (expires_at) WHERE (expires_at IS NOT NULL);

CREATE INDEX porte_sessions_user_idx ON porte_sessions USING btree (user_id);

CREATE INDEX push_devices_user_idx ON push_devices USING btree (user_id);

ALTER TABLE ONLY porte_api_key_usage
    ADD CONSTRAINT porte_api_key_usage_api_key_id_fkey FOREIGN KEY (api_key_id) REFERENCES porte_api_keys(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_identities
    ADD CONSTRAINT porte_identities_user_id_fkey FOREIGN KEY (user_id) REFERENCES porte_users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_login_codes
    ADD CONSTRAINT porte_login_codes_user_id_fkey FOREIGN KEY (user_id) REFERENCES porte_users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_magic_links
    ADD CONSTRAINT porte_magic_links_user_id_fkey FOREIGN KEY (user_id) REFERENCES porte_users(id) ON DELETE CASCADE;

ALTER TABLE ONLY porte_sessions
    ADD CONSTRAINT porte_sessions_user_id_fkey FOREIGN KEY (user_id) REFERENCES porte_users(id) ON DELETE CASCADE;

ALTER TABLE ONLY push_devices
    ADD CONSTRAINT push_devices_porte_session_id_fkey FOREIGN KEY (porte_session_id) REFERENCES porte_sessions(id) ON DELETE CASCADE;

ALTER TABLE ONLY push_devices
    ADD CONSTRAINT push_devices_user_id_fkey FOREIGN KEY (user_id) REFERENCES porte_users(id) ON DELETE CASCADE;

-- +goose Down
SELECT 'baseline migration is not reversible' WHERE false;
