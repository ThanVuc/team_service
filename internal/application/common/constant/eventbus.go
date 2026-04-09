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
	AI_TEAM_EXCHANGE eventbus.ExchangeName = "ai_team_exchange"
)

// ==========Routing key==========
var (
	AuthRoutingKey                = "sync.auth.user"
	TEAM_NOTIFICATION_ROUTING_KEY = fmt.Sprintf(
		"%s_%s",
		NOTIFICATION_SERVICE,
		TEAM,
	)
	AI_TEAM_SPRINT_GENERATION_ROUTING_KEY        = "ai_team_sprint-generation"
	TEAM_AI_SPRINT_GENERATION_RESULT_ROUTING_KEY = "team_ai_sprint-generation-result"
)
