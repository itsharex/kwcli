package sampledb

// Scenario represents a sample query scenario
type Scenario struct {
	Name        string
	Title       string
	Description string
	SQL         string
}

// AllScenarios contains all smart-meter scenario queries
var AllScenarios = []Scenario{
	{
		Name:        "top10-area-energy",
		Title:       "查询区域用电量TOP10",
		Description: "统计各区域总用电量，取前10名",
		SQL: `SELECT
  a.area_name,
  SUM(md.energy) AS total_energy
FROM tsdb.meter_data md
JOIN rdb.meter_info mi ON md.meter_id = mi.meter_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
GROUP BY a.area_name
ORDER BY total_energy DESC
LIMIT 10;`,
	},
	{
		Name:        "fault-meters",
		Title:       "查询故障电表及用户信息",
		Description: "查询状态为Fault的电表及其关联的用户信息",
		SQL: `SELECT
  mi.meter_id,
  u.user_name,
  u.contact,
  a.area_name
FROM rdb.meter_info mi
JOIN rdb.user_info u ON mi.user_id = u.user_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
WHERE mi.status = 'Fault';`,
	},
	{
		Name:        "meter-summary",
		Title:       "电表概要查询",
		Description: "查询指定电表的详细信息及数据点数量",
		SQL: `SELECT
  mi.meter_id,
  mi.voltage_level,
  mi.status,
  u.user_name,
  a.area_name,
  (SELECT COUNT(*)
   FROM tsdb.meter_data md
   WHERE md.meter_id = mi.meter_id) AS data_points
FROM rdb.meter_info mi
JOIN rdb.user_info u ON mi.user_id = u.user_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
WHERE mi.meter_id = 'M1';`,
	},
	{
		Name:        "alarm-detection",
		Title:       "告警检测查询",
		Description: "根据告警规则检测异常数据",
		SQL: `SELECT
  md.meter_id,
  md.ts,
  ar.rule_name,
  md.voltage,
  md.current,
  md.power
FROM tsdb.meter_data md
JOIN rdb.alarm_rules ar ON 1=1
WHERE (ar.metric = 'voltage'
       AND ((ar.operator = '>' AND md.voltage < ar.threshold)
            OR (ar.operator = '<' AND md.voltage > ar.threshold)))
   OR (ar.metric = 'current' AND md.current > ar.threshold)
   OR (ar.metric = 'power' AND md.power > ar.threshold)
ORDER BY md.ts DESC
LIMIT 100;`,
	},
	{
		Name:        "area-energy-stats",
		Title:       "区域用电量统计",
		Description: "按区域统计总用电量和平均功率",
		SQL: `SELECT
  a.region,
  a.area_name,
  SUM(md.energy) AS total_energy,
  AVG(md.power) AS avg_power
FROM tsdb.meter_data md
JOIN rdb.meter_info mi ON md.meter_id = mi.meter_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
GROUP BY a.region, a.area_name;`,
	},
	{
		Name:        "recent-24h-trend",
		Title:       "查询指定电表最近24小时用电趋势量",
		Description: "查询M1电表最近24小时的功率和能耗",
		SQL: `SELECT
  md.ts,
  md.power,
  md.energy
FROM tsdb.meter_data md
WHERE md.meter_id = 'M1'
  AND md.ts > NOW() - INTERVAL '24 hours'
ORDER BY md.ts;`,
	},
	{
		Name:        "time-bucket-stats",
		Title:       "分时负荷统计",
		Description: "按1小时粒度统计重点电表的功率指标",
		SQL: `SELECT
  meter_id,
  time_bucket(ts, '1h') AS bucket_start,
  COUNT(*) AS sample_count,
  AVG(power) AS avg_power,
  MAX(power) AS max_power
FROM tsdb.meter_data
WHERE meter_id IN ('M1', 'M2', 'M3')
GROUP BY meter_id, bucket_start
ORDER BY meter_id, bucket_start;`,
	},
	{
		Name:        "session-analysis",
		Title:       "用电会话分析",
		Description: "以30分钟空闲间隔划分会话窗口，查看连续用电过程",
		SQL: `SELECT
  meter_id,
  first(ts) AS session_start,
  last(ts) AS session_end,
  COUNT(*) AS sample_count,
  SUM(energy) AS total_energy
FROM tsdb.meter_data
WHERE meter_id = 'M1'
GROUP BY meter_id, session_window(ts, '30m')
ORDER BY session_start;`,
	},
	{
		Name:        "voltage-state",
		Title:       "电压状态持续分析",
		Description: "按是否高压切分连续区间，观察状态持续时间",
		SQL: `SELECT
  meter_id,
  first(ts) AS window_start,
  last(ts) AS window_end,
  COUNT(*) AS sample_count,
  MIN(voltage) AS min_voltage,
  MAX(voltage) AS max_voltage
FROM tsdb.meter_data
WHERE meter_id = 'M1'
GROUP BY meter_id, state_window(CASE WHEN voltage >= 225 THEN 'high' ELSE 'low' END)
ORDER BY window_start;`,
	},
	{
		Name:        "abnormal-current",
		Title:       "异常电流事件识别",
		Description: "识别电流升高到6A及以上并回落到5.3A及以下的异常事件",
		SQL: `SELECT
  meter_id,
  first(ts) AS event_start,
  last(ts) AS event_end,
  COUNT(*) AS sample_count,
  MAX(current) AS peak_current,
  AVG(power) AS avg_power
FROM tsdb.meter_data
WHERE meter_id = 'M1'
GROUP BY meter_id, event_window(current >= 6, current <= 5.3)
ORDER BY event_start;`,
	},
	{
		Name:        "sliding-window",
		Title:       "滑动采样趋势分析",
		Description: "每12条采样做窗口，以6条为滑动步长观察功率变化",
		SQL: `SELECT
  meter_id,
  first(ts) AS window_start,
  last(ts) AS window_end,
  COUNT(*) AS sample_count,
  AVG(power) AS avg_power,
  MAX(power) AS max_power
FROM tsdb.meter_data
WHERE meter_id = 'M1'
GROUP BY meter_id, count_window(12, 6)
ORDER BY window_start;`,
	},
}

// FindScenario finds a scenario by name
func FindScenario(name string) *Scenario {
	for i := range AllScenarios {
		if AllScenarios[i].Name == name {
			return &AllScenarios[i]
		}
	}
	return nil
}
