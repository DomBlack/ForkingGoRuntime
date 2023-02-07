package sql

// tracingQueryStart is called when a query starts
func tracingQueryStart(query string)

// tracingQueryEnd is called when a query finishes
func tracingQueryEnd(err error)
