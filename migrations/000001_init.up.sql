-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enums
CREATE TYPE user_role AS ENUM ('admin', 'supervisor', 'client');
CREATE TYPE project_status AS ENUM ('active', 'on_hold', 'completed', 'archived');
CREATE TYPE log_status AS ENUM ('draft', 'submitted', 'approved', 'rejected');
CREATE TYPE issue_priority AS ENUM ('low', 'medium', 'high');
CREATE TYPE issue_status AS ENUM ('open', 'in_progress', 'resolved');

-- Users
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name          VARCHAR(255) NOT NULL,
    email         VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role          user_role NOT NULL DEFAULT 'supervisor',
    avatar_url    TEXT,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);

-- Refresh tokens
CREATE TABLE refresh_tokens (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);

-- Projects
CREATE TABLE projects (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    location    VARCHAR(500),
    start_date  DATE,
    end_date    DATE,
    status      project_status NOT NULL DEFAULT 'active',
    created_by  UUID NOT NULL REFERENCES users(id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_created_by ON projects(created_by);

-- Project members
CREATE TABLE project_members (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role       VARCHAR(50) NOT NULL DEFAULT 'supervisor',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(project_id, user_id)
);

CREATE INDEX idx_project_members_project ON project_members(project_id);
CREATE INDEX idx_project_members_user ON project_members(user_id);

-- Workers (site labour, not app users)
CREATE TABLE workers (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       VARCHAR(255) NOT NULL,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    phone      VARCHAR(50),
    trade      VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_workers_project ON workers(project_id);

-- Attendance
CREATE TABLE attendance (
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    worker_id  UUID NOT NULL REFERENCES workers(id) ON DELETE CASCADE,
    check_in   TIMESTAMPTZ NOT NULL,
    check_out  TIMESTAMPTZ,
    notes      TEXT,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_attendance_project ON attendance(project_id);
CREATE INDEX idx_attendance_worker ON attendance(worker_id);
CREATE INDEX idx_attendance_checkin ON attendance(check_in);

-- Daily logs
CREATE TABLE daily_logs (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    created_by    UUID NOT NULL REFERENCES users(id),
    date          DATE NOT NULL,
    weather       VARCHAR(255),
    notes         TEXT,
    gps_lat       DOUBLE PRECISION,
    gps_lng       DOUBLE PRECISION,
    signature_url TEXT,
    status        log_status NOT NULL DEFAULT 'draft',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_daily_logs_project ON daily_logs(project_id);
CREATE INDEX idx_daily_logs_created_by ON daily_logs(created_by);
CREATE INDEX idx_daily_logs_date ON daily_logs(date);
CREATE INDEX idx_daily_logs_status ON daily_logs(status);

-- Media (images, videos, audio uploaded to GitHub)
CREATE TABLE media (
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id   UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    uploaded_by  UUID NOT NULL REFERENCES users(id),
    filename     VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    url          TEXT NOT NULL,
    github_sha   VARCHAR(255) NOT NULL,
    github_key   TEXT NOT NULL,
    size         BIGINT NOT NULL DEFAULT 0,
    log_id       UUID REFERENCES daily_logs(id) ON DELETE SET NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_media_project ON media(project_id);
CREATE INDEX idx_media_log ON media(log_id);
CREATE INDEX idx_media_uploaded_by ON media(uploaded_by);

-- Issues
CREATE TABLE issues (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id  UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    title       VARCHAR(500) NOT NULL,
    description TEXT,
    priority    issue_priority NOT NULL DEFAULT 'medium',
    status      issue_status NOT NULL DEFAULT 'open',
    assigned_to UUID REFERENCES users(id) ON DELETE SET NULL,
    reported_by UUID NOT NULL REFERENCES users(id),
    gps_lat     DOUBLE PRECISION,
    gps_lng     DOUBLE PRECISION,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_issues_project ON issues(project_id);
CREATE INDEX idx_issues_status ON issues(status);
CREATE INDEX idx_issues_priority ON issues(priority);
CREATE INDEX idx_issues_assigned_to ON issues(assigned_to);
