package main

import (
	"adventure-game/database"
	"adventure-game/database/models"
	"fmt"
	"log"
	"regexp"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// install
func init() {
	database.RunMigrations()
}

func main() {
	mainScreen()
}

func mainScreen() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	currentSelected := 0

	searchInputText := ""

	searchInput := widgets.NewParagraph()
	searchInput.Title = "Busca de jogadas"
	searchInput.Text = searchInputText
	searchInput.SetRect(1, 3, 180, 6)
	searchInput.BorderStyle.Fg = ui.ColorYellow

	table := widgets.NewTable()

	table.Title = "Lista de Jogadores"
	table.TextStyle = ui.NewStyle(ui.ColorWhite)
	table.RowSeparator = true
	table.BorderStyle = ui.NewStyle(ui.ColorGreen)
	table.SetRect(1, 7, 180, 30)
	table.FillRow = true

	shortCuts := widgets.NewParagraph()
	shortCuts.Title = "Atalhos"
	shortCuts.Text = `CTRL+c ou ESC - Encerra jogo | Seta para cima ou para baixo - Muda seleção de partida | CTRL+n - Nova partida | ENTER - Seleciona partida`
	shortCuts.SetRect(1, 40, 180, 43)
	shortCuts.BorderStyle.Fg = ui.ColorYellow

	updateTable := func(s string) {
		playerData := searchPlayers(s)
		rows := [][]string{
			[]string{"*", "Id", "Nome", "Parágrafo Atual", "Habilidade", "Sorte", "Energia", "Habilidade Inicial", "Sorte Inicial", "Energia Inicial", "Provisões"},
		}

		for idx, playerDataRow := range playerData {
			var rowSelected []string

			if idx == currentSelected {
				rowSelected = []string{"*"}
			} else {
				rowSelected = []string{""}
			}

			rowSelected = append(rowSelected, playerDataRow...)

			rows = append(rows, rowSelected)
		}

		table.Rows = rows

		table.RowStyles[0] = ui.NewStyle(ui.ColorWhite, ui.ColorBlack, ui.ModifierBold)

		for idx := range rows {
			if idx == 0 {
				continue
			}

			if idx % 2 == 0 {
				table.RowStyles[idx] = ui.NewStyle(ui.ColorYellow)
				continue
			}

			table.RowStyles[idx] = ui.NewStyle(ui.ColorWhite)
		}

		ui.Render(table)
		ui.Render(shortCuts)
	}

	updateInput := func(s string, backspace bool) {
		if backspace {
			if len(searchInputText) > 0 {
				searchInputText = searchInputText[:len(searchInputText)-1]
			}
		} else {
			searchInputText += s
		}

		searchInput.Text = searchInputText
		ui.Render(searchInput)
		updateTable(searchInput.Text)
	}

	updateTable(``)
	updateInput(``, false)

	uiEvents := ui.PollEvents()

	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	for {
		e := <-uiEvents

		switch e.ID {
		case "<C-c>", "<Escape>":
			return
		case "<Resize>":
			ui.Clear()
			updateTable(searchInput.Text)
			ui.Render(searchInput, table)
		case "<Backspace>":
			updateInput(" ", true)
		case "<Space>":
			updateInput(" ", false)
		case "<Up>":
			if currentSelected == 0 {
				continue
			}
			currentSelected--
			updateTable(searchInput.Text)
		case "<Down>":
			if currentSelected >= (len(table.Rows) - 2) {
				continue
			}
			currentSelected++
			updateTable(searchInput.Text)
		default:
			if regex.MatchString(e.ID) {
				currentSelected = 0
				updateInput(e.ID, false)
			}
		}
	}
}

// Função para buscar jogadores pelo nome
func searchPlayers(name string) [][]string {
	players := models.Player{Name: name}.Retrieve()
	return formatPlayerData(players)
}

// Função para formatar os dados dos jogadores para exibição na tabela
func formatPlayerData(players []models.Player) [][]string {
	var data [][]string
	for _, player := range players {
		data = append(data, []string{
			fmt.Sprintf("%d", player.Id),
			player.Name,
			fmt.Sprintf("%d", player.Hability),
			fmt.Sprintf("%d", player.Luck),
			fmt.Sprintf("%d", player.Energy),
			fmt.Sprintf("%d", player.Provisions),
			fmt.Sprintf("%d", player.InitialHability),
			fmt.Sprintf("%d", player.InitialLuck),
			fmt.Sprintf("%d", player.InitialEnergy),
			fmt.Sprintf("%d", player.Paragraph),
		})
	}
	return data
}
