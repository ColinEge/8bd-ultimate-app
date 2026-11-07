package app

import (
	"context"
	"errors"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/profiles"
	"github.com/ColinEge/8bd-ultimate-api/test/fixtures/keyboard"
	"github.com/ColinEge/8bd-ultimate-appgithub.com/ColinEge/8bitdo-ultimate-app/internal/service"
)

type localKeyboard struct{ client.Keyboard }

func (localKeyboard) GetProfile(context.Context) (string, error) { return "", nil }

type localClient struct{ client.Client }

func (localClient) WithKeyboard(ctx context.Context, id client.DeviceID, fn func(client.Keyboard) error) error {
	return fn(localKeyboard{})
}

func profileApp(t *testing.T) *App {
	t.Helper()
	store, err := profiles.NewStore(profiles.WithDirectory(t.TempDir()))
	if err != nil {
		t.Fatal(err)
	}
	return &App{ctx: context.Background(), log: slog.Default(), profileStore: store,
		selectedDevice: "001-003", selectedModel: uint16(client.ProductTypeKeyboard), service: service.NewService(localClient{}),
		profileConfig: filepath.Join(t.TempDir(), "profiles.json")}
}

func TestLocalProfiles(t *testing.T) {
	a := profileApp(t)
	if err := a.CreateProfile("Gaming"); err != nil {
		t.Fatal(err)
	}
	if err := a.CreateProfile("Gaming"); !errors.Is(err, fs.ErrExist) {
		t.Fatalf("duplicate: %v", err)
	}
	p := service.Profile{Title: "Gaming", Keys: []service.ProfileBinding{{Key: 4, Action: "07-00-05"}}}
	if err := a.SaveProfile(p); err != nil {
		t.Fatal(err)
	}
	bad := service.Profile{Title: "Gaming", Keys: []service.ProfileBinding{{Key: 4, Action: "broken"}}}
	if err := a.SaveProfile(bad); !errors.Is(err, client.ErrInvalidAction) {
		t.Fatalf("invalid action: %v", err)
	}
	bad.Keys[0] = service.ProfileBinding{Key: 0, Action: "07-00-05"}
	if err := a.SaveProfile(bad); !errors.Is(err, client.ErrInvalidKey) {
		t.Fatalf("invalid key: %v", err)
	}
	if err := a.RenameProfile("Gaming", "Work"); err != nil {
		t.Fatal(err)
	}
	list, err := a.ListProfiles()
	if err != nil || len(list) != 1 || list[0].Title != "Work" ||
		len(list[0].Keys) != 1 || list[0].Keys[0] != p.Keys[0] {
		t.Fatalf("profiles=%+v, error=%v", list, err)
	}
	a.selectedDevice = ""
	if err := a.ApplyProfile("Work"); err == nil {
		t.Fatal("apply without device succeeded")
	}
	if err := a.CaptureProfile("Snapshot"); err == nil {
		t.Fatal("capture without device succeeded")
	}
	a.selectedDevice = "001-003"
	if err := a.DeleteProfile("Work"); err != nil {
		t.Fatal(err)
	}
	list, err = a.ListProfiles()
	if err != nil || list == nil || len(list) != 0 {
		t.Fatalf("empty list=%v, error=%v", list, err)
	}
}

func TestProfileDirectoryPersistence(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	t.Setenv("XDG_DATA_HOME", t.TempDir())
	a := &App{selectedDevice: "001-003", selectedModel: uint16(client.ProductTypeKeyboard)}
	defaultPath, err := a.GetProfileDirectory()
	if err != nil {
		t.Fatal(err)
	}
	if err := a.CreateProfile("Default"); err != nil {
		t.Fatal(err)
	}
	chosen := t.TempDir()
	path, err := a.setProfileDirectory(chosen)
	if err != nil || path != chosen {
		t.Fatalf("override=%q, error=%v", path, err)
	}
	if err := a.CreateProfile("Custom"); err != nil {
		t.Fatal(err)
	}
	restarted := &App{selectedDevice: "001-003", selectedModel: uint16(client.ProductTypeKeyboard)}
	path, err = restarted.GetProfileDirectory()
	if err != nil || path != chosen {
		t.Fatalf("restarted=%q, error=%v", path, err)
	}
	path, err = restarted.ResetProfileDirectory()
	if err != nil || path != defaultPath {
		t.Fatalf("reset=%q, error=%v", path, err)
	}
	list, err := restarted.ListProfiles()
	if err != nil || len(list) != 1 || list[0].Title != "Default" {
		t.Fatalf("default files lost: %v, %v", list, err)
	}
	path, err = a.GetProfileDirectory()
	if err != nil || path != chosen {
		t.Fatalf("old store changed: %q, %v", path, err)
	}
	blocked := filepath.Join(t.TempDir(), "not-a-directory")
	if err := os.WriteFile(blocked, []byte("file"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := a.setProfileDirectory(blocked); err == nil {
		t.Fatal("invalid override succeeded")
	}
	path, _ = a.GetProfileDirectory()
	if path != chosen {
		t.Fatal("failed override changed active store")
	}
	a.profileConfig = filepath.Join(blocked, "profiles.json")
	if _, err := a.setProfileDirectory(t.TempDir()); err == nil {
		t.Fatal("settings persistence failure ignored")
	}
	path, _ = a.GetProfileDirectory()
	if path != chosen {
		t.Fatal("failed persistence changed active store")
	}
}

func TestProfileKeyboardWorkflow(t *testing.T) {
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
	packet := func(prefix ...byte) []byte {
		data := make([]byte, 33)
		copy(data, prefix)
		return data
	}
	f.Exchanges = []keyboard.Exchange{
		{Write: packet(0x52, 0x81), Response: packet(0x54, 0x81, 4, 7)},
		{Write: packet(0x52, 0x83, 4), Response: packet(0x54, 0x83, 4, 7, 0, 5)},
		{Write: packet(0x52, 0x70, 8, 0, 0, 'W', 0, 'o', 0, 'r', 0, 'k'), Response: []byte{0x54, 0xe4, 8}},
		{Write: packet(0x52, 0xfa, 3, 0x0c, 0, 0xaa, 9, 0x71, 4, 7, 0, 5), Response: []byte{0x54, 0xe4, 8}},
		{Write: packet(0x52, 0x80), Response: packet(0x54, 0x80, 8, 0, 0, 'W', 0, 'o', 0, 'r', 0, 'k')},
	}
	if err := a.CaptureProfile("Work"); err != nil {
		t.Fatal(err)
	}
	list, err := a.ListProfiles()
	if err != nil || len(list) != 1 || len(list[0].Keys) != 1 || list[0].Keys[0].Action != "07-00-05" {
		t.Fatalf("capture=%v, error=%v", list, err)
	}
	if len(f.Writes) != 2 {
		t.Fatal("capture wrote to keyboard")
	}
	if err := a.ApplyProfile("Work"); err != nil {
		t.Fatal(err)
	}
	title, err := a.GetProfileName()
	if err != nil || title != "Work" {
		t.Fatalf("device title=%q, error=%v", title, err)
	}
	f.Exchanges = append(f.Exchanges, keyboard.Exchange{
		Write: packet(0x52, 0x70, 8, 0, 0, 'W', 0, 'o', 0, 'r', 0, 'k'), Response: []byte{0x54, 0xe4, 7},
	})
	if err := a.ApplyProfile("Work"); !errors.Is(err, client.ErrMalformedResponse) {
		t.Fatalf("unknown status treated as success: %v", err)
	}
	list, err = a.ListProfiles()
	if err != nil || list[0].Keys[0].Action != "07-00-05" {
		t.Fatal("apply changed local profile")
	}
}
