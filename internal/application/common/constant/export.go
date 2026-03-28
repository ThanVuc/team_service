package appconstant

type SprintTaskStatus string

const (
	SprintTaskStatusNotDoneUnestimated SprintTaskStatus = "NOT_DONE_UNESTIMATED"
	SprintTaskStatusDoneUnestimated    SprintTaskStatus = "DONE_UNESTIMATED"
	SprintTaskStatusNotDone            SprintTaskStatus = "NOT_DONE"
	SprintTaskStatusDone               SprintTaskStatus = "DONE"
	SprintTaskStatusDoneEarly          SprintTaskStatus = "DONE_EARLY"
	SprintTaskStatusDoneOnTime         SprintTaskStatus = "DONE_ON_TIME"
	SprintTaskStatusDoneLate           SprintTaskStatus = "DONE_LATE"
	SprintTaskStatusDoneSpillover      SprintTaskStatus = "DONE_SPILLOVER"
)

const (
	SprintExportSheetName = "Sprint Burndown"

	SprintExportDefaultFileName  = "sprint_burndown.xlsx"
	SprintExportFileNamePrefix   = "sprint_"
	SprintExportUnestimatedLabel = "UNESTIMATED"
	SprintExportContentType      = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"

	SprintExportDateFormat     = "02/01/2006"
	SprintExportDayLabelFormat = "2-Jan"

	SprintExportMaxDurationDays = 30
)
