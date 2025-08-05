package server

const (
	BeforeInitApplication Stage = "BeforeInitApplication"
	AfterInitApplication  Stage = "AfterInitApplication"
	// BeforeInitEtaModules  Stage = "BeforeInitEtaModules"
	// AfterInitEtaModules   Stage = "AfterInitEtaModules"
	BeforeInitMiddleware Stage = "BeforeInitMiddleware"
	AfterInitMiddleware  Stage = "AfterInitMiddleware"
	BeforeInitRouter     Stage = "BeforeInitRouter"
	AfterInitRouter      Stage = "AfterInitRouter"
	BeforeMigrate        Stage = "BeforeMigrate"
	AfterMigrate         Stage = "AfterMigrate"
	BeforeStartServer    Stage = "BeforeStartServer"
	AfterStartServer     Stage = "AfterStartServer"
	BeforeShutdown       Stage = "BeforeShutdown"
	AfterShutdown        Stage = "AfterShutdown"
)

var (
	callbackMap = make(map[Stage][]Callback)
)

type Stage string
type Callback func() error

func AddCallback(stage Stage, fc ...Callback) {
	callbackMap[stage] = append(callbackMap[stage], fc...)
}

func callStageCallback(stage Stage) error {
	for _, fc := range callbackMap[stage] {
		if err := fc(); err != nil {
			return err
		}
	}
	return nil
}
