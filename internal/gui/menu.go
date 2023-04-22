package gui

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"

	"gov/internal/dlsdk"
	"gov/internal/gov"
	"gov/internal/progress"
	"gov/internal/resources"
	"gov/internal/version"
)

func RunGUI() {
	onExit := func() {}
	systray.Run(onReadyRef, onExit)
}

type menu struct {
	versionsAvailableMII map[string]*versionMenuItemAvailable
	versionsInstalledMII map[string]*versionMenuItemInstalled
	minorVersionsMII     map[string]*versionMenuItemMinor
	versionsInstalledMI  *systray.MenuItem
	versionsAvailableMI  *systray.MenuItem
}

type versionMenuItemInstalled struct {
	version  string
	menuItem *systray.MenuItem
	bin      string
	isMain   bool
}

type versionMenuItemAvailable struct {
	version  string
	menuItem *systray.MenuItem
}

type versionMenuItemMinor struct {
	version  string
	menuItem *systray.MenuItem
}

func onReadyRef() {
	defer systray.Quit()

	systray.SetTemplateIcon(resources.Icon, resources.Icon)
	systray.SetTitle("GoV")
	systray.SetTooltip("go Version changer")

	m := menu{
		versionsAvailableMII: map[string]*versionMenuItemAvailable{},
		versionsInstalledMII: map[string]*versionMenuItemInstalled{},
		minorVersionsMII:     map[string]*versionMenuItemMinor{},
	}

	m.versionsInstalledMI = systray.AddMenuItem("installed versions", "")
	m.versionsAvailableMI = systray.AddMenuItem("available versions", "")

	err := m.getVersionsInstalled()
	if err != nil {
		return
	}

	err = m.addVersionsAvailable()
	if err != nil {
		return
	}

	openUrl := systray.AddMenuItem("Open https://go.dev/dl/", "available versions")
	go func() {
		for range openUrl.ClickedCh {
			_ = open.Run("https://go.dev/dl/")
		}
	}()

	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "")
	mQuit.SetTemplateIcon(resources.Icon, resources.Icon)

	<-mQuit.ClickedCh
}

func (m *menu) getVersionsInstalled() error {
	vv, err := gov.ListInstalled()
	if err != nil {
		return err
	}

	for _, v := range vv {
		m.versionsInstalledMII[v.Version] = &versionMenuItemInstalled{
			version: v.Version,
			bin:     v.Bin,
		}

		if v.IsMain {
			m.versionsInstalledMII[v.Version].isMain = true
			systray.SetTitle(fmt.Sprintf("GoV %s", v.Version))
		}

		m.addInstalledMenuItem(v.Version)
	}

	return nil
}

func (m *menu) addInstalledMenuItem(version string) {
	mi, ok := m.versionsInstalledMII[version]
	if !ok {
		return
	}

	vMenuItem := m.versionsInstalledMI.AddSubMenuItem(version, "")

	mi.menuItem = vMenuItem

	vMenuItemUse := vMenuItem.AddSubMenuItem("use", "")
	vMenuItemDelete := vMenuItem.AddSubMenuItem("delete", "")
	if mi.isMain {
		vMenuItemDelete.Hide()
	}

	go func(mi *versionMenuItemInstalled, use <-chan struct{}, delete <-chan struct{}) {
		for {
			select {
			case <-use:
				mi.use()
			case <-delete:
				_ = mi.delete()
				m.markAsUninstalled(mi.version)
				return
			}
		}
	}(mi, vMenuItemUse.ClickedCh, vMenuItemDelete.ClickedCh)
}

func (m *menu) markAsUninstalled(version string) {
	miA, ok := m.versionsAvailableMII[version]
	if !ok {
		return
	}
	miA.menuItem.SetTitle(miA.version)
	miA.menuItem.Enable()
}

func (m *menu) addVersionsAvailable() error {
	vva := dlsdk.VersionsAvailable()
	for _, release := range vva {
		m.addVersionAvailable(version.ExtractFull(release.Version))
	}

	return nil
}

func (m *menu) addVersionAvailable(version string) {
	m.versionsAvailableMII[version] = &versionMenuItemAvailable{
		version: version,
	}

	m.addAvailableMenuItem(version)
}

func (m *menu) addAvailableMenuItem(v string) {
	mi, ok := m.versionsAvailableMII[v]
	if !ok {
		return
	}

	mv := version.ExtractMinor(v)
	minorMI, ok := m.minorVersionsMII[mv]
	if !ok {
		minorMI = &versionMenuItemMinor{
			version:  mv,
			menuItem: m.versionsAvailableMI.AddSubMenuItem(mv, ""),
		}
		m.minorVersionsMII[mv] = minorMI
	}

	mi.menuItem = minorMI.menuItem.AddSubMenuItem(mi.version, "")
	if _, ok := m.versionsInstalledMII[mi.version]; ok {
		mi.markAsInstalled()
	}

	go func(mi *versionMenuItemAvailable) {
		for range mi.menuItem.ClickedCh {
			bin, err := installVersion(v)
			if err != nil {
				return
			}

			m.versionsInstalledMII[mi.version] = &versionMenuItemInstalled{
				version: mi.version,
				bin:     bin,
			}

			m.addInstalledMenuItem(mi.version)
			mi.markAsInstalled()
		}
	}(mi)
}

func (v *versionMenuItemInstalled) use() {
	gov.Use(v.version)
	systray.SetTitle(fmt.Sprintf("GoV %s", v.version))
}

func (v *versionMenuItemInstalled) delete() error {
	gov.Remove(v.version)
	v.menuItem.Hide()

	return nil
}

func (v *versionMenuItemAvailable) markAsInstalled() {
	v.menuItem.SetTitle(fmt.Sprintf("%s (installed)", v.version))
	v.menuItem.Disable()
}

func installVersion(version string) (string, error) {
	p := progress.GUIProgress{}
	p.Run()
	p.Start()
	defer p.Stop()

	return gov.Install(version)
}
