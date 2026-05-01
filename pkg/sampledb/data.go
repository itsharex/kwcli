package sampledb

// GenerateRDBData returns SQL to populate relationship tables
func GenerateRDBData() string {
	return `
-- Clear existing data
TRUNCATE TABLE rdb.area_info CASCADE;
TRUNCATE TABLE rdb.user_info CASCADE;
TRUNCATE TABLE rdb.meter_info CASCADE;
TRUNCATE TABLE rdb.alarm_rules CASCADE;

-- Insert area_info (100 areas, 4 regions)
INSERT INTO rdb.area_info (area_id, area_name, manager, region)
SELECT
  'A' || s::text,
  'Area ' || s::text,
  'Manager ' || ((s-1)/25 + 1)::text,
  CASE WHEN s <= 25 THEN 'North'
       WHEN s <= 50 THEN 'South'
       WHEN s <= 75 THEN 'East'
       ELSE 'West' END
FROM generate_series(1, 100) AS s;

-- Insert user_info (100 users)
INSERT INTO rdb.user_info (user_id, user_name, address, contact)
SELECT
  'U' || s::text,
  'User ' || s::text,
  'Address ' || s::text || ', Street ' || (s%10 + 1)::text,
  '13800' || LPAD(s::text, 5, '0')
FROM generate_series(1, 100) AS s;

-- Insert meter_info (100 meters, 5 faults)
INSERT INTO rdb.meter_info (meter_id, install_date, voltage_level, manufacturer, status, area_id, user_id)
SELECT
  'M' || s::text,
  DATE '2020-01-01' + (s % 365)::int,
  CASE WHEN s % 2 = 0 THEN '380V' ELSE '220V' END,
  CASE WHEN s % 3 = 0 THEN 'Factory A'
       WHEN s % 3 = 1 THEN 'Factory B'
       ELSE 'Factory C' END,
  CASE WHEN s IN (20, 40, 60, 80, 100) THEN 'Fault' ELSE 'Normal' END,
  'A' || s::text,
  'U' || s::text
FROM generate_series(1, 100) AS s;

-- Insert alarm_rules
INSERT INTO rdb.alarm_rules (rule_name, metric, operator, threshold, severity, notify_method) VALUES
  ('高压告警', 'voltage', '>', 230, 'High', 'SMS'),
  ('低压告警', 'voltage', '<', 210, 'High', 'SMS'),
  ('过流告警', 'current', '>', 6.5, 'Critical', 'SMS+Email'),
  ('过载告警', 'power', '>', 1800, 'Medium', 'Email');
`
}

// GenerateTSDBData returns SQL to populate time-series table
func GenerateTSDBData() string {
	return `
-- Clear existing time-series data
DELETE FROM tsdb.meter_data WHERE 1=1;

-- Insert 10,000 time-series readings
INSERT INTO tsdb.meter_data(ts, voltage, current, power, energy, meter_id)
SELECT
  NOW() - (s * 10)::int * INTERVAL '1 minute',
  220.0 + (s % 10)::float,
  5.0 + (s % 15)::float * 0.1,
  1000.0 + (s % 20)::float * 50,
  5000.0 + s::float * 10,
  'M' || ((s % 100) + 1)::text
FROM generate_series(1, 10000) AS s;
`
}
