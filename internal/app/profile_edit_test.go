package app

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/profiles"
	"github.com/ColinEge/8bd-ultimate-api/test/fixtures/keyboard"
	"github.com/ColinEge/8bd-ultimate-appgithub.com/ColinEge/8bitdo-ultimate-app/internal/service"
)

func editFixture(t *testing.T) (*App, *keyboard.KeyboardFixture) {
	t.Helper()
	a := profileApp(t)
	f := &keyboard.KeyboardFixture{}
	if err := f.Reset(); err != nil {
		t.Fatal(err)
	}
	a.service = service.NewService(f.Keyboard)
	a.selectedDevice = f.Device.ID
	t.Cleanup(func() {
		if err := f.AssertComplete(); err != nil {
			t.Error(err)
		}
		if err := a.service.Close(); err != nil {
			t.Error(err)
		}
	})
	return a, f
}

func editPacket(prefix ...byte) []byte {
	p := make([]byte, 33)
	copy(p, prefix)
	return p
}

func titleExchange(title string) keyboard.Exchange {
	units := utf16.Encode([]rune(title))
	p := editPacket(0x54, 0x80, byte(len(units)*2), 0)
	for i, unit := range units {
		p[4+2*i], p[5+2*i] = byte(unit>>8), byte(unit)
	}
	return keyboard.Exchange{Write: editPacket(0x52, 0x80), Response: p}
}

func TestActiveAndInactiveProfileEdits(t *testing.T) {
	for _, active := range []bool{false, true} {
		t.Run(map[bool]string{false: "inactive", true: "active"}[active], func(t *testing.T) {
			a, f := editFixture(t)
			if err := a.CreateProfile("Work"); err != nil {
				t.Fatal(err)
			}
			name := "Other"
			if active {
				name = "Work"
			}
			f.Exchanges = []keyboard.Exchange{titleExchange(name)}
			if active {
				f.Exchanges = append(f.Exchanges,
					keyboard.Exchange{Write: editPacket(0x52, 0xfa, 3, 0x0c, 0, 0xaa, 9, 0x71, 4, 7, 0, 5), Response: []byte{0x54, 0xe4, 8}},
					keyboard.Exchange{Write: editPacket(0x52, 0x83, 4), Response: editPacket(0x54, 0x83, 4, 7, 0, 5)})
			}
			layout, err := a.UpdateProfileBinding("Work", 4, "07-00-05", false)
			if err != nil || layout["0x04"].ActionCode != "07-00-05" {
				t.Fatalf("layout=%v, error=%v", layout, err)
			}
			list, err := a.ListProfiles()
			if err != nil || len(list[0].Keys) != 1 || list[0].Keys[0].Action != "07-00-05" {
				t.Fatalf("stored=%v, %v", list, err)
			}
			f.Exchanges = append(f.Exchanges, titleExchange(name))
			if active {
				f.Exchanges = append(f.Exchanges,
					keyboard.Exchange{Write: editPacket(0x52, 0xfa, 3, 0x0c, 0, 0xaa, 9, 0x71, 4, 7, 0, 4), Response: []byte{0x54, 0xe4, 8}},
					keyboard.Exchange{Write: editPacket(0x52, 0x83, 4), Response: editPacket(0x54, 0x83, 4, 0, 0, 0)})
			}
			layout, err = a.UpdateProfileBinding("Work", 4, "", true)
			if err != nil || layout["0x04"].ActionCode != "07-00-04" {
				t.Fatalf("restore=%v, %v", layout, err)
			}
			list, err = a.ListProfiles()
			if err != nil || len(list[0].Keys) != 0 {
				t.Fatalf("override not removed: %v, %v", list, err)
			}
			if !active && len(f.Writes) != 2 {
				t.Fatal("inactive edits wrote bindings")
			}
		})
	}
}

func TestActiveEditFailurePreservesLocalSave(t *testing.T) {
	a, f := editFixture(t)
	if err := a.CreateProfile("Work"); err != nil {
		t.Fatal(err)
	}
	f.Exchanges = []keyboard.Exchange{titleExchange("Work"),
		{Write: editPacket(0x52, 0xfa, 3, 0x0c, 0, 0xaa, 9, 0x71, 4, 7, 0, 5), Response: []byte{0x54, 0xe4, 7}}}
	layout, err := a.UpdateProfileBinding("Work", 4, "07-00-05", false)
	if layout != nil || !errors.Is(err, client.ErrMalformedResponse) || !strings.Contains(err.Error(), "saved locally") {
		t.Fatalf("error=%v", err)
	}
	list, err := a.ListProfiles()
	if err != nil || list[0].Keys[0].Action != "07-00-05" {
		t.Fatalf("local save lost: %v, %v", list, err)
	}
}

func TestModelIsolationAndLegacyImport(t *testing.T) {
	a := profileApp(t)
	if err := a.profileStore.Create(profiles.Profile{Title: "Legacy"}); err != nil {
		t.Fatal(err)
	}
	if err := a.CreateProfile("Work"); err != nil {
		t.Fatal(err)
	}
	a.selectedModel = uint16(client.ProductType108Keyboard)
	list, err := a.ListProfiles()
	if err != nil || len(list) != 0 {
		t.Fatalf("models mixed: %v, %v", list, err)
	}
	if err := a.ImportSharedProfiles(); err != nil {
		t.Fatal(err)
	}
	list, err = a.ListProfiles()
	if err != nil || len(list) != 1 || list[0].Title != "Legacy" {
		t.Fatalf("import=%v, %v", list, err)
	}
	a.selectedModel = uint16(client.ProductTypeKeyboard)
	list, err = a.ListProfiles()
	if err != nil || len(list) != 1 || list[0].Title != "Work" {
		t.Fatalf("original model=%v, %v", list, err)
	}
	if _, err := a.profileStore.Load("Legacy"); err != nil {
		t.Fatal("legacy original lost")
	}
}

func TestActiveRename(t *testing.T) {
	a, f := editFixture(t)
	if err := a.CreateProfile("Work"); err != nil {
		t.Fatal(err)
	}
	f.Exchanges = []keyboard.Exchange{titleExchange("Work"),
		{Write: editPacket(0x52, 0x70, 6, 0, 0, 'N', 0, 'e', 0, 'w'), Response: []byte{0x54, 0xe4, 8}},
		titleExchange("New")}
	if err := a.RenameProfile("Work", "New"); err != nil {
		t.Fatal(err)
	}
	list, err := a.ListProfiles()
	if err != nil || list[0].Title != "New" {
		t.Fatalf("renamed=%v, %v", list, err)
	}
}

func TestActiveRefreshReadsKeyboard(t *testing.T) {
	a, f := editFixture(t)
	if err := a.CreateProfile("Work"); err != nil {
		t.Fatal(err)
	}
	f.Exchanges = []keyboard.Exchange{titleExchange("Work"),
		{Write: editPacket(0x52, 0x81), Response: editPacket(0x54, 0x81, 4, 7)},
		{Write: editPacket(0x52, 0x83, 4), Response: editPacket(0x54, 0x83, 4, 7, 0, 5)}}
	layout, err := a.GetProfileLayout("Work")
	if err != nil || layout["0x04"].ActionCode != "07-00-05" {
		t.Fatalf("refresh=%v, %v", layout, err)
	}
	f.Exchanges = append(f.Exchanges, titleExchange("Other"))
	layout, err = a.GetProfileLayout("Work")
	if err != nil || layout["0x04"].ActionCode != "07-00-04" {
		t.Fatalf("inactive=%v, %v", layout, err)
	}
}
