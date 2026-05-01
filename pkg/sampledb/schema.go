package sampledb

// RDBSchema creates the relational database and tables
const RDBSchema = `
CREATE DATABASE IF NOT EXISTS rdb;

CREATE TABLE IF NOT EXISTS rdb.meter_info (
  meter_id VARCHAR(50) PRIMARY KEY,
  install_date DATE,
  voltage_level VARCHAR(20),
  manufacturer VARCHAR(50),
  status VARCHAR(20),
  area_id VARCHAR(20),
  user_id VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS rdb.user_info (
  user_id VARCHAR(50) PRIMARY KEY,
  user_name VARCHAR(100),
  address VARCHAR(200),
  contact VARCHAR(20)
);

CREATE TABLE IF NOT EXISTS rdb.area_info (
  area_id VARCHAR(20) PRIMARY KEY,
  area_name VARCHAR(100),
  manager VARCHAR(50),
  region VARCHAR(50)
);

CREATE TABLE IF NOT EXISTS rdb.alarm_rules (
  rule_id SERIAL PRIMARY KEY,
  rule_name VARCHAR(100),
  metric VARCHAR(50),
  operator VARCHAR(10),
  threshold FLOAT8,
  severity VARCHAR(20),
  notify_method VARCHAR(50)
);
`

// TSDBSchema creates the time-series database and table
const TSDBSchema = `
CREATE TS DATABASE IF NOT EXISTS tsdb;

CREATE TABLE IF NOT EXISTS tsdb.meter_data (
  ts TIMESTAMPTZ(3) NOT NULL,
  voltage FLOAT8 NULL,
  current FLOAT8 NULL,
  power FLOAT8 NULL,
  energy FLOAT8 NULL
) TAGS (
  meter_id VARCHAR(50) NOT NULL
) PRIMARY TAGS(meter_id)
  retentions 0s
  activetime 1d;
`
