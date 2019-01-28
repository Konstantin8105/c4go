package cupaloy

// Configurator is a functional option that can be passed to cupaloy.New() to change snapshotting behaviour.
type Configurator func(*Config)

// EnvVariableName can be used to customize the environment variable that determines whether snapshots
// should be updated e.g.
//  cupaloy.New(EnvVariableName("UPDATE"))
// Will create an instance where snapshots will be updated if the UPDATE environment variable is set.
// Default: UPDATE_SNAPSHOTS
func EnvVariableName(name string) Configurator {
	return func(c *Config) {
		c.shouldUpdate = func() bool {
			return envVariableSet(name)
		}
	}
}

// ShouldUpdate can be used to provide custom logic to decide whether or not to update a snapshot
// e.g.
//   var update = flag.Bool("update", false, "update snapshots")
//   cupaloy.New(ShouldUpdate(func () bool { return *update })
// Will create an instance where snapshots are updated if the --update flag is passed to go test.
// Default: checks for the presence of the UPDATE_SNAPSHOTS environment variable
func ShouldUpdate(f func() bool) Configurator {
	return func(c *Config) {
		c.shouldUpdate = f
	}
}

// SnapshotSubdirectory can be used to customize the location that snapshots are stored in.
// e.g.
//  cupaloy.New(SnapshotSubdirectory("testdata"))
// Will create an instance where snapshots are stored in the "testdata" folder
// Default: .snapshots
func SnapshotSubdirectory(name string) Configurator {
	return func(c *Config) {
		c.subDirName = name
	}
}

// FailOnUpdate controls whether tests should be failed when snapshots are updated.
// By default this is true to prevent snapshots being accidentally updated in CI.
// Default: true
func FailOnUpdate(failOnUpdate bool) Configurator {
	return func(c *Config) {
		c.failOnUpdate = failOnUpdate
	}
}

// CreateNewAutomatically controls whether snapshots should be automatically created
// if no matching snapshot already exists.
// Default: true
func CreateNewAutomatically(createNewAutomatically bool) Configurator {
	return func(c *Config) {
		c.createNewAutomatically = createNewAutomatically
	}
}

// FatalOnMismatch controls whether failed tests should fail using t.Fatal which should
// immediately stop any remaining tests. Will use t.Error on false.
// Default: false
func FatalOnMismatch(fatalOnMismatch bool) Configurator {
	return func(c *Config) {
		c.fatalOnMismatch = fatalOnMismatch
	}
}

// SnapshotFileExtension allows you to change the extension of the snapshot files
// that are written. E.g. if you're snapshotting HTML then adding SnapshotFileExtension(".html")
// will allow for more easier viewing of snapshots.
// Default: "", no extension is added.
func SnapshotFileExtension(snapshotFileExtension string) Configurator {
	return func(c *Config) {
		c.snapshotFileExtension = snapshotFileExtension
	}
}

// Config provides the same snapshotting functions with additional configuration capabilities.
type Config struct {
	shouldUpdate           func() bool
	subDirName             string
	failOnUpdate           bool
	createNewAutomatically bool
	fatalOnMismatch        bool
	snapshotFileExtension  string
}

// NewDefaultConfig returns a new Config instance initialised with the same options as
// the original Global instance (i.e. before any config changes were made to it)
func NewDefaultConfig() *Config {
	return (&Config{}).WithOptions(
		SnapshotSubdirectory(".snapshots"),
		EnvVariableName("UPDATE_SNAPSHOTS"),
		FailOnUpdate(true),
		CreateNewAutomatically(true),
		FatalOnMismatch(false),
		SnapshotFileExtension(""),
	)
}

// Global is the Config instance used by `cupaloy.SnapshotT` and other package-level functions.
var Global = NewDefaultConfig()

func (c *Config) clone() *Config {
	return &Config{
		shouldUpdate:           c.shouldUpdate,
		subDirName:             c.subDirName,
		failOnUpdate:           c.failOnUpdate,
		createNewAutomatically: c.createNewAutomatically,
		fatalOnMismatch:        c.fatalOnMismatch,
		snapshotFileExtension:  c.snapshotFileExtension,
	}
}
