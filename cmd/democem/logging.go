package democem

import (
	"fmt"
	"time"
)

// Logging interface

func (d *DemoCem) log(level string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s", t.Format(time.RFC3339), level, fmt.Sprintln(args...))
}

func (d *DemoCem) logf(level, format string, args ...interface{}) {
	t := time.Now()
	fmt.Printf("%s: %s %s\n", t.Format(time.RFC3339), level, fmt.Sprintf(format, args...))
}

func (d *DemoCem) Trace(args ...interface{}) {
	d.log("TRACE", args...)
}

func (d *DemoCem) Tracef(format string, args ...interface{}) {
	d.logf("TRACE", format, args...)
}

func (d *DemoCem) Debug(args ...interface{}) {
	d.log("DEBUG", args...)
}

func (d *DemoCem) Debugf(format string, args ...interface{}) {
	d.logf("DEBUG", format, args...)
}

func (d *DemoCem) Info(args ...interface{}) {
	d.log("INFO", args...)
}

func (d *DemoCem) Infof(format string, args ...interface{}) {
	d.logf("INFO", format, args...)
}

func (d *DemoCem) Error(args ...interface{}) {
	d.log("ERROR", args...)
}

func (d *DemoCem) Errorf(format string, args ...interface{}) {
	d.logf("ERROR", format, args...)
}
