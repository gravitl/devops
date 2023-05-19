# Upgrading Netmaker SaaS

How to upgrade staging and prod environments:

- go over the upgrade process for AMB [documented here](https://github.com/gravitl/account-management-backend/blob/main/README.md)
- go over the upgrade process for NM operator [documented here](https://github.com/gravitl/netmaker-operator/blob/main/README.md)
- go over the upgrade process for NMUI and AMUI [documented here](https://github.com/gravitl/saas/blob/master/README.md)
- for netmaker saas images:
  - checkout GRA-1298 branch and pull
  - change ee/types.go api_endpoint const to "https://api.staging.accounts.netmaker.io/api/v1/license/validate"
  - `docker build -t gravitl/netmaker:saas . --push`
  - `docker build -t gravitl/netmaker:saas-ee --build-arg tags=ee . --push`
  - change ee/types.go api_endpoint const to "https://api.accounts.netmaker.io/api/v1/license/validate" (remove staging)
  - `docker build -t gravitl/netmaker:saasprod . --push`
  - `docker build -t gravitl/netmaker:saasprod-ee --build-arg tags=ee . --push`
- update instances, for every tenant namespace, with: `kubectl delete po nm-app-<id>-server-0 -n <id>`
