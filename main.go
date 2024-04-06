package main

import (
	"adventure-game/database"
	"adventure-game/database/models"
	"adventure-game/utils"
	"fmt"
	"log"
	"regexp"
	"time"
	"math/rand"

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
	var currentSelectedRow []string

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
	shortCuts.Text = `CTRL+c ou ESC - Encerra jogo | Seta para cima ou para baixo - Muda seleção de partida | CTRL+n - Nova partida | CTRL+x - Exclui partida | ENTER - Seleciona partida`
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
				currentSelectedRow = playerDataRow
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
		case "<C-x>":
			playerToDelete := getPlayer(currentSelectedRow[0])
			playerToDelete.Delete()
			ui.Clear()
			updateTable(searchInput.Text)
			ui.Render(searchInput, table)
		case "<Resize>":
			ui.Clear()
			updateTable(searchInput.Text)
			ui.Render(searchInput, table)
		case "<Enter>":
			playerToPlay := getPlayer(currentSelectedRow[0])
			game(playerToPlay)
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
		case "<C-n>":
			newGame()
			return
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

func newGame() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	shortCuts := widgets.NewParagraph()
	shortCuts.Title = "Atalhos"
	shortCuts.Text = `CTRL+c ou ESC - Volta para tela inicial | ENTER - Inicia partida`
	shortCuts.SetRect(1, 40, 180, 43)
	shortCuts.BorderStyle.Fg = ui.ColorYellow

	newGameInputText := ""

	newGameInput := widgets.NewParagraph()
	newGameInput.Title = "Nome"
	newGameInput.Text = newGameInputText
	newGameInput.SetRect(1, 3, 180, 6)
	newGameInput.BorderStyle.Fg = ui.ColorYellow

	ui.Clear()
	ui.Render(newGameInput, shortCuts)
	uiEvents := ui.PollEvents()

	regex := regexp.MustCompile(`^[a-zA-Z0-9]+$`)

	updateInput := func(s string, backspace bool) {
		if backspace {
			if len(newGameInputText) > 0 {
				newGameInputText = newGameInputText[:len(newGameInputText)-1]
			}
		} else {
			newGameInputText += s
		}

		newGameInput.Text = newGameInputText
		ui.Render(newGameInput)
	}

	for {
		e := <-uiEvents

		switch e.ID {
		case "<C-c>", "<Escape>":
			mainScreen()
			return
		case "<Resize>":
			ui.Clear()
			updateInput("", false)
			ui.Render(newGameInput, shortCuts)
		case "<Backspace>":
			updateInput(" ", true)
		case "<Space>":
			updateInput(" ", false)
		case "<Enter>":
			rand.Seed(time.Now().UnixNano())
			
			p := models.Player{
				Name: newGameInputText,
			}

			name := widgets.NewParagraph()
			name.Title = "Nome do personagem"
			name.Text = newGameInputText
			name.SetRect(1, 7, 180, 10)
			name.BorderStyle.Fg = ui.ColorYellow
			ui.Render(name)

			time.Sleep(2 * time.Second)

			dice := rand.Intn(6) + 1

			hability := widgets.NewParagraph()
			hability.Title = "Habilidade inicial (Dado D6 + 6)"
			hability.Text = fmt.Sprintf("%d + 6 = %d", dice, dice + 6)
			hability.SetRect(1, 12, 180, 15)
			hability.BorderStyle.Fg = ui.ColorYellow
			ui.Render(hability)

			p.InitialHability = dice + 6
			p.Hability = dice + 6

			time.Sleep(2 * time.Second)

			dice = rand.Intn(6) + 1

			luck := widgets.NewParagraph()
			luck.Title = "Sorte inicial (Dado D6 + 6)"
			luck.Text = fmt.Sprintf("%d + 6 = %d", dice, dice + 6)
			luck.SetRect(1, 17, 180, 20)
			luck.BorderStyle.Fg = ui.ColorYellow
			ui.Render(luck)

			p.InitialLuck = dice + 6
			p.Luck = dice + 6

			time.Sleep(2 * time.Second)

			dice = rand.Intn(6) + 1

			energy := widgets.NewParagraph()
			energy.Title = "Energia inicial (Dado D6 + 12)"
			energy.Text = fmt.Sprintf("%d + 12 = %d", dice, dice + 12)
			energy.SetRect(1, 22, 180, 25)
			energy.BorderStyle.Fg = ui.ColorYellow
			ui.Render(energy)

			p.InitialEnergy = dice + 12
			p.Energy = dice + 12

			time.Sleep(2 * time.Second)

			p.Paragraph = 1
			p.Provisions = 10

			p.Create()

			game(p)
			return
		default:
			if regex.MatchString(e.ID) {
				updateInput(e.ID, false)
			}
		}
	}
}

func game(p models.Player) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	defer ui.Close()

	shortCuts := widgets.NewParagraph()
	shortCuts.Title = "Atalhos"
	shortCuts.Text = `CTRL+c ou ESC - Volta para tela inicial | ENTER - Seleciona opção escrita no destino`
	shortCuts.SetRect(1, 40, 180, 43)
	shortCuts.BorderStyle.Fg = ui.ColorYellow


	currentParagraph := models.Paragraph{}
	currentParagraph.Id = p.Paragraph

	currentParagraph = currentParagraph.GetOneById(fmt.Sprintf(`%d`, currentParagraph.Id))

	currentParagraphShow := widgets.NewParagraph()
	currentParagraphShow.Title = fmt.Sprintf(`Parágrafo: %d`, p.Paragraph)
	currentParagraphShow.SetRect(1, 1, 150, 35)
	currentParagraphShow.Text = currentParagraph.Context
	currentParagraphShow.BorderStyle.Fg = ui.ColorYellow

	playerData := widgets.NewParagraph()
	playerData.Title = p.Name
	playerData.Text = fmt.Sprintf(`Habilidade: %d
Sorte: %d
Energia: %d
Habilidade INICIAL: %d
Sorte INICIAL: %d
Energia INICIAL: %d`,
		p.Hability,
		p.Luck,
		p.Energy,
		p.InitialHability,
		p.InitialLuck,
		p.InitialEnergy,
	)
	playerData.SetRect(155, 1, 180, 39)
	playerData.BorderStyle.Fg = ui.ColorYellow

	regex := regexp.MustCompile(`^[0-9]+$`)

	destinyParagraphText := ``

	destinyParagraphInput := widgets.NewParagraph()
	destinyParagraphInput.Title = "Destino"
	destinyParagraphInput.Text = destinyParagraphText
	destinyParagraphInput.SetRect(1, 36, 150, 39)
	destinyParagraphInput.BorderStyle.Fg = ui.ColorYellow

	updateScreen := func() {
		ui.Clear()
		ui.Render(playerData, shortCuts, destinyParagraphInput, currentParagraphShow)
	}

	updateScreen()
	uiEvents := ui.PollEvents()

	updateInput := func(s string, backspace bool) {
		if backspace {
			if len(destinyParagraphText) > 0 {
				destinyParagraphText = destinyParagraphText[:len(destinyParagraphText)-1]
			}
		} else {
			destinyParagraphText += s
		}

		destinyParagraphInput.Text = destinyParagraphText
		ui.Render(destinyParagraphInput)
	}

	for {
		e := <-uiEvents

		switch e.ID {
		case "<C-c>", "<Escape>":
			mainScreen()
			return
		case "<Backspace>":
			updateInput(" ", true)
		case "<Resize>":
			updateScreen()
		case "<Enter>":
			if currentParagraph.ValidateSelectedDestiny(destinyParagraphText) {
				currentParagraph = currentParagraph.GetOneById(destinyParagraphText)
				newParagraphId, _ := utils.ParseInt(destinyParagraphText)

				destinyParagraphText = ""
				updateScreen()
				updateInput(``, false)
				
				p.Paragraph = newParagraphId
				p.Update();
			} else {
				destinyParagraphText = ""
				currentParagraph = currentParagraph.GetOneById(destinyParagraphText)
				updateScreen()
			}
		default:
			if regex.MatchString(e.ID) {
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

func getPlayer(id string) models.Player {
	playerId, err := utils.ParseInt(id)
	if err != nil {
		log.Fatal("não foi possível definir a jogada a ser excluída")
		return models.Player{}
	}

	player := models.Player{}
	player.Id = playerId

	return player.RetrieveOneById()
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
