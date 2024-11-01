package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type UI struct {
	window      fyne.Window
	outputPath  *widget.Entry
	dropZone    *canvas.Rectangle
	dropText    *widget.Label
	dropIcon    *widget.Icon
	progress    *widget.ProgressBar
	startBtn    *widget.Button
	previewList *widget.List
	csvData     [][]string
	csvHeaders  []string
}

func NewUI() *UI {
	fmt.Println("Creating new UI")
	mainApp := app.New()
	window := mainApp.NewWindow("Test Media Generator")

	ui := &UI{
		window:     window,
		outputPath: widget.NewEntry(),
		progress:   widget.NewProgressBar(),
	}

	ui.setupUI()
	return ui
}
func (u *UI) setupUI() {
	fmt.Println("Setting up UI")
	// Output directory selection
	if cwd, err := os.Getwd(); err == nil {
		u.outputPath.SetText(cwd)
	}

	browseButton := widget.NewButton("Browse", func() {
		dialog.ShowFolderOpen(func(uri fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, u.window)
				return
			}
			if uri == nil {
				return
			}
			u.outputPath.SetText(uri.Path())
		}, u.window)
	})

	outputContainer := container.NewBorder(nil, nil, nil, browseButton, u.outputPath)

	// Create preview list
	u.previewList = widget.NewList(
		func() int { return len(u.csvData) },
		func() fyne.CanvasObject {
			return widget.NewLabel("Template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			label := o.(*widget.Label)
			row := make(map[string]string)
			for j, value := range u.csvData[i] {
				row[u.csvHeaders[j]] = value
			}

			isVideo := strings.ToLower(row["Content Type"]) == "video"
			ext := "mp4"
			if !isVideo {
				ext = "jpg"
			}

			filename := fmt.Sprintf("%s_%s_%s.%s",
				row["Resolution"],
				strings.ToLower(row["Content Type"]),
				row["Duration"],
				ext,
			)
			label.SetText(filename)
		},
	)
	u.previewList.Hide()

	// Create start button
	u.startBtn = widget.NewButton("Start Processing", func() {
		u.startProcessing()
	})
	u.startBtn.Hide()

	// Create drop zone
	u.dropZone = canvas.NewRectangle(theme.BackgroundColor())
	u.dropZone.StrokeWidth = 2
	u.dropZone.StrokeColor = theme.PrimaryColor()

	u.dropIcon = widget.NewIcon(theme.DocumentIcon())
	u.dropText = widget.NewLabel("Drop formats.csv here")
	u.dropText.Alignment = fyne.TextAlignCenter

	dropContent := container.NewCenter(container.NewVBox(
		u.dropIcon,
		u.dropText,
	))

	dropContainer := container.NewMax(u.dropZone, dropContent)
	dropContainer.Resize(fyne.NewSize(300, 200))

	// Progress bar
	u.progress.Hide()

	// Create main content layout with spacer for the list
	headerContent := container.NewVBox(
		widget.NewLabel("Output Directory:"),
		outputContainer,
		dropContainer,
		widget.NewLabel("Files to be generated:"),
	)

	// Create bottom button container
	bottomContent := container.NewVBox(
		u.progress,
		u.startBtn,
	)

	// Main layout using Border container to push button to bottom
	mainContent := container.NewBorder(
		headerContent,                      // top
		bottomContent,                      // bottom
		nil,                                // left
		nil,                                // right
		container.NewPadded(u.previewList), // center (will expand)
	)

	u.window.SetContent(mainContent)
	u.window.Resize(fyne.NewSize(600, 500))

	// Set up window drop handler
	u.window.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		fmt.Printf("Files dropped at position %v\n", pos)
		for i, uri := range uris {
			fmt.Printf("Dropped file %d: %s\n", i+1, uri.String())
		}

		if len(uris) > 0 {
			u.handleDrop(uris[0].Path())
		}
	})
}
func (u *UI) setDropZoneHighlight(highlighted bool) {
	fmt.Printf("Setting highlight: %v\n", highlighted)
	if highlighted {
		u.dropZone.StrokeColor = theme.FocusColor()
		u.dropZone.StrokeWidth = 3
		u.dropText.SetText("Release to load file")
	} else {
		u.dropZone.StrokeColor = theme.PrimaryColor()
		u.dropZone.StrokeWidth = 2
		u.dropText.SetText("Drop formats.csv here")
	}
	u.dropZone.Refresh()
	u.dropText.Refresh()
}

func (u *UI) setDropZoneSuccess() {
	fmt.Println("Setting success state")
	u.dropZone.StrokeColor = theme.SuccessColor()
	u.dropIcon.SetResource(theme.ConfirmIcon())
	u.dropText.SetText("CSV file loaded successfully!")
	u.dropZone.Refresh()
	u.dropIcon.Refresh()
	u.dropText.Refresh()
}

func (u *UI) handleDrop(path string) {
	fmt.Printf("Handling dropped file: %s\n", path)

	if filepath.Ext(path) != ".csv" {
		fmt.Println("Error: Not a CSV file")
		dialog.ShowError(fmt.Errorf("please drop a CSV file"), u.window)
		return
	}

	// Read and validate CSV file
	csvFile, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening CSV file: %v\n", err)
		dialog.ShowError(fmt.Errorf("error opening CSV file: %v", err), u.window)
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	headers, err := reader.Read()
	if err != nil {
		fmt.Printf("Error reading CSV headers: %v\n", err)
		dialog.ShowError(fmt.Errorf("error reading CSV headers: %v", err), u.window)
		return
	}

	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading CSV records: %v\n", err)
		dialog.ShowError(fmt.Errorf("error reading CSV records: %v", err), u.window)
		return
	}

	fmt.Printf("Successfully read CSV with %d records\n", len(records))

	// Store CSV data
	u.csvHeaders = headers
	u.csvData = records

	// Update UI
	u.setDropZoneSuccess()
	u.previewList.Show()
	u.previewList.Refresh()
	u.startBtn.Show()
	fmt.Println("Drop handling complete, UI updated")
}

func (u *UI) startProcessing() {
	outputDir := u.outputPath.Text
	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		dialog.ShowError(fmt.Errorf("output directory does not exist: %s", outputDir), u.window)
		return
	}

	// Download font
	fontPath, err := ensureFontExists(outputDir)
	if err != nil {
		dialog.ShowError(fmt.Errorf("error ensuring font exists: %v", err), u.window)
		return
	}

	// Ensure FFmpeg exists
	err = EnsureFfmpegExists()
	if err != nil {
		dialog.ShowError(fmt.Errorf("error ensuring FFmpeg exists: %v", err), u.window)
		return
	}

	go u.processCSV(fontPath, outputDir)
}

func (u *UI) processCSV(fontPath, outputDir string) {
	u.startBtn.Disable()
	u.progress.Show()
	defer func() {
		u.startBtn.Enable()
		u.progress.Hide()
	}()

	u.progress.Min = 0
	u.progress.Max = float64(len(u.csvData))

	for i, record := range u.csvData {
		row := make(map[string]string)
		for j, value := range record {
			row[u.csvHeaders[j]] = value
		}

		err := createTestMedia(row, fontPath, outputDir)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error creating media: %v", err), u.window)
			return
		}

		u.progress.SetValue(float64(i + 1))
	}

	dialog.ShowInformation("Success", "Processing complete!", u.window)
}

func (u *UI) Show() {
	fmt.Println("Showing UI window")
	u.window.ShowAndRun()
}
