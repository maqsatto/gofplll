package gofplll

import "time"

// ClientOption configures a Client.
type ClientOption interface {
	applyClient(*clientConfig)
}

type clientConfig struct {
	binaryPath string
	workDir    string
}

type withBinaryPath struct{ path string }

func (w withBinaryPath) applyClient(c *clientConfig) { c.binaryPath = w.path }

// WithBinaryPath sets the path to the fplll binary.
func WithBinaryPath(path string) ClientOption { return withBinaryPath{path: path} }

type withWorkDir struct{ dir string }

func (w withWorkDir) applyClient(c *clientConfig) { c.workDir = w.dir }

// WithWorkDir sets the working directory for the subprocess.
func WithWorkDir(dir string) ClientOption { return withWorkDir{dir: dir} }

// ReduceOption configures a single Reduce call.
type ReduceOption interface {
	applyReduce(*reduceConfig)
}

type reduceConfig struct {
	algorithm    Algorithm
	delta        string
	eta          string
	blockSize    int
	maxLoops     int
	maxTime      time.Duration
	autoAbort    bool
	noLLL        bool
	floatType    FloatType
	precision    int
	outputFormat OutputFormat
	verbose      bool
	timeout      time.Duration
}

type withAlgorithm struct{ v Algorithm }

func (w withAlgorithm) applyReduce(c *reduceConfig) { c.algorithm = w.v }

// WithAlgorithm sets the reduction algorithm.
func WithAlgorithm(a Algorithm) ReduceOption { return withAlgorithm{v: a} }

type withDelta struct{ v string }

func (w withDelta) applyReduce(c *reduceConfig) { c.delta = w.v }

// WithDelta sets the LLL delta parameter.
func WithDelta(d string) ReduceOption { return withDelta{v: d} }

type withEta struct{ v string }

func (w withEta) applyReduce(c *reduceConfig) { c.eta = w.v }

// WithEta sets the LLL eta parameter.
func WithEta(e string) ReduceOption { return withEta{v: e} }

type withBlockSize struct{ v int }

func (w withBlockSize) applyReduce(c *reduceConfig) { c.blockSize = w.v }

// WithBlockSize sets the BKZ block size.
func WithBlockSize(n int) ReduceOption { return withBlockSize{v: n} }

type withMaxLoops struct{ v int }

func (w withMaxLoops) applyReduce(c *reduceConfig) { c.maxLoops = w.v }

// WithMaxLoops sets the BKZ max loop count.
func WithMaxLoops(n int) ReduceOption { return withMaxLoops{v: n} }

type withMaxTime struct{ v time.Duration }

func (w withMaxTime) applyReduce(c *reduceConfig) { c.maxTime = w.v }

// WithMaxTime sets the BKZ max time limit.
func WithMaxTime(d time.Duration) ReduceOption { return withMaxTime{v: d} }

type withAutoAbort struct{}

func (w withAutoAbort) applyReduce(c *reduceConfig) { c.autoAbort = true }

// WithAutoAbort enables BKZ auto-abort.
func WithAutoAbort() ReduceOption { return withAutoAbort{} }

type withNoLLL struct{}

func (w withNoLLL) applyReduce(c *reduceConfig) { c.noLLL = true }

// WithNoLLL disables LLL pre-processing.
func WithNoLLL() ReduceOption { return withNoLLL{} }

type withFloatType struct{ v FloatType }

func (w withFloatType) applyReduce(c *reduceConfig) { c.floatType = w.v }

// WithFloatType sets the floating-point backend.
func WithFloatType(f FloatType) ReduceOption { return withFloatType{v: f} }

type withPrecision struct{ v int }

func (w withPrecision) applyReduce(c *reduceConfig) { c.precision = w.v }

// WithPrecision sets the MPFR precision in bits.
func WithPrecision(p int) ReduceOption { return withPrecision{v: p} }

type withOutputFormat struct{ v OutputFormat }

func (w withOutputFormat) applyReduce(c *reduceConfig) { c.outputFormat = w.v }

// WithOutputFormat sets what fplll returns.
func WithOutputFormat(f OutputFormat) ReduceOption { return withOutputFormat{v: f} }

type withVerbose struct{}

func (w withVerbose) applyReduce(c *reduceConfig) { c.verbose = true }

// WithVerbose enables verbose fplll output.
func WithVerbose() ReduceOption { return withVerbose{} }

type withTimeout struct{ v time.Duration }

func (w withTimeout) applyReduce(c *reduceConfig) { c.timeout = w.v }

// WithTimeout sets the subprocess timeout.
func WithTimeout(d time.Duration) ReduceOption { return withTimeout{v: d} }

func applyReduceOptions(opts []ReduceOption) reduceConfig {
	var cfg reduceConfig
	for _, o := range opts {
		o.applyReduce(&cfg)
	}
	return cfg
}
