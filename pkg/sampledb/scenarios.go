package sampledb

// Scenario represents a sample query scenario
type Scenario struct {
	Name        string
	Title       string
	Description string
	SQL         string
	Category    string // "basic", "cross-mode", "window"
}

// AllScenarios contains all smart-meter scenario queries
var AllScenarios = []Scenario{
	// === 基础查询场景 (Basic Queries) ===
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
		Category: "basic",
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
		Category: "basic",
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
		Category: "basic",
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
		Category: "basic",
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
		Category: "basic",
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
		Category: "basic",
	},
	// === 跨模查询场景 (Cross-Mode Queries) ===
	{
		Name:        "cross-mode-join",
		Title:       "跨模关联查询",
		Description: "关联时序数据与关系数据，分析电表与用户关联的用电行为",
		SQL: `SELECT
  mi.meter_id,
  u.user_name,
  a.area_name,
  a.region,
  COUNT(md.ts) AS reading_count,
  AVG(md.power) AS avg_power,
  SUM(md.energy) AS total_energy,
  MAX(md.ts) AS last_reading
FROM tsdb.meter_data md
JOIN rdb.meter_info mi ON md.meter_id = mi.meter_id
JOIN rdb.user_info u ON mi.user_id = u.user_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
WHERE md.ts > NOW() - INTERVAL '1 day'
GROUP BY mi.meter_id, u.user_name, a.area_name, a.region
ORDER BY total_energy DESC
LIMIT 20;`,
		Category: "cross-mode",
	},
	{
		Name:        "cross-mode-user-power",
		Title:       "用户用电排名分析",
		Description: "统计用户总用电量并进行排名，结合用户信息和区域信息",
		SQL: `SELECT
  u.user_id,
  u.user_name,
  u.contact,
  a.area_name,
  a.region,
  SUM(md.energy) AS total_energy,
  AVG(md.power) AS avg_power,
  MAX(md.power) AS peak_power
FROM tsdb.meter_data md
JOIN rdb.meter_info mi ON md.meter_id = mi.meter_id
JOIN rdb.user_info u ON mi.user_id = u.user_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
GROUP BY u.user_id, u.user_name, u.contact, a.area_name, a.region
ORDER BY total_energy DESC
LIMIT 30;`,
		Category: "cross-mode",
	},
	{
		Name:        "cross-mode-alarm-analysis",
		Title:       "跨模告警分析",
		Description: "分析告警触发的时序数据，结合电表信息和用户信息",
		SQL: `SELECT
  md.meter_id,
  mi.status AS meter_status,
  u.user_name,
  a.area_name,
  ar.rule_name,
  ar.severity,
  md.ts,
  md.voltage,
  md.current,
  md.power,
  CASE 
    WHEN ar.metric = 'voltage' AND ar.operator = '>' THEN md.voltage - ar.threshold
    WHEN ar.metric = 'voltage' AND ar.operator = '<' THEN ar.threshold - md.voltage
    WHEN ar.metric = 'current' THEN md.current - ar.threshold
    WHEN ar.metric = 'power' THEN md.power - ar.threshold
    ELSE 0
  END AS deviation
FROM tsdb.meter_data md
JOIN rdb.meter_info mi ON md.meter_id = mi.meter_id
JOIN rdb.user_info u ON mi.user_id = u.user_id
JOIN rdb.area_info a ON mi.area_id = a.area_id
JOIN rdb.alarm_rules ar ON 1=1
WHERE (ar.metric = 'voltage' AND ((ar.operator = '>' AND md.voltage > ar.threshold) OR (ar.operator = '<' AND md.voltage < ar.threshold)))
   OR (ar.metric = 'current' AND md.current > ar.threshold)
   OR (ar.metric = 'power' AND md.power > ar.threshold)
ORDER BY ar.severity DESC, deviation DESC
LIMIT 50;`,
		Category: "cross-mode",
	},
	{
		Name:        "cross-mode-region-comparison",
		Title:       "区域用电对比分析",
		Description: "对比不同区域的用电情况，结合区域管理员信息",
		SQL: `SELECT
  a.region,
  a.area_id,
  a.area_name,
  a.manager,
  COUNT(DISTINCT mi.meter_id) AS meter_count,
  COUNT(md.ts) AS total_readings,
  SUM(md.energy) AS total_energy,
  AVG(md.power) AS avg_power,
  MAX(md.power) AS max_power,
  MIN(md.voltage) AS min_voltage,
  MAX(md.voltage) AS max_voltage
FROM rdb.area_info a
LEFT JOIN rdb.meter_info mi ON a.area_id = mi.area_id
LEFT JOIN tsdb.meter_data md ON mi.meter_id = md.meter_id
WHERE md.ts > NOW() - INTERVAL '7 days'
GROUP BY a.region, a.area_id, a.area_name, a.manager
ORDER BY a.region, total_energy DESC;`,
		Category: "cross-mode",
	},
	// === 时间窗口函数 (Time Window Functions) ===
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
		Category: "window",
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
		Category: "window",
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
		Category: "window",
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
		Category: "window",
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
		Category: "window",
	},
	{
		Name:        "time-window-advanced",
		Title:       "高级时间窗口分析",
		Description: "使用时间窗口函数进行滑动统计，分析用电趋势",
		SQL: `SELECT
  meter_id,
  time_bucket(ts, '1h') AS bucket,
  AVG(power) AS avg_power,
  SUM(energy) AS energy_sum,
  COUNT(*) AS samples
FROM tsdb.meter_data
WHERE meter_id = 'M1' AND ts > NOW() - INTERVAL '24 hours'
GROUP BY meter_id, time_bucket(ts, '1h')
ORDER BY bucket;`,
		Category: "window",
	},
	{
		Name:        "count-window-example",
		Title:       "计数窗口函数示例",
		Description: "使用计数窗口进行滑动统计，每50条数据计算一次统计指标",
		SQL: `SELECT
  meter_id,
  first(ts) AS window_start,
  last(ts) AS window_end,
  COUNT(*) AS sample_count,
  AVG(voltage) AS avg_voltage,
  AVG(current) AS avg_current,
  AVG(power) AS avg_power
FROM tsdb.meter_data
WHERE meter_id = 'M1'
GROUP BY meter_id, count_window(50, 25)
ORDER BY window_start
LIMIT 20;`,
		Category: "window",
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
