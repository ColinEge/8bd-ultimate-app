package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/test/fixtures/keyboard"
)

func fixtureService(t *testing.T) (*EightBitDoService, *keyboard.KeyboardFixture) {
	t.Helper()
	f := &keyboard.KeyboardFixture{}
	if err := f.Reset(); err != nil {
		t.Fatal(err)
	}
	s := NewService(f.Keyboard)
	t.Cleanup(func() {
		if err := f.AssertComplete(); err != nil {
			t.Error(err)
		}
		if err := s.Close(); err != nil {
			t.Error(err)
		}
	})
	return s, f
}

func packet(prefix ...byte) []byte {
	data := make([]byte, 33)
	copy(data, prefix)
	return data
}

func bindPacket(key byte, action client.Action) []byte {
	data := packet(0x52, 0xfa, 0x03, 0x0c, 0x00, 0xaa, 0x09, 0x71, key)
	copy(data[9:], action[:])
	return data
}

func readExchange(key byte, action client.Action) keyboard.Exchange {
	reply := packet(0x54, 0x83, key)
	copy(reply[3:], action[:])
	return keyboard.Exchange{Write: packet(0x52, 0x83, key), Response: reply}
}

func TestDiscoveryAndCatalog(t *testing.T) {
	s, f := fixtureService(t)
	ctx := context.Background()
	devices, err := s.ScanDevices(ctx)
	if err != nil || len(devices) != 1 || devices[0].DeviceID != string(f.Device.ID) {
		t.Fatalf("devices=%v, error=%v", devices, err)
	}
	actions, err := s.GetActions(ctx, f.Device.ID)
	if err != nil || len(actions) == 0 {
		t.Fatalf("actions=%v, error=%v", actions, err)
	}
	found := false
	for _, category := range actions {
		if category.Name == "" || len(category.Actions) == 0 {
			t.Fatalf("empty category: %+v", category)
		}
		for _, action := range category.Actions {
			if action.ActionCode == client.ActionB().String() {
				found = true
			}
		}
	}
	if !found {
		t.Fatal("missing B action")
	}
	f.Connected = false
	devices, err = s.ScanDevices(ctx)
	if err != nil || devices == nil || len(devices) != 0 {
		t.Fatalf("empty scan=%v, error=%v", devices, err)
	}
	if _, err := s.GetDetails(ctx, f.Device.ID); !errors.Is(err, client.ErrDeviceNotFound) {
		t.Fatalf("disconnected selection: %v", err)
	}
}

func TestLayout(t *testing.T) {
	s, f := fixtureService(t)
	f.Exchanges = []keyboard.Exchange{
		{Write: packet(0x52, 0x81), Response: packet(0x54, 0x81, 4, 7)},
		readExchange(4, client.ActionB()),
	}
	kbd, err := s.GetDetails(context.Background(), f.Device.ID)
	if err != nil {
		t.Fatal(err)
	}
	if kbd.DeviceID != string(f.Device.ID) || kbd.HasKeypad {
		t.Fatalf("metadata: %+v", kbd)
	}
	a := kbd.Keys["0x04"]
	if a.PhysicalKey != 4 || a.KeyLabel != "A" || a.ActionCode != "07-00-05" || a.BindingState != "action" {
		t.Fatalf("A remapped to B: %+v", a)
	}
	for id, action := range map[string]string{"0x64": "07-e0-00", "0x65": "07-e1-00", "0x28": "07-00-28"} {
		key, ok := kbd.Keys[id]
		if !ok || key.ActionCode != action || key.BindingState != "default" {
			t.Fatalf("default %s: %+v", id, key)
		}
	}
	if _, ok := kbd.Keys["0xe0"]; ok {
		t.Fatal("HID modifier usage exposed as physical ID")
	}
	if f.SessionCloses != 1 {
		t.Fatalf("connections closed=%d", f.SessionCloses)
	}
	data, err := json.Marshal(a)
	if err != nil || !strings.Contains(string(data), `"physicalKey":4`) ||
		!strings.Contains(string(data), `"actionCode":"07-00-05"`) {
		t.Fatalf("frontend DTO=%s, error=%v", data, err)
	}
}

func TestSingleKeyWorkflow(t *testing.T) {
	for _, tt := range []struct {
		name      string
		key       byte
		action    client.Action
		readBack  client.Action
		state     string
		effective string
	}{
		{"bind", 4, client.ActionB(), client.ActionB(), "action", "07-00-05"},
		{"disable", 4, client.ActionDisabled(), client.ActionDisabled(), "disabled", "07-00-00"},
		{"restore Ctrl", 0x64, client.ActionLeftCtrl(), client.Action{}, "default", "07-e0-00"},
		{"restore Shift explicit default", 0x65, client.ActionLeftShift(), client.ActionLeftShift(), "action", "07-e1-00"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, f := fixtureService(t)
			f.Exchanges = []keyboard.Exchange{
				{Write: bindPacket(tt.key, tt.action), Response: []byte{0x54, 0xe4, 8}},
				readExchange(tt.key, tt.readBack),
			}
			ctx := context.Background()
			var key *Key
			var err error
			switch tt.name {
			case "bind":
				key, err = s.BindAction(ctx, f.Device.ID, tt.key, tt.action.String())
			case "disable":
				key, err = s.DisableKey(ctx, f.Device.ID, tt.key)
			default:
				key, err = s.RestoreKey(ctx, f.Device.ID, tt.key)
			}
			if err != nil {
				t.Fatal(err)
			}
			if key.PhysicalKey != tt.key || key.BindingState != tt.state || key.ActionCode != tt.effective {
				t.Fatalf("read-back: %+v", key)
			}
			if f.SessionCloses != 1 {
				t.Fatalf("connections closed=%d", f.SessionCloses)
			}
		})
	}
}

func TestBindingErrors(t *testing.T) {
	readErr := errors.New("USB read failed")
	for _, tt := range []struct {
		name      string
		action    string
		exchanges []keyboard.Exchange
		cause     error
	}{
		{"parse", "not-an-action", nil, client.ErrInvalidAction},
		{"unsupported", "07-00-ff", nil, client.ErrUnsupported},
		{"device rejected", "07-00-05", []keyboard.Exchange{
			{Write: bindPacket(4, client.ActionB()), Response: []byte{0x54, 0xe4, 9}},
		}, client.ErrDeviceRejected},
		{"unknown status", "07-00-05", []keyboard.Exchange{
			{Write: bindPacket(4, client.ActionB()), Response: []byte{0x54, 0xe4, 7}},
		}, client.ErrMalformedResponse},
		{"read-back failed", "07-00-05", []keyboard.Exchange{
			{Write: bindPacket(4, client.ActionB()), Response: []byte{0x54, 0xe4, 8}},
			{Write: packet(0x52, 0x83, 4), ReadErr: readErr},
		}, readErr},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s, f := fixtureService(t)
			f.Exchanges = tt.exchanges
			key, err := s.BindAction(context.Background(), f.Device.ID, 4, tt.action)
			if key != nil || !errors.Is(err, tt.cause) {
				t.Fatalf("key=%v, error=%v, want=%v", key, err, tt.cause)
			}
			if tt.name == "read-back failed" && !strings.Contains(err.Error(), "write acknowledged") {
				t.Fatalf("missing uncertain outcome guidance: %v", err)
			}
		})
	}
}

func TestReadBindingAndCleanupErrors(t *testing.T) {
	s, f := fixtureService(t)
	ctx := context.Background()
	if _, err := s.GetBinding(ctx, f.Device.ID, 0); !errors.Is(err, client.ErrInvalidKey) {
		t.Fatalf("invalid key: %v", err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := s.ScanDevices(cancelled); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation: %v", err)
	}
	f.Exchanges = []keyboard.Exchange{readExchange(4, client.ActionB())}
	closeErr := errors.New("connection close failed")
	f.SessionCloseErr = closeErr
	key, err := s.GetBinding(ctx, f.Device.ID, 4)
	if key != nil || !errors.Is(err, closeErr) {
		t.Fatalf("cleanup: key=%v, error=%v", key, err)
	}
	f.SessionCloseErr = nil
}
