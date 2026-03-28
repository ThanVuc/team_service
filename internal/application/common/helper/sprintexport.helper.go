package apphelper

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	appconstant "team_service/internal/application/common/constant"
	appdto "team_service/internal/application/common/dto"
	errorbase "team_service/internal/domain/common/apperror"
	errdict "team_service/internal/domain/common/apperror/err"
	"team_service/internal/domain/entity"
	t "time"

	"github.com/xuri/excelize/v2"
)

type SprintExportHelper struct{}

func NewSprintExportHelper() *SprintExportHelper {
	return &SprintExportHelper{}
}

func (h *SprintExportHelper) BuildSprintBurndownExcel(input appdto.SprintExportInput) (*appdto.SprintExportOutput, errorbase.AppError) {
	if err := validateSprintExportInput(input); err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	index, err := f.NewSheet(appconstant.SprintExportSheetName)
	if err != nil {
		return nil, wrapSprintExportError("failed to create workbook sheet", err)
	}

	f.SetActiveSheet(index)
	_ = f.DeleteSheet("Sheet1")

	styles, appErr := createSprintStyles(f)
	if appErr != nil {
		return nil, appErr
	}

	configureSprintSheet(f)
	createSprintHeader(f, styles, input.Sprint)
	createSprintStatusLegend(f, styles)

	memberStart, memberEnd := createSprintMemberTable(f, styles, input.Members)
	burndownStart := memberEnd + 2
	taskStart := burndownStart + 2

	taskViews := buildSprintTaskViews(input.Works, input.Members, input.Sprint)
	dates := sprintDates(input.Sprint.StartDate, input.Sprint.EndDate)
	stats := computeSprintStatistic(input.Works, input.Sprint)

	_, taskDataStart, taskEnd, colorCells := createSprintTaskTable(f, styles, taskStart, taskViews, dates)
	fillSprintMemberTotals(f, memberStart, memberEnd, taskDataStart, taskEnd)
	createSprintBurndown(f, styles, burndownStart, dates, stats.TotalEstimated, taskDataStart, taskEnd)
	createSprintMetricSection(f, styles, taskEnd+2, stats)
	dailyHeaderRow, dailyEndRow := createSprintDailyProgressTable(f, styles, taskEnd+2, input.Sprint, dates, stats.TotalEstimated, burndownStart)
	createSprintBurndownChartFromDailyTable(f, dailyHeaderRow, dailyEndRow)

	applySprintFullBorder(f, styles, burndownStart, len(dates), taskEnd)
	applySprintTaskWrapText(f, taskDataStart, taskEnd)
	applySprintTaskCompletionColors(f, styles, colorCells)

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, wrapSprintExportError("failed to encode workbook to bytes", err)
	}

	fileName := strings.TrimSpace(input.FileName)
	if fileName == "" {
		fileName = appconstant.SprintExportDefaultFileName
	}

	return &appdto.SprintExportOutput{
		FileName: fileName,
		Content:  buf.Bytes(),
	}, nil
}

func validateSprintExportInput(input appdto.SprintExportInput) errorbase.AppError {
	if input.Sprint == nil {
		return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("sprint is required"))
	}

	if strings.TrimSpace(input.Sprint.ID) == "" {
		return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("sprint id is required"))
	}

	if strings.TrimSpace(input.Sprint.Name) == "" {
		return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("sprint name is required"))
	}

	if input.Sprint.StartDate.IsZero() || input.Sprint.EndDate.IsZero() {
		return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("sprint start and end dates are required"))
	}

	start := normalizeDateUTC(input.Sprint.StartDate)
	end := normalizeDateUTC(input.Sprint.EndDate)

	if !start.Before(end) {
		return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("sprint end date must be after start date"))
	}

	if end.Sub(start).Hours() > float64(24*appconstant.SprintExportMaxDurationDays) {
		return errorbase.New(
			errdict.ErrSprintExportInvalidInput,
			errorbase.WithDetail(fmt.Sprintf("sprint duration must not exceed %d days", appconstant.SprintExportMaxDurationDays)),
		)
	}

	for _, member := range input.Members {
		if member == nil || strings.TrimSpace(member.ID) == "" {
			return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("member id is required"))
		}
	}

	for _, work := range input.Works {
		if work == nil || strings.TrimSpace(work.ID) == "" {
			return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("work id is required"))
		}

		if strings.TrimSpace(work.Name) == "" {
			return errorbase.New(errdict.ErrSprintExportInvalidInput, errorbase.WithDetail("work name is required"))
		}
	}

	return nil
}

func wrapSprintExportError(detail string, cause error) errorbase.AppError {
	return errorbase.New(
		errdict.ErrSprintExportGenerateFailed,
		errorbase.WithDetail(detail),
		errorbase.WithCause(cause),
	)
}

func createSprintStyles(f *excelize.File) (appdto.SprintExportStyles, errorbase.AppError) {
	border := []excelize.Border{
		{Type: "left", Style: 1, Color: "#000000"},
		{Type: "right", Style: 1, Color: "#000000"},
		{Type: "top", Style: 1, Color: "#000000"},
		{Type: "bottom", Style: 1, Color: "#000000"},
	}

	base, err := f.NewStyle(&excelize.Style{
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create base style", err)
	}

	bold, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create bold style", err)
	}

	boldBorder, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create bold border style", err)
	}

	header, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#D9E1F2"},
			Pattern: 1,
		},
		Border: border,
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create header style", err)
	}

	totalLabel, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Border:    border,
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E2EFDA"},
			Pattern: 1,
		},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create total label style", err)
	}

	totalData, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Border:    border,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#E2EFDA"},
			Pattern: 1,
		},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create total data style", err)
	}

	doneEarly, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Border:    border,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#9DC3E6"},
			Pattern: 1,
		},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create done-early style", err)
	}

	doneOnTime, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Border:    border,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#A9D18E"},
			Pattern: 1,
		},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create done-on-time style", err)
	}

	doneLate, err := f.NewStyle(&excelize.Style{
		Alignment: &excelize.Alignment{Horizontal: "left"},
		Border:    border,
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#F8696B"},
			Pattern: 1,
		},
	})
	if err != nil {
		return appdto.SprintExportStyles{}, wrapSprintExportError("failed to create done-late style", err)
	}

	return appdto.SprintExportStyles{
		Base:       base,
		Bold:       bold,
		BoldBorder: boldBorder,
		Header:     header,
		TotalLabel: totalLabel,
		TotalData:  totalData,
		DoneEarly:  doneEarly,
		DoneOnTime: doneOnTime,
		DoneLate:   doneLate,
	}, nil
}

func configureSprintSheet(f *excelize.File) {
	_ = f.SetColWidth(appconstant.SprintExportSheetName, "A", "A", 50)
	_ = f.SetColWidth(appconstant.SprintExportSheetName, "B", "B", 30)
	_ = f.SetColWidth(appconstant.SprintExportSheetName, "C", "C", 14)
	_ = f.SetColWidth(appconstant.SprintExportSheetName, "D", "E", 14)
	_ = f.SetColWidth(appconstant.SprintExportSheetName, "F", "AZ", 10)
}

func createSprintHeader(f *excelize.File, s appdto.SprintExportStyles, sprint *entity.Sprint) {
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A1", "Sprint name:")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B1", sprint.Name)

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A2", "Start date:")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B2", sprint.StartDate.Format(appconstant.SprintExportDateFormat))

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A3", "End date:")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B3", sprint.EndDate.Format(appconstant.SprintExportDateFormat))

	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A1", "A3", s.Bold)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A1", "B3", s.Base)
}

func createSprintStatusLegend(f *excelize.File, s appdto.SprintExportStyles) {
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "F1", "Status")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "G1", "Color")

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "F2", "Done Early")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "F3", "Done On Time")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "F4", "Done Late")

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "G2", "Blue")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "G3", "Green")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "G4", "Red")

	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "F1", "G1", s.Header)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "F2", "F4", s.Base)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "G2", "G2", s.DoneEarly)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "G3", "G3", s.DoneOnTime)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "G4", "G4", s.DoneLate)
}

func createSprintMemberTable(f *excelize.File, s appdto.SprintExportStyles, members []*entity.User) (int, int) {
	startRow := 6
	headers := []string{"No", "Member", "Actual", "Expected"}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 5)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, cell, h)
	}

	memberNo := 0
	for i, m := range members {
		if m == nil {
			continue
		}
		memberNo++

		row := startRow + i
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(row), memberNo)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(row), strings.TrimSpace(m.Email))
	}

	endRow := startRow + len(members) - 1
	if len(members) == 0 {
		endRow = startRow
	}

	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A5", "D"+strconv.Itoa(endRow), s.Base)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A5", "D5", s.BoldBorder)

	return startRow, endRow
}

func fillSprintMemberTotals(f *excelize.File, memberStartRow, memberEndRow, taskStartRow, taskEndRow int) {
	if memberEndRow < memberStartRow {
		return
	}

	for row := memberStartRow; row <= memberEndRow; row++ {
		actualFormula := fmt.Sprintf(
			"SUMIFS($D$%d:$D$%d,$B$%d:$B$%d,B%d)",
			taskStartRow,
			taskEndRow,
			taskStartRow,
			taskEndRow,
			row,
		)
		expectedFormula := fmt.Sprintf(
			"SUMIFS($E$%d:$E$%d,$B$%d:$B$%d,B%d)",
			taskStartRow,
			taskEndRow,
			taskStartRow,
			taskEndRow,
			row,
		)

		_ = f.SetCellFormula(appconstant.SprintExportSheetName, "C"+strconv.Itoa(row), actualFormula)
		_ = f.SetCellFormula(appconstant.SprintExportSheetName, "D"+strconv.Itoa(row), expectedFormula)
	}
}

func createSprintTaskTable(
	f *excelize.File,
	s appdto.SprintExportStyles,
	startRow int,
	tasks []appdto.SprintTaskView,
	dates []t.Time,
) (int, int, int, []appdto.SprintExportCellColor) {
	headerRow := startRow
	dataStart := headerRow + 1
	endRow := dataStart + len(tasks) - 1
	if len(tasks) == 0 {
		endRow = dataStart
	}

	colorCells := make([]appdto.SprintExportCellColor, 0)

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(headerRow), "Task Name")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(headerRow), "Assignee")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "C"+strconv.Itoa(headerRow), "Completed At")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "D"+strconv.Itoa(headerRow), "Actual")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "E"+strconv.Itoa(headerRow), "Expected")

	for i, day := range dates {
		col, _ := excelize.ColumnNumberToName(6 + i)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, col+strconv.Itoa(headerRow), day.Format(appconstant.SprintExportDayLabelFormat))
	}

	for i, task := range tasks {
		row := dataStart + i
		hasStoryPoint := task.StoryPoint != nil && *task.StoryPoint > 0

		_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(row), task.Name)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(row), task.AssigneeName)

		if task.CompletedAt != nil {
			_ = f.SetCellValue(
				appconstant.SprintExportSheetName,
				"C"+strconv.Itoa(row),
				task.CompletedAt.Format(appconstant.SprintExportDateFormat),
			)
		}

		if !hasStoryPoint {
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "E"+strconv.Itoa(row), appconstant.SprintExportUnestimatedLabel)
		} else {
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "E"+strconv.Itoa(row), *task.StoryPoint)
		}

		if hasStoryPoint && task.CompletedAt != nil {
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "D"+strconv.Itoa(row), *task.StoryPoint)
		}

		if hasStoryPoint {
			sprintEndDate := normalizeDateUTC(dates[len(dates)-1])

			for dayIdx, day := range dates {
				col, _ := excelize.ColumnNumberToName(6 + dayIdx)
				cell := col + strconv.Itoa(row)
				dayDate := normalizeDateUTC(day)

				if task.CompletedAt == nil {
					_ = f.SetCellValue(appconstant.SprintExportSheetName, cell, *task.StoryPoint)
					continue
				}

				completedDate := normalizeDateUTC(*task.CompletedAt)
				if completedDate.After(sprintEndDate) {
					_ = f.SetCellValue(appconstant.SprintExportSheetName, cell, *task.StoryPoint)
					continue
				}

				if dayDate.Before(completedDate) {
					_ = f.SetCellValue(appconstant.SprintExportSheetName, cell, *task.StoryPoint)
					continue
				}

				_ = f.SetCellValue(appconstant.SprintExportSheetName, cell, 0)

				if dayDate.Equal(completedDate) {
					colorCells = append(colorCells, appdto.SprintExportCellColor{Cell: cell, Status: task.Status})
				}
			}
		}
	}

	endCol, _ := excelize.ColumnNumberToName(5 + len(dates))

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"A"+strconv.Itoa(headerRow),
		endCol+strconv.Itoa(endRow),
		s.Base,
	)

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"A"+strconv.Itoa(headerRow),
		endCol+strconv.Itoa(headerRow),
		s.Header,
	)

	if len(tasks) > 0 {
		wrapStyle, _ := f.NewStyle(&excelize.Style{
			Border: []excelize.Border{
				{Type: "left", Style: 1, Color: "#000000"},
				{Type: "right", Style: 1, Color: "#000000"},
				{Type: "top", Style: 1, Color: "#000000"},
				{Type: "bottom", Style: 1, Color: "#000000"},
			},
			Alignment: &excelize.Alignment{Horizontal: "left", WrapText: true},
		})

		_ = f.SetCellStyle(
			appconstant.SprintExportSheetName,
			"A"+strconv.Itoa(dataStart),
			"B"+strconv.Itoa(endRow),
			wrapStyle,
		)
	}

	return headerRow, dataStart, endRow, colorCells
}

func createSprintBurndown(
	f *excelize.File,
	s appdto.SprintExportStyles,
	startRow int,
	dates []t.Time,
	totalEstimated int,
	taskStartRow int,
	taskEndRow int,
) {
	labelRow := startRow
	actualRow := startRow + 1
	expectedDaily := buildSprintExpectedDailyFormula(totalEstimated, len(dates))

	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(labelRow), "Total")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(labelRow), "Expected")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(actualRow), "Actual")
	_ = f.SetCellFormula(appconstant.SprintExportSheetName, "C"+strconv.Itoa(labelRow), "SUM(E"+strconv.Itoa(taskStartRow)+":E"+strconv.Itoa(taskEndRow)+")")
	_ = f.SetCellFormula(appconstant.SprintExportSheetName, "C"+strconv.Itoa(actualRow), "SUM(D"+strconv.Itoa(taskStartRow)+":D"+strconv.Itoa(taskEndRow)+")")

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"A"+strconv.Itoa(labelRow),
		"C"+strconv.Itoa(actualRow),
		s.TotalLabel,
	)

	for i := 0; i < len(dates); i++ {
		col, _ := excelize.ColumnNumberToName(6 + i)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, col+strconv.Itoa(labelRow), expectedDaily[i])
		_ = f.SetCellFormula(
			appconstant.SprintExportSheetName,
			col+strconv.Itoa(actualRow),
			"SUM("+col+strconv.Itoa(taskStartRow)+":"+col+strconv.Itoa(taskEndRow)+")",
		)
	}

	endCol, _ := excelize.ColumnNumberToName(5 + len(dates))

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"C"+strconv.Itoa(labelRow),
		endCol+strconv.Itoa(actualRow),
		s.TotalData,
	)
}

func createSprintBurndownChartFromDailyTable(f *excelize.File, dailyHeaderRow, dailyEndRow int) {
	if dailyEndRow <= dailyHeaderRow {
		return
	}

	anchorCell := "H" + strconv.Itoa(dailyHeaderRow)

	categories := fmt.Sprintf("'%s'!$D$%d:$D$%d", appconstant.SprintExportSheetName, dailyHeaderRow+1, dailyEndRow)
	expectedValues := fmt.Sprintf("'%s'!$E$%d:$E$%d", appconstant.SprintExportSheetName, dailyHeaderRow+1, dailyEndRow)
	actualValues := fmt.Sprintf("'%s'!$F$%d:$F$%d", appconstant.SprintExportSheetName, dailyHeaderRow+1, dailyEndRow)
	expectedName := fmt.Sprintf("'%s'!$E$%d", appconstant.SprintExportSheetName, dailyHeaderRow)
	actualName := fmt.Sprintf("'%s'!$F$%d", appconstant.SprintExportSheetName, dailyHeaderRow)

	_ = f.AddChart(appconstant.SprintExportSheetName, anchorCell, &excelize.Chart{
		Type: excelize.Line,
		Series: []excelize.ChartSeries{
			{
				Name:       expectedName,
				Categories: categories,
				Values:     expectedValues,
			},
			{
				Name:       actualName,
				Categories: categories,
				Values:     actualValues,
			},
		},
		Title: []excelize.RichTextRun{{Text: "Burndown Chart"}},
		Legend: excelize.ChartLegend{
			Position: "bottom",
		},
		Dimension: excelize.ChartDimension{
			Width:  720,
			Height: 320,
		},
	})
}

func buildSprintExpectedDailyFormula(totalEstimated int, totalDays int) []int {
	result := make([]int, totalDays)
	if totalDays == 0 {
		return result
	}

	step := float64(totalEstimated) / float64(totalDays)
	for i := 0; i < totalDays; i++ {
		n := float64(i + 1)
		value := float64(totalEstimated) - step*n
		if value < 0 {
			value = 0
		}
		result[i] = int(math.Round(value))
	}

	return result
}

func createSprintMetricSection(f *excelize.File, s appdto.SprintExportStyles, startRow int, st appdto.SprintStatistic) {
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(startRow), "Metric")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(startRow), "Value")

	rows := []struct {
		name  string
		value any
	}{
		{name: "Total Estimated", value: st.TotalEstimated},
		{name: "Total Completed", value: st.TotalCompleted},
		{name: "Completion Rate", value: fmt.Sprintf("%.2f%%", st.CompletionRate*100)},
		{name: "Spillover", value: st.Spillover},
		{name: "Unestimated Count", value: st.UnestimatedCount},
	}

	for i, item := range rows {
		row := startRow + 1 + i
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "A"+strconv.Itoa(row), item.name)
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "B"+strconv.Itoa(row), item.value)
	}

	endRow := startRow + len(rows)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A"+strconv.Itoa(startRow), "B"+strconv.Itoa(endRow), s.Base)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "A"+strconv.Itoa(startRow), "B"+strconv.Itoa(startRow), s.Header)
}

func createSprintDailyProgressTable(
	f *excelize.File,
	s appdto.SprintExportStyles,
	startRow int,
	sprint *entity.Sprint,
	dates []t.Time,
	totalEstimated int,
	burndownStartRow int,
) (int, int) {
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "D"+strconv.Itoa(startRow), "Date")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "E"+strconv.Itoa(startRow), "Expected")
	_ = f.SetCellValue(appconstant.SprintExportSheetName, "F"+strconv.Itoa(startRow), "Actual")

	rows := len(dates) + 1
	for i := 0; i < rows; i++ {
		row := startRow + 1 + i

		if i == 0 {
			baselineDate := sprint.StartDate.AddDate(0, 0, -1)
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "D"+strconv.Itoa(row), baselineDate.Format(appconstant.SprintExportDayLabelFormat))
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "E"+strconv.Itoa(row), totalEstimated)
			_ = f.SetCellValue(appconstant.SprintExportSheetName, "F"+strconv.Itoa(row), totalEstimated)
			continue
		}

		day := dates[i-1]
		_ = f.SetCellValue(appconstant.SprintExportSheetName, "D"+strconv.Itoa(row), day.Format(appconstant.SprintExportDayLabelFormat))

		burndownCol, _ := excelize.ColumnNumberToName(6 + (i - 1))
		_ = f.SetCellFormula(appconstant.SprintExportSheetName, "E"+strconv.Itoa(row), burndownCol+strconv.Itoa(burndownStartRow))
		_ = f.SetCellFormula(appconstant.SprintExportSheetName, "F"+strconv.Itoa(row), burndownCol+strconv.Itoa(burndownStartRow+1))
	}

	endRow := startRow + rows
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "D"+strconv.Itoa(startRow), "F"+strconv.Itoa(endRow), s.Base)
	_ = f.SetCellStyle(appconstant.SprintExportSheetName, "D"+strconv.Itoa(startRow), "F"+strconv.Itoa(startRow), s.Header)

	return startRow, endRow
}

func buildSprintTaskViews(works []*entity.Work, members []*entity.User, sprint *entity.Sprint) []appdto.SprintTaskView {
	memberMap := map[string]string{}
	for _, m := range members {
		if m == nil {
			continue
		}

		memberMap[m.ID] = strings.TrimSpace(m.Email)
	}

	views := make([]appdto.SprintTaskView, 0, len(works))
	for _, w := range works {
		if w == nil {
			continue
		}

		if strings.TrimSpace(w.SprintID) != "" && w.SprintID != sprint.ID {
			continue
		}

		var storyPoint *int
		if w.StoryPoint != nil && *w.StoryPoint > 0 {
			value := int(*w.StoryPoint)
			storyPoint = &value
		}

		assigneeName := strings.TrimSpace(memberMap[w.AssigneeID])
		if assigneeName == "" {
			assigneeName = "Unassigned"
		}

		views = append(views, appdto.SprintTaskView{
			Name:         strings.TrimSpace(w.Name),
			AssigneeName: assigneeName,
			StoryPoint:   storyPoint,
			CompletedAt:  w.CompletedAt,
			Status:       classifySprintWork(w, sprint),
		})
	}

	sort.SliceStable(views, func(i, j int) bool {
		left := views[i].CompletedAt
		right := views[j].CompletedAt

		if left == nil && right == nil {
			return views[i].Name < views[j].Name
		}
		if left == nil {
			return false
		}
		if right == nil {
			return true
		}
		return left.Before(*right)
	})

	return views
}

func classifySprintWork(w *entity.Work, sprint *entity.Sprint) appconstant.SprintTaskStatus {
	hasStoryPoint := w.StoryPoint != nil && *w.StoryPoint > 0

	if !hasStoryPoint && w.CompletedAt == nil {
		return appconstant.SprintTaskStatusNotDoneUnestimated
	}

	if !hasStoryPoint && w.CompletedAt != nil {
		return appconstant.SprintTaskStatusDoneUnestimated
	}

	if hasStoryPoint && w.CompletedAt == nil {
		return appconstant.SprintTaskStatusNotDone
	}

	completedDate := normalizeDateUTC(*w.CompletedAt)
	sprintEndDate := normalizeDateUTC(sprint.EndDate)
	if completedDate.After(sprintEndDate) {
		return appconstant.SprintTaskStatusDoneSpillover
	}

	if w.DueDate == nil {
		return appconstant.SprintTaskStatusDone
	}

	dueDate := normalizeDateUTC(*w.DueDate)

	if completedDate.Before(dueDate) {
		return appconstant.SprintTaskStatusDoneEarly
	}

	if completedDate.Equal(dueDate) {
		return appconstant.SprintTaskStatusDoneOnTime
	}

	return appconstant.SprintTaskStatusDoneLate
}

func computeSprintStatistic(works []*entity.Work, sprint *entity.Sprint) appdto.SprintStatistic {
	st := appdto.SprintStatistic{}

	for _, w := range works {
		if w == nil {
			continue
		}

		if strings.TrimSpace(w.SprintID) != "" && w.SprintID != sprint.ID {
			continue
		}

		if w.StoryPoint == nil || *w.StoryPoint <= 0 {
			st.UnestimatedCount++
			continue
		}

		sp := int(*w.StoryPoint)
		st.TotalEstimated += sp

		if w.CompletedAt != nil {
			st.TotalCompleted += sp
			if normalizeDateUTC(*w.CompletedAt).After(normalizeDateUTC(sprint.EndDate)) {
				st.Spillover += sp
			}
		}
	}

	if st.TotalEstimated > 0 {
		st.CompletionRate = float64(st.TotalCompleted) / float64(st.TotalEstimated)
	}

	return st
}

func sprintDates(start t.Time, end t.Time) []t.Time {
	if end.Before(start) {
		return []t.Time{}
	}

	days := int(end.Sub(start).Hours()/24) + 1
	result := make([]t.Time, 0, days)

	for i := 0; i < days; i++ {
		result = append(result, start.AddDate(0, 0, i))
	}

	return result
}

func applySprintTaskCompletionColors(f *excelize.File, s appdto.SprintExportStyles, cells []appdto.SprintExportCellColor) {
	for _, item := range cells {
		switch item.Status {
		case appconstant.SprintTaskStatusDoneEarly:
			_ = f.SetCellStyle(appconstant.SprintExportSheetName, item.Cell, item.Cell, s.DoneEarly)
		case appconstant.SprintTaskStatusDoneOnTime:
			_ = f.SetCellStyle(appconstant.SprintExportSheetName, item.Cell, item.Cell, s.DoneOnTime)
		case appconstant.SprintTaskStatusDoneLate:
			_ = f.SetCellStyle(appconstant.SprintExportSheetName, item.Cell, item.Cell, s.DoneLate)
		}
	}
}

func applySprintTaskWrapText(f *excelize.File, taskStartRow, taskEndRow int) {
	if taskEndRow < taskStartRow {
		return
	}

	wrapStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Style: 1, Color: "#000000"},
			{Type: "right", Style: 1, Color: "#000000"},
			{Type: "top", Style: 1, Color: "#000000"},
			{Type: "bottom", Style: 1, Color: "#000000"},
		},
		Alignment: &excelize.Alignment{Horizontal: "left", Vertical: "top", WrapText: true},
	})

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"A"+strconv.Itoa(taskStartRow),
		"B"+strconv.Itoa(taskEndRow),
		wrapStyle,
	)
}

func applySprintFullBorder(f *excelize.File, s appdto.SprintExportStyles, startRow int, days int, endRow int) {
	endCol, _ := excelize.ColumnNumberToName(5 + days)

	_ = f.SetCellStyle(
		appconstant.SprintExportSheetName,
		"A"+strconv.Itoa(startRow),
		endCol+strconv.Itoa(endRow),
		s.Base,
	)

	for row := startRow; row <= endRow; row++ {
		cell := endCol + strconv.Itoa(row)
		_ = f.SetCellStyle(appconstant.SprintExportSheetName, cell, cell, s.Base)
	}
}

func normalizeDateUTC(value t.Time) t.Time {
	year, month, day := value.UTC().Date()
	return t.Date(year, month, day, 0, 0, 0, 0, t.UTC)
}
