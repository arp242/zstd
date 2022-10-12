module zgo.at/zstd

go 1.19

// Exception to the "stdlib-only" rule, since this should be in stdlib soon, and
// adding some of these methods to zslice only to deprecate them in a few months
// isn't too helpful.
require golang.org/x/exp v0.0.0-20221012134508-3640c57a48ea
