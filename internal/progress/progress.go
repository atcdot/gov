package progress

import (
	"time"

	"github.com/getlantern/systray"

	"github.com/atcdot/gov/internal/resources"
)

type GUIProgress struct {
	blocked chan interface{}
}

func (p *GUIProgress) Start() {
	if p.blocked != nil {
		close(p.blocked)
		p.blocked = nil
	}
}

func (p *GUIProgress) Stop() {
	if p.blocked == nil {
		p.blocked = make(chan interface{})
	}
}

func (p *GUIProgress) Run() {
	p.blocked = make(chan interface{})

	go func(p *GUIProgress) {
		for {
			if p.blocked != nil {
				<-p.blocked
			}

			systray.SetTemplateIcon(resources.IconV2, resources.IconV2)
			time.Sleep(300 * time.Millisecond)
			systray.SetTemplateIcon(resources.Icon, resources.Icon)
			time.Sleep(300 * time.Millisecond)
		}
	}(p)
}
