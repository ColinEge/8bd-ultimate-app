package app

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/ColinEge/8bd-ultimate-api/client"
	"github.com/ColinEge/8bd-ultimate-api/profiles"
	"github.com/ColinEge/8bd-ultimate-appgithub.com/ColinEge/8bitdo-ultimate-app/internal/service"
)

type App struct {
	mu             sync.Mutex
	ctx            context.Context
	log            *slog.Logger
	service        *service.EightBitDoService
	selectedDevice client.DeviceID
	selectedModel  uint16
	profileStore   *profiles.Store
	profileConfig  string
}

func New(log *slog.Logger) *App {
	return &App{
		ctx: context.Background(), log: log,
		service: service.NewService(client.New(client.WithTimeout(2 * time.Second))),
	}
}

func (a *App) Startup(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.ctx = ctx
}

func (a *App) Shutdown(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.service.Close(); err != nil {
		a.log.Error("Failed to close keyboard client", "error", err)
	}
}

func (a *App) ScanDevices() ([]service.Device, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.service.ScanDevices(a.ctx)
}

func (a *App) SelectDevice(id string) (*service.Keyboard, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	kbd, err := a.service.GetDetails(a.ctx, client.DeviceID(id))
	if err != nil {
		return nil, err
	}
	a.selectedDevice = client.DeviceID(id)
	a.selectedModel = kbd.Model
	return kbd, nil
}

func (a *App) requireSelection() error {
	if a.selectedDevice == "" {
		return errors.New("no device is selected, please select a device first")
	}
	return nil
}

func (a *App) GetSelectedDevice() (*service.Keyboard, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.GetDetails(a.ctx, a.selectedDevice)
}

func (a *App) GetActions() ([]service.ActionCategory, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.GetActions(a.ctx, a.selectedDevice)
}

func (a *App) GetBinding(key uint8) (*service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.GetBinding(a.ctx, a.selectedDevice, key)
}

func (a *App) BindAction(key uint8, action string) (*service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.BindAction(a.ctx, a.selectedDevice, key, action)
}

func (a *App) DisableKey(key uint8) (*service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.DisableKey(a.ctx, a.selectedDevice, key)
}

func (a *App) RestoreKey(key uint8) (*service.Key, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if err := a.requireSelection(); err != nil {
		return nil, err
	}
	return a.service.RestoreKey(a.ctx, a.selectedDevice, key)
}
