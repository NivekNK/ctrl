package main

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

func listAllTableOS(list []database.ListAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Status", "Version", "Avaiable"}
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
			app.Status,
			version,
			available,
		})
	}
	return getTable(columns, rows)
}

func listAllTable(list []database.ListAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Status", "Version", "Avaiable"}
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
			app.Status,
			version,
			available,
		})
	}
	return getTable(columns, rows)
}

func listNotInstalledOS(list []database.FindNotInstalledAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Version", "Avaiable", "Last Updated"}
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
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func listNotInstalled(list []database.FindNotInstalledAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Version", "Avaiable", "Last Updated"}
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
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func listInstalledOS(list []database.FindInstalledAppsOSRow) *table.Table {
	columns := []string{"ID", "Name", "Version", "Avaiable", "Last Updated"}
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
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}

func listInstalled(list []database.FindInstalledAppsRow) *table.Table {
	columns := []string{"ID", "Name", "OS", "Version", "Avaiable", "Last Updated"}
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
			version,
			available,
			app.LastUpdated,
		})
	}
	return getTable(columns, rows)
}
