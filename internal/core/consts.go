package core

import "time"

// If the suitable media manifest is not found within this time, the request times out
const TimeoutValue = time.Second * 30
