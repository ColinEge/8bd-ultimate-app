# 8BitDo Keyboard Mapper

A desktop app for managing key bindings and local profiles on supported
8BitDo keyboards.

## Installing

To build from source, install Go (see [go.mod](./go.mod)), Node.js/npm,
libusb development libraries, and the
[Wails platform prerequisites](https://wails.io/docs/gettingstarted/installation).
Then run from this directory:

```sh
go install github.com/wailsapp/wails/v2/cmd/wails@v2.16.0
wails doctor
wails build
```

Wails installs the frontend dependencies and builds the app into `build/bin`.
Launch the generated `8bdmapper` executable or platform app bundle.
On Linux systems using WebKitGTK 4.1 instead of 4.0, use
`wails build -tags webkit2_41`.

### Linux USB access

Install the supplied [udev rules](./50-8bitdo-kdb.rules) to allow USB
configuration access without running the app as root:

```sh
sudo install -m 0644 50-8bitdo-kdb.rules /etc/udev/rules.d/50-8bitdo-kdb.rules
sudo udevadm control --reload-rules
```

Unplug and reconnect the keyboard after installing the rules.

## Usage

1. Connect the keyboard over USB with the wireless selector off.
2. Launch the app and select your keyboard. Use **Rescan** if it is not listed.
3. Choose **New profile** to create an empty profile or copy the keyboard's
   current bindings.
4. Open a profile's three-dot menu and choose **Edit profile**. Select a key,
   then **Rebind**, **Disable**, or **Reset to default**.
5. Choose **Set as active** to apply a saved profile to the keyboard.

Edits save immediately. Inactive profile edits will not modify the keyboard;
active profile edits also update the keyboard and read back the result.
Applying a profile changes only its saved keys, leaving omitted keys unchanged.
Renaming or deleting an active profile also updates the keyboard's profile title.

Back up your configuration before changing it. Writes have no rollback or
automatic retry. If an operation fails, use **Refresh profile** in the editor
or **Refresh** in the profile list before deciding whether to retry.
Super A/B support key binding edits; macros are not currently supported.

Profiles are stored locally and grouped by keyboard model. Under
**Profile storage**, **Choose folder** changes the storage location and
**Use default** resets it. Changing folders does not move or delete profiles.

## Contributing

Keep changes focused and add or update tests for changed behavior.
The app uses the [8BitDo keyboard API](../8bd-ultimate-api); its dependency is
declared in [go.mod](./go.mod).

For live development, run from this directory:

```sh
wails dev
```

Use `wails dev -tags webkit2_41` where WebKitGTK 4.1 is required.
Before submitting changes, format changed Go files with `gofmt` and run:

```sh
go test ./internal/...
cd frontend
npm install
npm run check
npm run build
```

Automated tests use scripted keyboards and do not access USB. Keep physical
keyboard writes out of default tests and obtain explicit approval before
hardware testing. Update this README when installation or usage changes.
