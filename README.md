training-log-analysis
=====================
[![GoDoc](https://godoc.org/github.com/jonfk/training-log-analysis?status.svg)](https://godoc.org/github.com/jonfk/training-log-analysis)
[![Build Status](https://travis-ci.org/jonfk/training-log-analysis.svg)](https://travis-ci.org/jonfk/training-log-analysis)

Programs for validating and analyzing my training log found
at [jonfk/training-log](https://github.com/jonfk/training-log)

##Notes
This project is built using the [gb tool](http://getgb.io/).

```bash
# to build all commands
$ gb build
```

##Modules
- cmd/
  - data-exporter
  - statistics-projector
  - validator
- training-log/
  - common
  - projections
- work-in-progress/