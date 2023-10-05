# Patching a release 
* not a recommended practise; should only be done in extreme circumstances
## Netmaker-UI
* create a branch from release branch
* submit PR to release branch
* merge PR to release branch
* delete existing release and version tag
* create new release from through GitHub UI  
* run Publish Docker workflow *(run workflow from release branch)*
* create PRs
  * release branch to master
  * release branch to develop
## Netmaker
* create a branch from release branch
* submit PR to release branch
* merge PR to release branch
* update release assets
  * checkout updated release branch
  * retag the release branch with version number
    * e.g. git tag -f v0.21.0; git push
  * run goreleaser release --clean --release-notes release.md
* copy release assets to fileserver
* run Publish Docker workflow  *(run workflow from release branch)*
* manually build nmctl linux/homebrew packages with revision bump
* create PRs
  * release branch to master
  * release branch to develop
## Netclient
* create a branch from release branch
* submit PR to release branch
* merge PR to release branch
* update release assets
  * checkout the release branch
  * update the netmaker import (if netmaker release has changed as well)
    * go get github.com/gravitl/netmaker@v.0.21.0; go mod tidy; git commit -am 'update go.mod'; git push 
  * retag the release branch
     * e.g. git tag -f v0.21.0; git push
  * run goreleaser release --clean --release-notes release.md
  * copy release assets to fileserver
  * manually build darwin/windows binaries/packages and upload to release assets
* run Publish Netclient Docker workflow  *(run workflow from release branch)*
* run Publish Netclient-Userspace Docker workflow  *(run workflow from release branch)*
* manually build netclient and remote-client linux/homebrew packages with revision bump
* create PRs
  * release branch to master
  * release branch to develop
