CREATE TABLE IF NOT EXISTS {prefix}webperf_rum_own_hostnames (
    username                        LowCardinality(String),
    hostname                        LowCardinality(String),
    subscription_id                 String,
    subscription_expire_at          DateTime64(3) NOT NULL
)
ENGINE = MergeTree
ORDER BY (username, hostname)
PARTITION BY username
