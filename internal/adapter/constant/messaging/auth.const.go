package adaptermessagingconst

import (
	"time"

	"github.com/thanvuc/go-core-lib/eventbus"
)

const (
	HandlerName = "auth_handler"
)

// =================================
// Sync user data
// =================================
const (
	AuthExchangeName eventbus.ExchangeName = "sync_database"
)

const (
	AuthQueueName  = "sync_user_queue_team"
	AuthRoutingKey = "sync.auth.user"
)

const (
	InstanceNumber = 1
	RetryInterval  = 3 * time.Second
)
