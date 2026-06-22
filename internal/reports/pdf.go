package reports

import (
	"fmt"

	"github.com/go-pdf/fpdf"
)

func GeneratePDF(data *ReportData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	// ── Header ──────────────────────────────────────────────────────────────
	pdf.SetFont("Arial", "B", 20)
	pdf.SetTextColor(30, 80, 160)
	pdf.CellFormat(0, 10, "SITELOGIX", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 6, "Construction Site Management Report", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "Generated: "+data.GeneratedAt.Format("02 Jan 2006 15:04"), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// ── Project Info ─────────────────────────────────────────────────────────
	sectionHeader(pdf, "PROJECT INFORMATION")
	infoRow(pdf, "Project", data.Project.Name)
	infoRow(pdf, "Status", data.Project.Status)
	if data.Project.Location != nil {
		infoRow(pdf, "Location", *data.Project.Location)
	}
	if data.Project.StartDate != nil {
		infoRow(pdf, "Start Date", *data.Project.StartDate)
	}
	if data.Project.EndDate != nil {
		infoRow(pdf, "End Date", *data.Project.EndDate)
	}
	if data.DateFrom != "" || data.DateTo != "" {
		infoRow(pdf, "Report Period", data.DateFrom+" to "+data.DateTo)
	}
	pdf.Ln(4)

	// ── Daily Logs ───────────────────────────────────────────────────────────
	sectionHeader(pdf, fmt.Sprintf("DAILY LOGS (%d)", len(data.Logs)))
	if len(data.Logs) == 0 {
		noDataRow(pdf, "No logs found for this period")
	} else {
		tableHeader(pdf, []string{"Date", "Submitted By", "Weather", "Status", "Notes"}, []float64{28, 35, 40, 22, 65})
		for _, l := range data.Logs {
			weather := ""
			if l.Weather != nil {
				weather = truncate(*l.Weather, 30)
			}
			notes := ""
			if l.Notes != nil {
				notes = truncate(*l.Notes, 45)
			}
			tableRow(pdf, []string{l.Date, l.CreatedBy, weather, l.Status, notes}, []float64{28, 35, 40, 22, 65})
		}
	}
	pdf.Ln(4)

	// ── Issues ────────────────────────────────────────────────────────────────
	sectionHeader(pdf, fmt.Sprintf("ISSUES (%d)", len(data.Issues)))
	if len(data.Issues) == 0 {
		noDataRow(pdf, "No issues recorded")
	} else {
		tableHeader(pdf, []string{"Title", "Priority", "Status", "Reported By", "Date"}, []float64{60, 22, 25, 40, 28})
		for _, i := range data.Issues {
			assigned := "-"
			if i.AssignedTo != nil {
				assigned = *i.AssignedTo
			}
			_ = assigned
			tableRow(pdf, []string{truncate(i.Title, 40), i.Priority, i.Status, i.ReportedBy, i.CreatedAt}, []float64{60, 22, 25, 40, 28})
		}
	}
	pdf.Ln(4)

	// ── Attendance ────────────────────────────────────────────────────────────
	sectionHeader(pdf, fmt.Sprintf("ATTENDANCE RECORDS (%d)", len(data.Attendance)))
	if len(data.Attendance) == 0 {
		noDataRow(pdf, "No attendance records for this period")
	} else {
		tableHeader(pdf, []string{"Worker", "Trade", "Check In", "Check Out"}, []float64{50, 35, 52, 52})
		for _, a := range data.Attendance {
			trade := "-"
			if a.Trade != nil {
				trade = *a.Trade
			}
			checkOut := "-"
			if a.CheckOut != nil {
				checkOut = *a.CheckOut
			}
			tableRow(pdf, []string{a.WorkerName, trade, a.CheckIn, checkOut}, []float64{50, 35, 52, 52})
		}
	}

	// ── Footer ────────────────────────────────────────────────────────────────
	pdf.SetY(-15)
	pdf.SetFont("Arial", "I", 8)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 10, "SiteLogix – Confidential Site Report", "", 0, "C", false, 0, "")

	var buf []byte
	if err := pdf.OutputAndClose(newBytesWriter(&buf)); err != nil {
		return nil, err
	}
	return buf, nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func sectionHeader(pdf *fpdf.Fpdf, title string) {
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(30, 80, 160)
	pdf.SetTextColor(255, 255, 255)
	pdf.CellFormat(0, 8, "  "+title, "", 1, "L", true, 0, "")
	pdf.SetTextColor(0, 0, 0)
	pdf.Ln(1)
}

func infoRow(pdf *fpdf.Fpdf, label, value string) {
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(40, 6, label+":", "", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 6, value, "", 1, "L", false, 0, "")
}

func tableHeader(pdf *fpdf.Fpdf, headers []string, widths []float64) {
	pdf.SetFont("Arial", "B", 8)
	pdf.SetFillColor(220, 230, 245)
	pdf.SetTextColor(30, 80, 160)
	for i, h := range headers {
		pdf.CellFormat(widths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)
	pdf.SetTextColor(0, 0, 0)
}

func tableRow(pdf *fpdf.Fpdf, cells []string, widths []float64) {
	pdf.SetFont("Arial", "", 8)
	pdf.SetFillColor(245, 248, 255)
	for i, cell := range cells {
		pdf.CellFormat(widths[i], 6, cell, "1", 0, "L", false, 0, "")
	}
	pdf.Ln(-1)
}

func noDataRow(pdf *fpdf.Fpdf, msg string) {
	pdf.SetFont("Arial", "I", 9)
	pdf.SetTextColor(150, 150, 150)
	pdf.CellFormat(0, 7, "  "+msg, "1", 1, "L", false, 0, "")
	pdf.SetTextColor(0, 0, 0)
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

// bytesWriter wraps a []byte slice to satisfy io.WriteCloser.
type bytesWriter struct {
	buf *[]byte
}

func newBytesWriter(buf *[]byte) *bytesWriter { return &bytesWriter{buf: buf} }
func (w *bytesWriter) Write(p []byte) (int, error) {
	*w.buf = append(*w.buf, p...)
	return len(p), nil
}
func (w *bytesWriter) Close() error { return nil }
