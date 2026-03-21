package appconstant

import (
	"fmt"

	"github.com/thanvuc/go-core-lib/eventbus"
)

// ==========Common==========
const (
	HandlerName          = "auth_handler"
	NOTIFICATION_SERVICE = "notification"

	// Feature
	TEAM = "team"

	// Common
	EXCHANGE    = "exchange"
	ROUTING_KEY = "routing_key"
)

// ==========Exchange==========
var (
	TEAM_NOTIFICATION_EXCHANGE eventbus.ExchangeName = eventbus.ExchangeName(fmt.Sprintf(
		"%s_%s_%s",
		NOTIFICATION_SERVICE,
		TEAM,
		EXCHANGE,
	))
)

// ==========Routing key==========
var (
	AuthRoutingKey                = "sync.auth.user"
	TEAM_NOTIFICATION_ROUTING_KEY = fmt.Sprintf(
		"%s_%s",
		NOTIFICATION_SERVICE,
		TEAM,
	)
)
