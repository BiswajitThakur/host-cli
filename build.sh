#!/bin/bash

source config.env

if [ ! -f /usr/bin/go ] && [ ! -f /bin/go ] && [ ! -f /usr/local/bin/go ]; then
    sudo apt update
    sudo apt install golang
fi

if [ $OS == "windows" ]; then
    GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o $PKG_NAME.exe
    exit 0
fi

function deb_control {
    local ctrl=$PKG_NAME/DEBIAN/control
    echo "Package: $PKG_NAME" > $ctrl
    echo "Version: $VERSION" >> $ctrl
    echo "Maintainer: Biswajit Thakur" >> $ctrl
    echo "Architecture: $ARCH" >> $ctrl
    echo "Homepage: $HOME_PAGE" >> $ctrl
    echo "Installed-Size: $(du -s $PKG_NAME | awk '{print $1}')" >> $ctrl
    echo "Description: CLI based Ads & Websites Blocker application." >> $ctrl
    chmod +x $ctrl
}

function deb_preinst {
    local preinst=$PKG_NAME/DEBIAN/preinst
    echo "#!/bin/bash" > $preinst
    echo "if [ -L /usr/bin/$PKG_NAME ] || [ -f /usr/bin/$PKG_NAME ]; then" >> $preinst
    echo "    sudo rm /usr/bin/$PKG_NAME" >> $preinst
    echo "    echo \"Removed old /usr/bin/$PKG_NAME\"" >> $preinst
    echo "fi" >> $preinst
    echo "if [ -d /usr/share/$PKG_NAME ]; then" >> $preinst
    echo "    rm -rf /usr/share/$PKG_NAME" >> $preinst
    echo "    echo \"Removed old /usr/share/$PKG_NAME\"" >> $preinst
    echo "fi" >> $preinst
    chmod +x $preinst
}

function build_deb {
    if [ ! -d ./$PKG_NAME/DEBIAN ]; then
        mkdir -p ./$PKG_NAME/DEBIAN
    fi
    if [ ! -d ./$PKG_NAME/usr/bin ]; then
        mkdir -p ./$PKG_NAME/usr/bin
    fi
    if [ ! -d ./$PKG_NAME/usr/share/$PKG_NAME ]; then
        mkdir -p ./$PKG_NAME/usr/share/$PKG_NAME
    fi
    # local allowPath="./$PKG_NAME/usr/share/$PKG_NAME/allow.txt"
    # local blockPath="./$PKG_NAME/usr/share/$PKG_NAME/block.txt"
    # local sourcePath="./$PKG_NAME/usr/share/$PKG_NAME/sources.txt"
    # touch $allowPath
    # touch $blockPath
    # touch ./$PKG_NAME/usr/share/$PKG_NAME/redirect.txt
    # echo "https://adaway.org/hosts.txt" > $sourcePath
    # echo "https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext" >> $sourcePath
    # echo "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts" >> $sourcePath
    deb_preinst
    if [ -f /usr/bin/go ] || [ -f /bin/go ] || [ -f /usr/local/bin/go ]; then
        GOOS=$OS GOARCH=$ARCH go build -ldflags="-s -w" -o ./$PKG_NAME/usr/bin/$PKG_NAME
        # go build -ldflags="-s -w" -o ./$PKG_NAME/usr/bin/$PKG_NAME
        deb_control
        dpkg-deb --build $PKG_NAME
        rm -rf $PKG_NAME
        local u="_"
        mv $PKG_NAME.deb $PKG_NAME$u$VERSION$u$ARCH.deb
    fi
}

build_deb

