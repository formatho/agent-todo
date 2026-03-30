-- +migrate Down
DROP TABLE IF EXISTS subscription_features;
DROP TABLE IF EXISTS usage_records;
DROP TABLE IF EXISTS payment_history;
DROP TABLE IF EXISTS user_subscriptions;
