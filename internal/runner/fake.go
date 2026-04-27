package runner

import (
	"fmt"
	"reflect"
)

// registration holds a pre-configured result and error for a command name.
type registration struct {
	result Result
	err    error
}

// FakeRunner is a test double that records Run calls and returns pre-configured
// results. It never executes real processes.
type FakeRunner struct {
	registered map[string]registration
	calls      []Command
}

// NewFakeRunner returns a FakeRunner with no registrations.
func NewFakeRunner() *FakeRunner {
	return &FakeRunner{
		registered: make(map[string]registration),
	}
}

// Register pre-configures the result and error to return when Run is called
// with a command whose Name matches name. Each name may be registered once;
// later calls overwrite earlier ones.
func (f *FakeRunner) Register(name string, result Result, err error) {
	f.registered[name] = registration{result: result, err: err}
}

// Run records the command and returns the pre-configured result for cmd.Name.
// If cmd.Name has not been registered, Run returns an error.
func (f *FakeRunner) Run(cmd Command) (Result, error) {
	f.calls = append(f.calls, cmd)
	reg, ok := f.registered[cmd.Name]
	if !ok {
		return Result{}, fmt.Errorf("fake runner: no registration for command %q", cmd.Name)
	}
	return reg.result, reg.err
}

// Calls returns a copy of all commands that have been passed to Run.
func (f *FakeRunner) Calls() []Command {
	out := make([]Command, len(f.calls))
	copy(out, f.calls)
	return out
}

// AssertCalled returns nil if at least one call to Run had Name == name and
// Args deeply equal to args. Otherwise it returns a descriptive error.
func (f *FakeRunner) AssertCalled(name string, args []string) error {
	for _, c := range f.calls {
		if c.Name == name && reflect.DeepEqual(c.Args, args) {
			return nil
		}
	}
	return fmt.Errorf("fake runner: no call to %q with args %v (recorded calls: %v)", name, args, f.calls)
}
