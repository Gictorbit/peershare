package client

import (
	"fmt"
	"github.com/gictorbit/peershare/api"
	"github.com/jedib0t/go-pretty/v6/table"
	"os"
)

func (pc *PeerClient) PrintFileInfo(file *api.File) {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"file", "info"})
	t.AppendRows([]table.Row{
		{"Name", file.Name},
		{"Extension", file.Extension},
		{"Size", fmt.Sprintf("%d Bytes", file.Size)},
		{"CheckSum", file.Md5Sum},
	})
	t.AppendSeparator()
	t.Render()
}

func (pc *PeerClient) PrintCode(code string) {
	t := table.NewWriter()
	t.SetStyle(table.StyleLight)
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"file share code"})
	t.AppendRows([]table.Row{
		{code},
	})
	t.AppendSeparator()
	t.Render()
}
