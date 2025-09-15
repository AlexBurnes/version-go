# Linux Installation

## Method 1: Using Simple Installer Scripts

Download and run the installer script directly:

```bash
# Download and install to /usr/local/bin (requires sudo)
wget -O - https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh | sh

# Download and install to custom directory
INSTALL_DIR=/opt/version wget -O - https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh | sh

# Or download first, then install
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version-1.0.0-linux-amd64-install.sh
chmod +x version-1.0.0-linux-amd64-install.sh
./version-1.0.0-linux-amd64-install.sh /usr/local/bin
```

## Method 2: Manual Installation from Archive

Download and extract the archive manually:

```bash
# Download and extract archive
wget https://github.com/AlexBurnes/version-go/releases/download/v1.0.0/version_1.0.0_linux_amd64.tar.gz
tar -xzf version_1.0.0_linux_amd64.tar.gz
cd version_1.0.0_linux_amd64

# Install using the included install.sh
./install.sh /usr/local/bin

# Or install to user directory
./install.sh ~/.local/bin
```

## Windows Installation

Windows users can use Scoop:
```bash
scoop install version
```