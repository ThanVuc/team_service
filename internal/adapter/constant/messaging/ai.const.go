package adaptermessagingconst

import (
	"time"

	"github.com/thanvuc/go-core-lib/eventbus"
)

const (
	AISprintGenerationHandlerName = "ai_sprint_generation_handler"
)

const (
	AISprintGenerationExchangeName eventbus.ExchangeName = "team_ai_exchange"
)

const (
	AISprintGenerationQueueName  = "team_ai_sprint-generation-result_queue"
	AISprintGenerationRoutingKey = "team_ai_sprint-generation-result"
)

const (
	AISprintGenerationInstanceNumber = 1
	AISprintGenerationRetryInterval  = 3 * time.Second
)
