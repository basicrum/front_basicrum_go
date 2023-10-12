CREATE TABLE IF NOT EXISTS {prefix}webperf_rum_own_hostnames (
    username                        LowCardinality(String),
    hostname                        LowCardinality(String),
    subscription_id                 String,
    subscription_expire_at          DateTime64(3) NOT NULL,
    updated_at                      DateTime64(3) DEFAULT now(),
    INDEX index_username username TYPE bloom_filter GRANULARITY 1
)
ENGINE = ReplacingMergeTree(updated_at)
ORDER BY hostname
PARTITION BY hostname
PRIMARY KEY hostname
