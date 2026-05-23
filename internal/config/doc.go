// Package config handles loading and validation of watchlog configuration files.
//
// Configuration is expressed as a JSON file with the following top-level fields:
//
//	{
//	  "paths":           ["/var/log/app.log"],   // required
//	  "level_field":     "level",                // default: "level"
//	  "message_field":   "message",              // default: "message"
//	  "timestamp_field": "time",                 // default: "time"
//	  "min_level":       "warn",                 // optional
//	  "no_color":        false,                  // optional
//	  "filters": [                               // optional
//	    {"field": "service", "pattern": "api"}
//	  ]
//	}
//
// Use Load to read and validate a config file from disk. Defaults are applied
// automatically for any omitted optional fields before validation.
package config
