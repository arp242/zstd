package zsync

import (
	"fmt"
	"os"
	"sync"

	"zgo.at/zstd/zdebug"
)

// LogMutex is like sync.Mutex, but will log a message to stderr on Lock() and
// Unlock().
//
// This can be a simple but effective way to debug locking issues.
type LogMutex struct{ mu sync.Mutex }

func (d *LogMutex) Lock()   { fmt.Fprintln(os.Stderr, "  LOCK", zdebug.Loc(1)); d.mu.Lock() }
func (d *LogMutex) Unlock() { fmt.Fprintln(os.Stderr, "UNLOCK", zdebug.Loc(1)); d.mu.Unlock() }

// LogRWMutex is like sync.RWMutex, but will log a message to stderr on Lock(),
// Unlock(), RLock(), and RUnlock().
//
// This can be a simple but effective way to debug locking issues.
type LogRWMutex struct{ mu sync.RWMutex }

func (d *LogRWMutex) Lock()    { fmt.Fprintln(os.Stderr, "   LOCK", zdebug.Loc(1)); d.mu.Lock() }
func (d *LogRWMutex) Unlock()  { fmt.Fprintln(os.Stderr, " UNLOCK", zdebug.Loc(1)); d.mu.Unlock() }
func (d *LogRWMutex) RLock()   { fmt.Fprintln(os.Stderr, "  RLOCK", zdebug.Loc(1)); d.mu.RLock() }
func (d *LogRWMutex) RUnlock() { fmt.Fprintln(os.Stderr, "RUNLOCK", zdebug.Loc(1)); d.mu.RUnlock() }
