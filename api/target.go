package api

// Target is part of Panel
type Target struct {
	Alias               interface{}    `json:"alias,omitempty"`
	BucketAggs          interface{}    `json:"bucketAggs,omitempty"`
	DataSource          DashDataSource `json:"datasource"`
	Format              string         `json:"format,omitempty"`
	DisableTextWrap     interface{}    `json:"disableTextWrap,omitempty"`
	EditorMode          string         `json:"editorMode,omitempty"`
	Expr                string         `json:"expr,omitempty"`
	FullMetaSearch      interface{}    `json:"fullMetaSearch,omitempty"`
	Hide                interface{}    `json:"hide,omitempty"`
	LuceneQueryType     interface{}    `json:"luceneQueryType,omitempty"`
	Metrics             interface{}    `json:"metrics,omitempty"`
	Query               interface{}    `json:"query,omitempty"`
	QueryType           interface{}    `json:"queryType,omitempty"`
	IncludeNullMetadata interface{}    `json:"includeNullMetadata,omitempty"`
	Instant             interface{}    `json:"instant,omitempty"`
	Interval            interface{}    `json:"interval,omitempty"`
	LegendFormat        interface{}    `json:"legendFormat,omitempty"`
	Range               interface{}    `json:"range,omitempty"`
	RefId               string         `json:"refId"`
	TimeField           interface{}    `json:"timeField,omitempty"`
	UseBackend          interface{}    `json:"useBackend,omitempty"`
}
