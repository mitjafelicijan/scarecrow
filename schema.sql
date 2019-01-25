create table if not exists stats (
    ts                        int     default (current_timestamp) 
                                      primary key,
    uptime_sec                real,
    total_status_code_count   text,
    total_count               integer,
    total_response_time_sec   real,
    average_response_time_sec real
);
