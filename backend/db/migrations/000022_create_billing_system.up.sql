-- +migrate Up
-- Create user_subscriptions table
CREATE TABLE IF NOT EXISTS user_subscriptions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tier VARCHAR(50) NOT NULL DEFAULT 'free',
    status VARCHAR(50) NOT NULL DEFAULT 'active',
    stripe_customer_id VARCHAR(255),
    stripe_subscription_id VARCHAR(255),
    stripe_session_id VARCHAR(255),
    agents_limit INTEGER NOT NULL DEFAULT 3,
    organizations_limit INTEGER NOT NULL DEFAULT 1,
    current_period_start TIMESTAMP,
    current_period_end TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    INDEX idx_user_subscriptions_user_id (user_id),
    INDEX idx_user_subscriptions_status (status),
    INDEX idx_user_subscriptions_stripe_customer_id (stripe_customer_id),
    INDEX idx_user_subscriptions_stripe_subscription_id (stripe_subscription_id)
);

-- Create payment_history table
CREATE TABLE IF NOT EXISTS payment_history (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id VARCHAR(36) NOT NULL REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    stripe_payment_intent_id VARCHAR(255),
    amount DECIMAL(10, 2) NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'usd',
    status VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    
    INDEX idx_payment_history_user_id (user_id),
    INDEX idx_payment_history_subscription_id (subscription_id),
    INDEX idx_payment_history_status (status)
);

-- Create usage_records table
CREATE TABLE IF NOT EXISTS usage_records (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feature VARCHAR(100) NOT NULL,
    count INTEGER NOT NULL DEFAULT 1,
    period VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_usage_records_user_id (user_id),
    INDEX idx_usage_records_period (period),
    UNIQUE INDEX idx_usage_records_user_feature_period (user_id, feature, period)
);

-- Create subscription_features table
CREATE TABLE IF NOT EXISTS subscription_features (
    id VARCHAR(36) PRIMARY KEY,
    tier VARCHAR(50) NOT NULL,
    feature VARCHAR(100) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_subscription_features_tier (tier),
    UNIQUE INDEX idx_subscription_features_tier_feature (tier, feature)
);

-- Insert default subscription features
INSERT INTO subscription_features (id, tier, feature, description, enabled) VALUES
-- Free tier
('free-basic-tasks', 'free', 'basic_tasks', 'Basic task management', TRUE),
('free-basic-agents', 'free', 'basic_agents', 'Up to 3 agents', TRUE),
('free-community-support', 'free', 'community_support', 'Community support forum', TRUE),

-- Starter tier
('starter-all-free', 'starter', 'all_free_features', 'All free tier features', TRUE),
('starter-agents', 'starter', 'agents_limit', 'Up to 10 agents', TRUE),
('starter-orgs', 'starter', 'organizations_limit', 'Up to 3 organizations', TRUE),
('starter-analytics', 'starter', 'basic_analytics', 'Basic usage analytics', TRUE),
('starter-email-support', 'starter', 'email_support', 'Email support with 24h response', TRUE),

-- Pro tier
('pro-all-starter', 'pro', 'all_starter_features', 'All starter tier features', TRUE),
('pro-agents', 'pro', 'unlimited_agents', 'Unlimited agents', TRUE),
('pro-orgs', 'pro', 'unlimited_organizations', 'Unlimited organizations', TRUE),
('pro-analytics', 'pro', 'advanced_analytics', 'Advanced analytics and insights', TRUE),
('pro-integrations', 'pro', 'custom_integrations', 'Custom integrations', TRUE),
('pro-api', 'pro', 'api_access', 'Full API access', TRUE),
('pro-priority-support', 'pro', 'priority_support', 'Priority support with 4h response', TRUE),
('pro-priority-execution', 'pro', 'priority_execution', 'Priority agent execution', TRUE),

-- Enterprise tier
('enterprise-all-pro', 'enterprise', 'all_pro_features', 'All pro tier features', TRUE),
('enterprise-sso', 'enterprise', 'sso_saml', 'Single Sign-On (SAML)', TRUE),
('enterprise-security', 'enterprise', 'advanced_security', 'Advanced security features', TRUE),
('enterprise-dedicated', 'enterprise', 'dedicated_support', 'Dedicated support manager', TRUE),
('enterprise-white-label', 'enterprise', 'white_label', 'White-label options', TRUE),
('enterprise-sla', 'enterprise', 'sla_guarantee', 'SLA guarantee', TRUE),
('enterprise-audit', 'enterprise', 'audit_logs', 'Comprehensive audit logs', TRUE);

-- +migrate Down
DROP TABLE IF EXISTS subscription_features;
DROP TABLE IF EXISTS usage_records;
DROP TABLE IF EXISTS payment_history;
DROP TABLE IF EXISTS user_subscriptions;
