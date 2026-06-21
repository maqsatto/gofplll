package gofplll

// Algorithm specifies the lattice reduction algorithm.
type Algorithm string

const (
	AlgLLL Algorithm = "lll"
	AlgBKZ Algorithm = "bkz"
	AlgHKZ Algorithm = "hkz"
	AlgSVP Algorithm = "svp"
	AlgCVP Algorithm = "cvp"
)

// FloatType specifies the floating-point backend.
type FloatType string

const (
	FloatMPFR       FloatType = "mpfr"
	FloatDouble     FloatType = "double"
	FloatLongDouble FloatType = "longdouble"
	FloatDD         FloatType = "dd"
	FloatQD         FloatType = "qd"
)

// OutputFormat specifies what fplll outputs.
type OutputFormat string

const (
	OutputBasis  OutputFormat = "b"
	OutputStatus OutputFormat = "t"
	OutputSVP    OutputFormat = "s"
	OutputCVP    OutputFormat = "c"
)
