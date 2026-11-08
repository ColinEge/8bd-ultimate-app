package service

import (
	"context"
	"fmt"

	"github.com/ColinEge/8bd-ultimate-api/client"
)

type EightBitDoService struct {
	client client.Client
}

func NewService(c client.Client) *EightBitDoService {
	return &EightBitDoService{client: c}
}

func (s *EightBitDoService) Close() error {
	return s.client.Close()
}

func (s *EightBitDoService) ScanDevices(ctx context.Context) ([]Device, error) {
	devs, err := s.client.Scan(ctx)
	if err != nil {
		return nil, err
	}
	devices := make([]Device, len(devs))
	for i, dev := range devs {
		devices[i] = Device{DeviceID: string(dev.ID), ProductName: dev.Product.String()}
	}
	return devices, nil
}

func (s *EightBitDoService) GetDetails(ctx context.Context, id client.DeviceID) (*Keyboard, error) {
	var result *Keyboard
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		info, err := k.Info(ctx)
		if err != nil {
			return err
		}
		bindings, err := k.Bindings(ctx)
		if err != nil {
			return err
		}
		keys := make(map[string]Key)
		for _, entry := range k.Keys() {
			key := displayKey(k, client.Binding{
				Key: entry.ID, State: client.BindingDefault, Effective: entry.Default,
			})
			keys[key.ID] = key
		}
		for _, binding := range bindings {
			key := displayKey(k, binding)
			keys[key.ID] = key
		}
		result = &Keyboard{
			Model:    uint16(info.Device.Product),
			DeviceID: string(info.Device.ID), ProductName: info.Device.Product.String(),
			HasKeypad: info.Device.Product == client.ProductType108Keyboard, Keys: keys,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *EightBitDoService) GetActions(ctx context.Context, id client.DeviceID) ([]ActionCategory, error) {
	categories := make([]ActionCategory, 0)
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		index := make(map[string]int)
		for _, entry := range k.Actions() {
			name := string(entry.Category)
			i, ok := index[name]
			if !ok {
				i = len(categories)
				index[name] = i
				categories = append(categories, ActionCategory{Name: name, Actions: make([]Action, 0)})
			}
			categories[i].Actions = append(categories[i].Actions, Action{
				ActionCode: entry.Action.String(), ActionLabel: entry.Label,
			})
		}
		return nil
	})
	return categories, err
}

func displayKey(k client.Keyboard, binding client.Binding) Key {
	key := Key{
		ID: binding.Key.String(), PhysicalKey: uint8(binding.Key),
		KeyLabel: binding.Key.String(), BindingState: string(binding.State),
		ActionCode: binding.Effective.String(), ActionLabel: binding.Effective.String(),
	}
	for _, entry := range k.Keys() {
		if entry.ID == binding.Key {
			key.KeyLabel = entry.Label
			break
		}
	}
	for _, entry := range k.Actions() {
		if entry.Action == binding.Effective {
			key.ActionLabel = entry.Label
			break
		}
	}
	return key
}

func (s *EightBitDoService) GetBinding(ctx context.Context, id client.DeviceID, key uint8) (*Key, error) {
	return s.withBinding(ctx, id, key, nil)
}

func (s *EightBitDoService) BindAction(ctx context.Context, id client.DeviceID, key uint8, text string) (*Key, error) {
	action, err := client.ParseAction(text)
	if err != nil {
		return nil, err
	}
	return s.withBinding(ctx, id, key, func(k client.Keyboard) error {
		return k.Bind(ctx, client.KeyID(key), action)
	})
}

func (s *EightBitDoService) DisableKey(ctx context.Context, id client.DeviceID, key uint8) (*Key, error) {
	return s.withBinding(ctx, id, key, func(k client.Keyboard) error {
		return k.Disable(ctx, client.KeyID(key))
	})
}

func (s *EightBitDoService) ResetToDefault(ctx context.Context, id client.DeviceID, key uint8) (*Key, error) {
	return s.withBinding(ctx, id, key, func(k client.Keyboard) error {
		return k.ResetToDefault(ctx, client.KeyID(key))
	})
}

func (s *EightBitDoService) withBinding(ctx context.Context, id client.DeviceID, key uint8, write func(client.Keyboard) error) (*Key, error) {
	var result *Key
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		if write != nil {
			if err := write(k); err != nil {
				return fmt.Errorf("key %s write failed; read its binding before retrying: %w", client.KeyID(key), err)
			}
		}
		binding, err := k.Binding(ctx, client.KeyID(key))
		if err != nil {
			if write != nil {
				return fmt.Errorf("write acknowledged but read-back failed; read the binding before retrying: %w", err)
			}
			return err
		}
		key := displayKey(k, binding)
		result = &key
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
