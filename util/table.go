package util

import (
	"ctrl/database"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

func getTable(columns []string, rows [][]string) *table.Table {
	t := table.New().
		Border(lipgloss.HiddenBorder()).
		Headers(columns...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return lipgloss.NewStyle().
					Foreground(lipgloss.Color("212")).
					Border(lipgloss.NormalBorder()).
					BorderTop(false).
					BorderLeft(false).
					BorderRight(false).
					BorderBottom(true).
					Bold(true)
			}
			if row%2 == 0 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("246"))
			}
			return lipgloss.NewStyle()
		})
	return t
}

func ListAllTableOS(list []database.ListAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Status", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		status := "Not Installed"
		if app.Installed {
			status = "Installed"
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			status,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func ListAllTable(list []database.ListAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Status", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		status := "Not Installed"
		if app.Installed {
			status = "Installed"
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			app.Os,
			status,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func ListNotInstalledOS(list []database.FindNotInstalledAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func ListNotInstalled(list []database.FindNotInstalledAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			app.Os,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func ListInstalledOS(list []database.FindInstalledAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func ListInstalled(list []database.FindInstalledAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Source", "Version", "Avaiable", "Last Updated"}
	var rows [][]string
	for _, app := range list {
		version := "-"
		if app.Version.Valid {
			version = app.Version.String
		}

		available := "-"
		if app.Available.Valid {
			available = app.Available.String
		}

		rows = append(rows, []string{
			app.ID,
			app.Name,
			app.Os,
			app.Source,
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}
