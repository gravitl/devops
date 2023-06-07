## apt/rpm packages (netclient)

* ssh fileserver.netmaker.org (fileserver.clustercat.com)
* cd packages
* git restore .
* git pull
* export VERSION=<netmaker version> # do not include leading v
* export REVISION=<package revision> # revision of package 0, 1, 2 ,etc
* ./apt.builder.sh
* ./rpm.builder.sh
* git restore .
  
## apt/rpm packages (nmcli)
* ssh fileserver.netmaker.org (fileserver.clustercat.com)
* cd packages/nmctl
* git restore .
* git pull
* export VERSION=<netmaker version> # do not include leading v
* export REVISION=<package revision> # revision of package 0, 1, 2 ,etc
* ./apt.builder.sh
* ./rpm.builder.sh
* git restore .

## mac homebrew package (netclient)
* ssh fileserver.netmaker.org (fileserver.clustercat.com)
* cd homebrew
* git pull
* cd build
* export VERSION=<netmaker version> # do not include leading v
* export REVISION=<package revision> # revision of package 0, 1, 2 ,etc
* ./build_tarfiles.sh
* git commit -am "new version"
* git push
  

## mac homebrew package (nmcli)
* ssh fileserver.netmaker.org (fileserver.clustercat.com)
* cd homebrew-nmctl
* git pull
* cd build
* export VERSION=<netmaker version> # do not include leading v
* export REVISION=<package revision> # revision of package 0, 1, 2 ,etc
* ./build_tarfiles.sh
* git commit -am "new version"
* git push

## AUR PACKAGE
* checkout aur.archlinux.org/netclient
* update version in PKGBUILD
* update SHA by running updpkgsums
* run makepkg --printscrinfo > .SCRINFO
* git commit/git push
