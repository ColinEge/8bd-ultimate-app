package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/profiles"
	"github.com/ColinEge/8bd-ultimate-app/internal/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type profileSettings struct {
	Directory string `json:"directory"`
}

func (a *App) store() (*profiles.Store, error) {
	if a.profileStore != nil {
		return a.profileStore, nil
	}
	base, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	a.profileConfig = filepath.Join(base, "8bd-kbd-app", "profiles.json")
	data, err := os.ReadFile(a.profileConfig)
	var options []profiles.StoreOption
	if err == nil {
		var settings profileSettings
		if err := json.Unmarshal(data, &settings); err != nil {
			return nil, fmt.Errorf("read profile settings: %w", err)
		}
		if settings.Directory != "" {
			options = append(options, profiles.WithDirectory(settings.Directory))
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	store, err := profiles.NewStore(options...)
	if err != nil {
		return nil, err
	}
	a.profileStore = store
	return store, nil
}

func (a *App) GetProfileDirectory() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.store()
	if err != nil {
		return "", err
	}
	return store.Directory(), nil
}

func (a *App) setProfileDirectory(directory string) (string, error) {
	var options []profiles.StoreOption
	if directory != "" {
		options = append(options, profiles.WithDirectory(directory))
	}
	store, err := profiles.NewStore(options...)
	if err != nil {
		return "", err
	}
	if a.profileConfig == "" {
		base, err := os.UserConfigDir()
		if err != nil {
			return "", err
		}
		a.profileConfig = filepath.Join(base, "8bd-kbd-app", "profiles.json")
	}
	if directory != "" {
		directory = store.Directory()
	}
	if err := os.MkdirAll(filepath.Dir(a.profileConfig), 0700); err != nil {
		return "", err
	}
	data, err := json.Marshal(profileSettings{Directory: directory})
	if err != nil {
		return "", err
	}
	file, err := os.CreateTemp(filepath.Dir(a.profileConfig), ".profiles-")
	if err != nil {
		return "", err
	}
	name := file.Name()
	_, writeErr := file.Write(data)
	if writeErr == nil {
		writeErr = file.Sync()
	}
	err = errors.Join(writeErr, file.Close())
	if err == nil {
		err = os.Rename(name, a.profileConfig)
	}
	if err != nil {
		return "", errors.Join(err, os.Remove(name))
	}
	a.profileStore = store
	return store.Directory(), nil
}

func (a *App) ChooseProfileDirectory() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.store()
	if err != nil {
		return "", err
	}
	directory, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Choose profile save directory", DefaultDirectory: store.Directory(),
	})
	if err != nil {
		return "", err
	}
	if directory == "" {
		return store.Directory(), nil
	}
	return a.setProfileDirectory(directory)
}

func (a *App) ResetProfileDirectory() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.setProfileDirectory("")
}

func (a *App) ListProfiles() ([]service.Profile, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return nil, err
	}
	list, err := store.List()
	if err != nil {
		return nil, err
	}
	result := make([]service.Profile, 0, len(list))
	for _, profile := range list {
		result = append(result, service.ProfileModel(profile))
	}
	return result, nil
}

func (a *App) CreateProfile(title string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return err
	}
	return store.Create(profiles.Profile{Title: title})
}

func (a *App) SaveProfile(profile service.Profile) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	model, err := service.ParseProfile(profile)
	if err != nil {
		return err
	}
	store, err := a.modelStore()
	if err != nil {
		return err
	}
	active, err := a.service.GetProfileName(a.ctx, a.selectedDevice)
	if err != nil {
		return err
	}
	if active == model.Title {
		return fmt.Errorf("use the key editor to edit the active profile")
	}
	return store.Save(model)
}

func (a *App) RenameProfile(title, newTitle string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return err
	}
	active, err := a.service.GetProfileName(a.ctx, a.selectedDevice)
	if err != nil {
		return err
	}
	if err := store.Rename(title, newTitle); err != nil {
		return err
	}
	if active == title {
		if err := a.service.SetProfileName(a.ctx, a.selectedDevice, newTitle); err != nil {
			return fmt.Errorf("local profile renamed to %q, but device rename failed: %w", newTitle, err)
		}
	}
	return nil
}

func (a *App) DeleteProfile(title string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return err
	}
	active, err := a.service.GetProfileName(a.ctx, a.selectedDevice)
	if err != nil {
		return err
	}
	if err := store.Delete(title); err != nil {
		return err
	}
	if active == title {
		if err := a.service.DeleteProfile(a.ctx, a.selectedDevice); err != nil {
			return fmt.Errorf("local profile deleted, but device profile deletion failed: %w", err)
		}
	}
	return nil
}

// Capture creates a new local profile without changing the device title.
func (a *App) CaptureProfile(title string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return err
	}
	store, err := a.modelStore()
	if err != nil {
		return err
	}
	profile, err := a.service.CaptureProfile(a.ctx, a.selectedDevice, title)
	if err != nil {
		return err
	}
	return store.Create(profile)
}

func (a *App) ApplyProfile(title string) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return err
	}
	store, err := a.modelStore()
	if err != nil {
		return err
	}

	profile, err := store.Load(title)
	if err != nil {
		return err
	}
	return a.service.ApplyProfile(a.ctx, a.selectedDevice, profile)
}

func (a *App) modelStore() (*profiles.Store, error) {
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	if a.selectedModel != uint16(client.ProductTypeKeyboard) && a.selectedModel != uint16(client.ProductType108Keyboard) {
		return nil, fmt.Errorf("unknown selected keyboard model")
	}
	root, err := a.store()
	if err != nil {
		return nil, err
	}
	return profiles.NewStore(profiles.WithDirectory(filepath.Join(root.Directory(), fmt.Sprintf("model-%04x", a.selectedModel))))
}

// ImportSharedProfiles copies legacy shared profiles without moving or deleting them.
func (a *App) ImportSharedProfiles() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	target, err := a.modelStore()
	if err != nil {
		return err
	}
	list, err := a.profileStore.List()
	if err != nil {
		return err
	}
	for _, p := range list {
		if err := target.Create(p); err != nil {
			return fmt.Errorf("import stopped at %q (earlier profiles may be copied): %w", p.Title, err)
		}
	}
	return nil
}

func (a *App) GetProfileLayout(title string) (map[string]service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return nil, err
	}
	p, err := store.Load(title)
	if err != nil {
		return nil, err
	}
	active, err := a.service.GetProfileName(a.ctx, a.selectedDevice)
	if err != nil {
		return nil, err
	}
	if active == title {
		kbd, err := a.service.GetDetails(a.ctx, a.selectedDevice)
		if err != nil {
			return nil, err
		}
		return kbd.Keys, nil
	}
	return a.service.ProfileLayout(a.ctx, a.selectedDevice, p)
}

// Each edit persists locally first. Active edits also write only the changed key.
func (a *App) UpdateProfileBinding(title string, key uint8, action string, restore bool) (map[string]service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	store, err := a.modelStore()
	if err != nil {
		return nil, err
	}
	p, err := store.Load(title)
	if err != nil {
		return nil, err
	}
	if err := client.ValidateBindings([]client.BindingChange{{Key: client.KeyID(key), Action: client.ActionA()}}); err != nil {
		return nil, err
	}
	var parsed client.Action
	if !restore {
		parsed, err = client.ParseAction(action)
		if err != nil {
			return nil, err
		}
		if err := client.ValidateBindings([]client.BindingChange{{Key: client.KeyID(key), Action: parsed}}); err != nil {
			return nil, err
		}
	}
	active, err := a.service.GetProfileName(a.ctx, a.selectedDevice)
	if err != nil {
		return nil, err
	}
	changes := make([]client.BindingChange, 0, len(p.Keys)+1)
	for _, change := range p.Keys {
		if change.Key != client.KeyID(key) {
			changes = append(changes, change)
		}
	}
	if !restore {
		changes = append(changes, client.BindingChange{Key: client.KeyID(key), Action: parsed})
	}
	p.Keys = changes
	if err := store.Save(p); err != nil {
		return nil, err
	}
	if active == title {
		var readBack *service.Key
		if restore {
			readBack, err = a.service.RestoreKey(a.ctx, a.selectedDevice, key)
		} else {
			readBack, err = a.service.BindAction(a.ctx, a.selectedDevice, key, action)
		}
		if err != nil {
			return nil, fmt.Errorf("profile saved locally but active-key update failed; refresh before retrying: %w", err)
		}
		layout, err := a.service.ProfileLayout(a.ctx, a.selectedDevice, p)
		if err != nil {
			return nil, fmt.Errorf("profile saved and key written, but layout refresh failed: %w", err)
		}
		expected := parsed.String()
		if restore {
			expected = layout[readBack.ID].ActionCode
		}
		if readBack.ActionCode != expected {
			return nil, fmt.Errorf("profile saved locally but key read-back %s does not match %s; refresh before retrying", readBack.ActionCode, expected)
		}
		layout[readBack.ID] = *readBack
		return layout, nil
	}
	return a.service.ProfileLayout(a.ctx, a.selectedDevice, p)
}
func (a *App) GetProfileName() (string, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return "", err
	}
	return a.service.GetProfileName(a.ctx, a.selectedDevice)
}
