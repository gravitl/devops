# Release Process
## netmaker must be released before netclient
### Prerequites
- version numbers updated
- release.md updated to match release
## netmaker-ui-2
### workflows
#### .github/workflows/release.yml --> gravitl/devops/.github/workflows/UIRelease.yml
* on github, select action tabs and run Release Workflow 
  * input version, e.g. v0.21.1
* workflow steps:
  * creates a release branch (release-{version} and tag {version} and push branch/tag (issue: NET-460)
  * creates docker images (amd64, arm64, arm/v7) and uploads to docker hub
  * creates a pull request from release branch to master
## netmaker
### workflows
#### .github/workflows/release.yml --> gravitl/devops/.github/workflows/netmakerRelease.yml
* on github, select action tabs and run Release Workflow 
  * input version, e.g. v0.21.1
* workflow steps:
  * creates a release branch (release-{version} and tag {version} and push branch/tag
  * creates release using goreleaser
    * netmaker : amd64
    * nmctl: linux_amd64, linux_arm64, darwin_amd64, darwin_arm64, freebsd_amd64, windows_amd64
  * creates ce/ee docker images (amd64, arm64, arm/v7) and uploads to docker hub
  * builds nmctl apt/rpm packages amd push to apt.netmaker.io/rpm.netmaker.io (issue: NET-376)
  * creates a pull request from release branch to master
  * copies release assets to fileserver
## netclient
* **do not run until the netmaker release workflow has created the netmaker release branch**
### workflows
#### .github/workflows/release.yml --> gravitl/devops/.github/workflows/netclientRelease.yml
* on github, select action tabs and run Release Workflow 
  * input version, e.g. v0.21.1
* workflow steps:
  * creates a release branch (release-{version} and tag {version}
    * updtes go.mod to to point to github.com/gravitl/netmaker@release_branch
    * push branch/tag
  * assets
    * gorleaser: builds/uploads binaries to release assets
      * netclient: linux_amd64, linux_arm64, linux_arm_5, linux_arm_6, linux_arm_7, linux_mips_hardfloat, linux_mips_softfloat, linux_mipsle_softfloat, linux_mipsle_hardfloat
  * freebsd
    * ssh to freebsd droplets and build linux-freebsd13-amd64, linux-freebsd14-amd64 and upload to release
  * creates netclient docker images (amd64, arm64, arm/v7) and uploads to docker hub
  * builds nmctl apt/rpm packages amd push to apt.netmaker.io/rpm.netmaker.io (issue: NET-376)
  * creates a pull request from release branch to master
### gravitl/devlops/.github/workflows/copyReleasefiles.yml 
* to be run after windows/darwin netclient binaries/packages are manually built and uploaded to release assets
* copies netclient release assets to fileserver
### manual steps
#### windows binaries/packages
* build and upload to release assets
#### darwin binaries/packages
* build and upload to release assets
#### homebrew package (depends on darwin binaries uploaded to release assets)
* see github.com/gravitl/devops/Notes/build_packages.md
#### aur package
* see github.com/gravitl/devops/Notes/build_packages.md
