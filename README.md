Trayscale (Fork with Auto Service Management)
==============================================

[![Go Report Card](https://goreportcard.com/badge/deedles.dev/trayscale)](https://goreportcard.com/report/deedles.dev/trayscale)

Trayscale is an unofficial GUI interface for the Tailscale daemon, particularly for use on Linux, that provides many features that the official Linux client does not.

> [!NOTE]
> **Fork Highlights (Auto Service Lifecycle Management):**
> Versi fork ini menambahkan integrasi kontrol `tailscaled.service` (systemd) secara otomatis:
> - **Auto-Start saat ON:** Jika daemon `tailscaled.service` mati, mengaktifkan switch (ON) akan otomatis memanggil `pkexec systemctl start tailscaled.service` (dialog otentikasi Polkit desktop akan muncul) sebelum menghubungkan ke jaringan Tailscale.
> - **Auto-Stop saat OFF:** Saat mematikan switch (OFF), Trayscale memutuskan koneksi (`tailscale down`) dan langsung menghentikan service `tailscaled.service` (`pkexec systemctl stop tailscaled.service`), memastikan tidak ada resource background yang aktif.
> - *Persyaratan:* Memerlukan `systemd` dan `policykit-1` / `polkit` (`pkexec`). Fitur ini ditujukan untuk **Native Build**.

_Disclaimer: This project is in a beta state. There may still be undiscovered bugs or compatibility issues. Use at your own risk._

![screenshot](https://github.com/user-attachments/assets/30f69ec9-630b-4d9a-98a9-cb31712c60e7)

Tailscale Config
----------------

Trayscale interfaces with the Tailscale daemon, `tailscaled`, to perform many of its operations. In order for this to work, the daemon must have been configured with the current user as the "operator". To do this, run `sudo tailscale set --operator=$USER` from the command-line at least once manually.

In Preferences, **Exit node selection** chooses what happens when **Use Exit Node** turns on: the previous exit node, or Tailscale's [recommended exit node](https://tailscale.com/docs/features/exit-nodes/auto-exit-nodes). The self page shows which exit node is in use and can open that node's page.

Building and Installation
-------------------------

### 1. Requirements

Pastikan dependensi sistem sudah terinstall:

* **Go** >= 1.24
* **GTK** >= 4.0
* **Libadwaita** >= 1.6
* **Polkit** (`pkexec`) dan **systemd**

Di Ubuntu / Debian:
```bash
sudo apt update
sudo apt install -y golang libgtk-4-dev libadwaita-1-dev policykit-1
```

### 2. Build & Install ke Menu Aplikasi

Gunakan script `dist.sh` yang sudah disediakan untuk melakukan build dan instalasi otomatis ke direktori user (`~/.local`):

```bash
# 1. Build binary
./dist.sh build v0.21.0-custom

# 2. Install binary, icon, .desktop launcher, dan schema ke ~/.local
./dist.sh install "$HOME/.local"

# 3. Compile GSettings schema & perbarui database menu
glib-compile-schemas "$HOME/.local/share/glib-2.0/schemas"
update-desktop-database "$HOME/.local/share/applications"
gtk-update-icon-cache -f -t "$HOME/.local/share/icons/hicolor"
```

Setelah langkah di atas selesai, aplikasi **Trayscale** akan langsung muncul di menu aplikasi desktop Linux Anda.

---

### Alternative: Flathub & AUR (Upstream)

If you would like, you can also copy the `.desktop` file, the icon, and other pieces of extra metadata into the places that they need to be put to function properly:

* `dev.deedles-trayscale.desktop` -> `$HOME/.local/share/applications/`
* `dev.deedles.Trayscale.png` -> `$HOME/.local/share/icons/hicolor/256x256/apps/`

Note that without copying both of these files into the correct locations, notifications will likely not function correctly in GNOME. Also keep in mind that if the `trayscale` binary is not in your `$PATH` in a way that the desktop environment can locate then the `.desktop` file will not be considered valid. If this is an issue, modify the file manually and change the `Exec=` line to point directly to the binary with an absolute path.

### macOS

First of all, please note that macOS support is not maintained or tested by the primary developer of Trayscale. Also note that the tray icon will _not_ work in macOS.

The simplest way to install Trayscale on macOS is via [Nixpkgs' unstable channel](https://search.nixos.org/packages?channel=unstable&query=trayscale). Note that it will need to be built from source if a binary cache is unavailable and the build can take some time.

To build it locally without Nix, first make sure that you have installed the Go compiler, as well as necessary dependencies. The dependencies are available as the following Homebrew packages:
* `libadwaita`
* `gtk4`
* `gobject-introspection`
* `harfbuzz`

Then, run the following in the root of the repository:

```
$ ./dist.sh build
$ ./dist.sh install-macos
```


Localization
------------

Trayscale can display UI strings in other languages via gettext catalogs under [`po/`](po/). Go-constructed strings use gotext; GtkBuilder labels marked `translatable="yes"` use the same catalogs through libintl. See [`po/README.md`](po/README.md). Prefer `LANGUAGE=ja` (or install the matching system locale); bare `LANG=ja` without a generated libc locale is sanitized so GTK does not hang.

Donate
------

<a href="https://www.buymeacoffee.com/DeedleFake" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/default-green.png" alt="Buy Me A Coffee" style="height: 60px !important;width: 217px !important;" ></a>
