package ui

import (
	"fmt"
	"github.com/JoelOtter/terraform-workspace-i/internal/terraform"
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
	"io"
	"os"
	"strings"
)

type ui struct {
	screen          tcell.Screen
	workspaces      []terraform.Workspace
	modulePath      string
	pointer         int
	deleteWorkspace string
	quit            chan struct{}
}

func (u *ui) drawStr(x int, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		u.screen.SetContent(x, y, c, comb, style)
		x += w
	}
}

func (u *ui) draw() {
	u.screen.Clear()
	u.drawStr(1, 1, tcell.StyleDefault.Bold(true), u.modulePath)
	for i, workspace := range u.workspaces {
		if workspace.Current {
			u.screen.SetCell(1, i+3, tcell.StyleDefault, '*')
		}
		style := tcell.StyleDefault
		if workspace.Current {
			style = style.Bold(true)
		}
		if i == u.pointer {
			style = style.Reverse(true)
		}
		u.drawStr(3, i+3, style, workspace.Name)
	}
	if u.deleteWorkspace != "" {
		w, h := u.screen.Size()
		for i := 1; i < w-1; i++ {
			u.screen.SetCell(i, h-2, tcell.StyleDefault.Background(tcell.ColorRed))
		}
		u.drawStr(
			2,
			h-2,
			tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorBlack),
			fmt.Sprintf("Delete workspace %s (y/n)? ", u.deleteWorkspace),
		)
	}
	u.screen.Show()
}

func (u *ui) keyDown() {
	u.pointer = (u.pointer + 1) % len(u.workspaces)
	u.draw()
}

func (u *ui) keyUp() {
	u.pointer = u.pointer - 1
	if u.pointer < 0 {
		u.pointer = len(u.workspaces) - 1
	}
	u.draw()
}

func (u *ui) run(uiOut io.Writer, uiErr *error) {
	defer close(u.quit)
	for {
		ev := u.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				return
			case tcell.KeyEnter:
				*uiErr = terraform.ChangeWorkspace(u.workspaces[u.pointer].Name, uiOut)
				return
			case tcell.KeyUp, tcell.KeyPgUp, tcell.KeyCtrlP:
				u.keyUp()
			case tcell.KeyDown, tcell.KeyPgDn, tcell.KeyCtrlN:
				u.keyDown()
			case tcell.KeyDelete, tcell.KeyBackspace, tcell.KeyDEL:
				u.deleteWorkspace = u.workspaces[u.pointer].Name
				u.draw()
			case tcell.KeyRune:
				switch ev.Rune() {
				case 'j':
					u.keyDown()
				case 'k':
					u.keyUp()
				case 'y':
					if u.deleteWorkspace != "" {
						u.workspaces, *uiErr = terraform.DeleteWorkspace(u.deleteWorkspace, uiOut)
						if *uiErr != nil {
							return
						}
						u.deleteWorkspace = ""
						u.pointer = u.pointer - 1
						if u.pointer < 0 {
							u.pointer = 0
						}
						u.draw()
					}
				case 'n':
					if u.deleteWorkspace != "" {
						u.deleteWorkspace = ""
					}
					u.draw()
				case 'd':
					u.deleteWorkspace = u.workspaces[u.pointer].Name
					u.draw()
				}
			}
		case *tcell.EventResize:
			u.screen.Sync()
			u.draw()
		}
	}
}

func ShowUI(workspaces []terraform.Workspace) error {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to get screen: %w", err)
	}
	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to init screen: %w", err)
	}

	var uiErr error
	uiOut := &strings.Builder{}
	defer func() {
		if uiOut.Len() > 0 {
			fmt.Print(uiOut.String())
		}
	}()

	moduleRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get module root: %w", err)
	}

	u := &ui{
		screen:          screen,
		workspaces:      workspaces,
		modulePath:      moduleRoot,
		pointer:         getInitialPointer(workspaces),
		deleteWorkspace: "",
		quit:            make(chan struct{}),
	}
	u.draw()

	defer screen.Fini()

	go u.run(uiOut, &uiErr)

	for {
		select {
		case <-u.quit:
			return uiErr
		}
	}
}

func getInitialPointer(workspaces []terraform.Workspace) int {
	for i, workspace := range workspaces {
		if workspace.Current {
			return i
		}
	}
	return 0
}
