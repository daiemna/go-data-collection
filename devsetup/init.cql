CREATE KEYSPACE IF NOT EXISTS ks_data_collection WITH REPLICATION = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
CREATE TABLE IF NOT EXISTS ks_data_collection.ts_timeseries (
    tsid uuid,
    time timestamp,
    value float,
    PRIMARY KEY (tsid, time)
) WITH CLUSTERING ORDER BY (time ASC);