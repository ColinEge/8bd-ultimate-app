package service

import (
	"context"
	"fmt"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/profiles"
)

func ProfileModel(p profiles.Profile) Profile {
	result := Profile{Title: p.Title, Keys: make([]ProfileBinding, 0, len(p.Keys))}
	for _, key := range p.Keys {
		result.Keys = append(result.Keys, ProfileBinding{Key: uint8(key.Key), Action: key.Action.String()})
	}
	return result
}

func ParseProfile(p Profile) (profiles.Profile, error) {
	result := profiles.Profile{Title: p.Title, Keys: make([]client.BindingChange, 0, len(p.Keys))}
	for _, key := range p.Keys {
		action, err := client.ParseAction(key.Action)
		if err != nil {
			return profiles.Profile{}, err
		}
		result.Keys = append(result.Keys, client.BindingChange{Key: client.KeyID(key.Key), Action: action})
	}
	return result, profiles.Validate(result)
}

func (s *EightBitDoService) CaptureProfile(ctx context.Context, id client.DeviceID, title string) (profiles.Profile, error) {
	if err := client.ValidateProfileName(title); err != nil {
		return profiles.Profile{}, err
	}
	var result profiles.Profile
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		bindings, err := k.Bindings(ctx)
		if err != nil {
			return err
		}
		result = profiles.Profile{Title: title, Keys: make([]client.BindingChange, 0, len(bindings))}
		for _, binding := range bindings {
			if binding.State != client.BindingDefault {
				result.Keys = append(result.Keys, client.BindingChange{Key: binding.Key, Action: binding.Action})
			}
		}
		return profiles.Validate(result)
	})
	return result, err
}

func (s *EightBitDoService) ApplyProfile(ctx context.Context, id client.DeviceID, profile profiles.Profile) error {
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		return profiles.Apply(ctx, k, profile)
	})
	if err != nil {
		return fmt.Errorf("profile apply failed; refresh the keyboard before retrying: %w", err)
	}
	return nil
}

func (s *EightBitDoService) GetProfileName(ctx context.Context, id client.DeviceID) (string, error) {
	var title string
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		var err error
		title, err = k.GetProfile(ctx)
		return err
	})
	return title, err
}

func (s *EightBitDoService) SetProfileName(ctx context.Context, id client.DeviceID, title string) error {
	return s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		if err := k.SetProfile(ctx, title); err != nil {
			return fmt.Errorf("profile title write failed; refresh before retrying: %w", err)
		}
		got, err := k.GetProfile(ctx)
		if err != nil {
			return fmt.Errorf("title acknowledged but read-back failed: %w", err)
		}
		if got != title {
			return fmt.Errorf("profile title read-back %q does not match %q", got, title)
		}
		return nil
	})
}

func (s *EightBitDoService) ProfileLayout(ctx context.Context, id client.DeviceID, profile profiles.Profile) (map[string]Key, error) {
	keys := make(map[string]Key)
	err := s.client.WithKeyboard(ctx, id, func(k client.Keyboard) error {
		for _, entry := range k.Keys() {
			key := displayKey(k, client.Binding{Key: entry.ID, State: client.BindingDefault, Effective: entry.Default})
			keys[key.ID] = key
		}
		for _, change := range profile.Keys {
			state := client.BindingAction
			if change.Action == client.ActionDisabled() {
				state = client.BindingDisabled
			}
			key := displayKey(k, client.Binding{Key: change.Key, State: state, Effective: change.Action})
			keys[key.ID] = key
		}
		return nil
	})
	return keys, err
}
