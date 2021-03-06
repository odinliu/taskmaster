package taskmaster

import "time"

type SuperFunc func()

type Loggable interface {
	Printf(format string, a ...interface{})
}

// SuperOption is used to init supervisor
type SuperOption struct {
	// when goruntinue panic, if need restart
	// default: false
	NeedRestart bool
	// if goruntinue panic, after `RestartDelay` to restart
	// default: 60 sec
	RestartDelay time.Duration
	// max restart num, 0 for unlimited.
	// default: 3
	MaxFailureTime int
	// max goruntine num
	// defualt: 1
	MaxWorkerNum int
	// a logger should implement `Loggable`
	// silentLogger for default
	Logger Loggable
	// if `IsDebug` is true, supervisor won't recovery from panics.
	// for debug, if you want to see where panic happens
	IsDebug bool
}

type silentLogger struct{}

func (silentLogger) Printf(format string, a ...interface{}) {}

var defaultOption = SuperOption{
	NeedRestart:    false,
	RestartDelay:   60 * time.Second,
	MaxFailureTime: 3,
	MaxWorkerNum:   1,
	Logger:         &silentLogger{},
	IsDebug:        false,
}

type Supervisor struct {
	runnable   SuperFunc
	options    SuperOption
	running    bool
	ch         chan int
	panicTimes int
}

// create a supervisor with defaut options
func DefaultSupervisor(runnable SuperFunc) *Supervisor {
	return NewSupervisor(runnable, defaultOption)
}

// create a supervisor with  customized option
func NewSupervisor(runnable SuperFunc, options SuperOption) *Supervisor {
	return &Supervisor{
		runnable:   runnable,
		options:    options,
		running:    false,
		ch:         nil,
		panicTimes: 0,
	}
}

// call `Start()` to run goruntine with specific options
func (s *Supervisor) Start() {
	s.running = true
	s.ch = make(chan int, s.options.MaxWorkerNum)
	go func() {
		for i := 0; i < s.options.MaxWorkerNum; i++ {
			if s.options.IsDebug {
				go s.runnable()
			} else {
				go s.runWithRecover()
			}
		}
		for s.running {
			select {
			case <-s.ch:
				s.options.Logger.Printf("receive panics...\n")
				s.panicTimes++
				if s.options.NeedRestart {
					go func() {
						select {
						case <-time.After(s.options.RestartDelay):
							s.options.Logger.Printf("starting new work...\n")
							if s.options.IsDebug {
								go s.runnable()
							} else {
								go s.runWithRecover()
							}
						}
					}()
				}
			case <-time.After(60 * time.Second):
				s.options.Logger.Printf("total panic times[%d]\n", s.panicTimes)
			}
		}
	}()
}

// if you don't need supervisor,
// just call `Stop()` to make sure supervisor exit its background goruntine
func (s *Supervisor) Stop() {
	s.running = false
}

func (s *Supervisor) runWithRecover() {
	RecoverableFunc(s.runnable, func(r interface{}) {
		s.options.Logger.Printf(
			"panic[%v], recovering...",
			r,
		)
		s.ch <- 1
	})()
}
