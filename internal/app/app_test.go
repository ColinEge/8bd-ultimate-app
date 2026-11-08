package app

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/test/fixtures/keyboard"
	"github.com/ColinEge/8bd-ultimate-app/internal/service"
)

func TestAppSelectionAndRebinding(t *testing.T) {
	f := &keyboard.KeyboardFixture{}
	if err := f.Reset(); err != nil {
		t.Fatal(err)
	}
	a := &App{ctx: context.Background(), log: slog.Default(), service: service.NewService(f.Keyboard)}
	t.Cleanup(func() {
		if err := f.AssertComplete(); err != nil {
			t.Error(err)
		}
		a.Shutdown(context.Background())
		if f.ProviderCloses != 1 {
			t.Errorf("provider closes=%d, want 1", f.ProviderCloses)
		}
	})
	for name, operation := range map[string]func() error{
		"selected device":  func() error { _, err := a.GetSelectedDevice(); return err },
		"actions":          func() error { _, err := a.GetActions(); return err },
		"read":             func() error { _, err := a.GetBinding(4); return err },
		"bind":             func() error { _, err := a.BindAction(4, "07-00-05"); return err },
		"disable":          func() error { _, err := a.DisableKey(4); return err },
		"reset-to-default": func() error { _, err := a.ResetToDefault(4); return err },
	} {
		if err := operation(); err == nil {
			t.Fatalf("%s succeeded without selection", name)
		}
	}
	devices, err := a.ScanDevices()
	if err != nil || len(devices) != 1 || a.selectedDevice != "" {
		t.Fatalf("scan must not select: devices=%v, error=%v", devices, err)
	}
	packet := func(prefix ...byte) []byte {
		data := make([]byte, 33)
		copy(data, prefix)
		return data
	}
	f.Exchanges = []keyboard.Exchange{
		{Write: packet(0x52, 0x81), Response: packet(0x54, 0x81)},
		{Write: packet(0x52, 0xfa, 3, 0x0c, 0, 0xaa, 9, 0x71, 4, 7, 0, 5), Response: []byte{0x54, 0xe4, 8}},
		{Write: packet(0x52, 0x83, 4), Response: packet(0x54, 0x83, 4, 7, 0, 5)},
	}
	kbd, err := a.SelectDevice(devices[0].DeviceID)
	if err != nil || kbd.Keys["0x04"].BindingState != "default" {
		t.Fatalf("selection: keyboard=%v, error=%v", kbd, err)
	}
	key, err := a.BindAction(4, client.ActionB().String())
	if err != nil || key.PhysicalKey != 4 || key.ActionCode != "07-00-05" || key.KeyLabel != "A" {
		t.Fatalf("physical A, action B: key=%v, error=%v", key, err)
	}
	if _, err := a.SelectDevice("missing"); !errors.Is(err, client.ErrDeviceNotFound) {
		t.Fatalf("failed selection: %v", err)
	}
	if a.selectedDevice != f.Device.ID {
		t.Fatal("failed selection replaced the selected device")
	}
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	a.Startup(cancelled)
	if _, err := a.ScanDevices(); !errors.Is(err, context.Canceled) {
		t.Fatalf("startup context not propagated: %v", err)
	}
}
