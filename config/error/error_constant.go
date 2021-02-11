package error

const (
	ChangeToRootDirFailure    = "Failed to change to root directory: "
	ConnectToDatabaseFailure  = "Failed to connect database: "
	DeleteRedisJobFailure     = "Failed to delete redis job: "
	EnqueueJobFailure         = "Failed to enqueue job: "
	HashPasswordFailure       = "Failed to hash password: "
	JSONMarshalFailure        = "Failed to marshal json: "
	JSONUnmarshalFailure      = "Failed to unmarshal json: "
	MigrateDatabaseFailure    = "Failed to migrate database: "
	ReadResponseBodyFailure   = "Failed to read response body: "
	RecorderInitializeFailure = "Failed to init recorder: "
	RecordStopFailure         = "Failed to stop record: "
	RequestInitializeFailure  = "Failed to init request: "
	RequestFailure            = "Failed to make a request: "
	ScanRowFailure            = "Failed to scan row: "
	StartOAuthServerFailure   = "Failed to setup oauth server: "
)
