package vo

// Root structure
type AlarmEvent struct {
	Time   string `json:"time"`
	Region string `json:"region"`
	Detail Detail `json:"detail"`
}

// Detail structure
type Detail struct {
	AlarmName     string `json:"alarmName"`
	State         State  `json:"state"`
	Configuration Config `json:"configuration"`
}

// State structure
type State struct {
	Value      string `json:"value"`
	Reason     string `json:"reason"`
	ReasonData string `json:"reasonData"`
	Timestamp  string `json:"timestamp"`
}

// ReasonData structure
type ReasonData struct {
	Version             string               `json:"version"`
	QueryDate           string               `json:"queryDate"`
	StartDate           string               `json:"startDate"`
	Statistic           string               `json:"statistic"`
	Period              int                  `json:"period"`
	RecentDatapoints    []float64            `json:"recentDatapoints"`
	Threshold           float64              `json:"threshold"`
	EvaluatedDatapoints []EvaluatedDatapoint `json:"evaluatedDatapoints"`
}

// EvaluatedDatapoint structure
type EvaluatedDatapoint struct {
	Timestamp   string  `json:"timestamp"`
	SampleCount float64 `json:"sampleCount"`
	Value       float64 `json:"value"`
}

// Config structure
type Config struct {
	Metrics []Metric `json:"metrics"`
}

// Metric structure
type Metric struct {
	ID         string     `json:"id"`
	MetricStat MetricStat `json:"metricStat"`
	ReturnData bool       `json:"returnData"`
}

// MetricStat structure
type MetricStat struct {
	Metric MetricDetails `json:"metric"`
	Period int           `json:"period"`
	Stat   string        `json:"stat"`
}

// MetricDetails structure
type MetricDetails struct {
	Namespace  string            `json:"namespace"`
	Name       string            `json:"name"`
	Dimensions map[string]string `json:"dimensions"`
}
