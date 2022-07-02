CREATE TABLE IF NOT EXISTS integration_test_webperf_rum_events (
    event_date Date DEFAULT toDate(created_at),
    created_at DateTime,
    event_type                      LowCardinality(String),
    browser_name                    LowCardinality(String),
    browser_version                 String,
    device_manufacturer             LowCardinality(String),
    device_type                     LowCardinality(String),
    user_agent                      String,
    next_hop_protocol               LowCardinality(String),
    visibility_state                LowCardinality(String),

    session_id                      FixedString(43),
    session_length                  UInt8,
    url                             String,
    connect_duration                Nullable(UInt16),
    dns_duration                    Nullable(UInt16),
    first_byte_duration             Nullable(UInt16),
    redirect_duration               Nullable(UInt16),
    redirects_count                 UInt8,
    
    first_contentful_paint          Nullable(UInt16),
    first_paint                     Nullable(UInt16),

    cumulative_layout_shift         Nullable(Float32),
    first_input_delay               Nullable(UInt16),
    largest_contentful_paint        Nullable(UInt16),

    country_code                    FixedString(2)
)
    ENGINE = MergeTree()
    PARTITION BY toYYYYMMDD(event_date)
    ORDER BY (device_type, event_date)
    SETTINGS index_granularity = 8192