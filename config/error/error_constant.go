package error

const (
	ChangeToRootDirFailure     = "Failed to change to root directory: "
	ConnectToDatabaseFailure   = "Failed to connect database: "
	DeleteRedisJobFailure      = "Failed to delete redis job: "
	EnqueueJobFailure          = "Failed to enqueue job: "
	HashPasswordFailure        = "Failed to hash password: "
	JSONMarshalFailure         = "Failed to marshal json: "
	JSONUnmarshalFailure       = "Failed to unmarshal json: "
	MigrateDatabaseFailure     = "Failed to migrate database: "
	PerformJobFailure          = "Failed to perform job: "
	ReadResponseBodyFailure    = "Failed to read response body: "
	RecorderInitializeFailure  = "Failed to init recorder: "
	RecorderStopFailure        = "Failed to stop recorder: "
	RequestInitializeFailure   = "Failed to init request: "
	RequestFailure             = "Failed to make a request: "
	SaveSessionFailure         = "Failed to save session: "
	ScanRowFailure             = "Failed to scan row: "
	StartOAuthServerFailure    = "Failed to setup oauth server: "
	UpdateKeywordStatusFailure = "Failed to update keyword status: "
)
