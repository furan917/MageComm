# Changelog

## [0.1.14](https://github.com/furan917/MageComm/compare/v0.1.13...v0.1.14) (2025-04-04)


### Bug Fixes

* handle escape quotes in message strings ([fdf6bf3](https://github.com/furan917/MageComm/commit/fdf6bf33366d64a577085022f139695dd4ad8a05))

## [0.1.13](https://github.com/furan917/MageComm/compare/v0.1.12...v0.1.13) (2024-03-19)


### Bug Fixes

* Ensure logrus always has the right logfile before logging ([ba34b51](https://github.com/furan917/MageComm/commit/ba34b51a2266490cf20f9cd06812ee1426e42797))

## [0.1.12](https://github.com/furan917/MageComm/compare/v0.1.11...v0.1.12) (2023-11-27)


### Bug Fixes

* Slack command output disable affected full message instead of command output only ([0f1773d](https://github.com/furan917/MageComm/commit/0f1773de706b88305230188fa24f11e7da54e355))

## [0.1.11](https://github.com/furan917/MageComm/compare/v0.1.10...v0.1.11) (2023-11-24)


### Bug Fixes

* Add release pleas comment to version file ([b31cc8d](https://github.com/furan917/MageComm/commit/b31cc8d9c0641913fe134a24f543f6478251fb26))
* add version command properly ([db3df42](https://github.com/furan917/MageComm/commit/db3df42d65a7603f9039fe6a5a5cae4d69ecceeb))
* ensure version output is displayed without comment ([e79e230](https://github.com/furan917/MageComm/commit/e79e230ea70d70069517f79072cc5cba7e59a0da))

## [0.1.10](https://github.com/furan917/MageComm/compare/v0.1.9...v0.1.10) (2023-11-23)


### Bug Fixes

* log formatting + add disallow override file option to config ([544767e](https://github.com/furan917/MageComm/commit/544767eb153363ee7707b6ceb30ca8cf432cfa99))
* magerun command now handles config override correctly ([66eb5ea](https://github.com/furan917/MageComm/commit/66eb5eaad38b6dfaa885061779aa3218ccf162d5))
* small fix for graceful interupts + some whitespacing fixes ([229a649](https://github.com/furan917/MageComm/commit/229a649443ffa8da1dd849ff26231e7ede389885))

## [0.1.9](https://github.com/furan917/MageComm/compare/v0.1.8...v0.1.9) (2023-10-11)


### Bug Fixes

* Clear line properly for waiting indicator ([9bfe765](https://github.com/furan917/MageComm/commit/9bfe7652c394322cc49bc75fd12ad13d343b1b3d))
* error output now properly included in return ([7da6a34](https://github.com/furan917/MageComm/commit/7da6a347aa8f9325bd68c1e234383c25fa1e09ce))
* fix no output returns + functionality to strip content ([869942f](https://github.com/furan917/MageComm/commit/869942f2239373df5605e50005e2880becbb1e1b))
* Fix tests for basic magerun command send ([5928164](https://github.com/furan917/MageComm/commit/59281648548b6a10a97e598ce3b755250feb71aa))

## [0.1.8](https://github.com/furan917/MageComm/compare/v0.1.7...v0.1.8) (2023-09-18)


### Bug Fixes

* Command Escaping added + Correct config file pickup ([5fbd2fe](https://github.com/furan917/MageComm/commit/5fbd2feeb96ccb8076567311689e163bd3e73235))
* Map and Slices can now be read from config file ([839b3c6](https://github.com/furan917/MageComm/commit/839b3c6c40171a0a22008b84f65216871220889c))

## [0.1.7](https://github.com/furan917/MageComm/compare/v0.1.6...v0.1.7) (2023-08-31)


### Bug Fixes

* CGO ENV setting ([7e0b462](https://github.com/furan917/MageComm/commit/7e0b462062386bdbc1c9027af190fcc45b0bbd6c))

## [0.1.6](https://github.com/furan917/MageComm/compare/v0.1.5...v0.1.6) (2023-07-20)


### Bug Fixes

* Allow limitation of listener queues to avoid abuse or mistakes via user misspelling ([6bb74ba](https://github.com/furan917/MageComm/commit/6bb74ba93d3e8b2f5b36b9b6c856965d660c15e7))

## [0.1.5](https://github.com/furan917/MageComm/compare/v0.1.4...v0.1.5) (2023-06-21)


### Bug Fixes

* Added configuration overriding and made sweeping fixes to configuration settings ([b56fa10](https://github.com/furan917/MageComm/commit/b56fa10e50c486555c029577ad52d6e2cdd9b43b))

## [0.1.4](https://github.com/furan917/MageComm/compare/v0.1.3...v0.1.4) (2023-06-15)


### Bug Fixes

* improved  n98 command handler and updated slack output notifier ([5b1d7ec](https://github.com/furan917/MageComm/commit/5b1d7ec767dfb945722acecaeb2fe8371bab8d8d))

## [0.1.3](https://github.com/furan917/MageComm/compare/v0.1.2...v0.1.3) (2023-06-14)


### Bug Fixes

* Added way to listen to magerun_output queue, added way to exit listening on output return ([879a106](https://github.com/furan917/MageComm/commit/879a106ac22a905349a8e61d40fa621388dc936d))

## [0.1.2](https://github.com/furan917/MageComm/compare/v0.1.1...v0.1.2) (2023-06-14)


### Bug Fixes

* Updating CTX wait timeout of SQS ([99398e8](https://github.com/furan917/MageComm/commit/99398e805f16a8346b4c4bd4c5f36e03998131e5))

## [0.1.1](https://github.com/furan917/MageComm/compare/v0.1.0...v0.1.1) (2023-04-30)


### Bug Fixes

* readme showed deploy command existed ([65cc4a7](https://github.com/furan917/MageComm/commit/65cc4a7e0fd68a143feba505bd49babc2281ba7a))

## [0.1.0](https://github.com/furan917/MageComm/compare/v0.0.15...v0.1.0) (2023-04-30)


### Features

* Fixed SQS and added slack notification ability ([3e4a869](https://github.com/furan917/MageComm/commit/3e4a869aaf5828f024707e39ec10d2c187c69836))

## [0.0.15](https://github.com/furan917/MageComm/compare/v0.0.14...v0.0.15) (2023-04-26)


### Bug Fixes

* SQS was set to short polling, updated to wait 60s ([eb4be08](https://github.com/furan917/MageComm/commit/eb4be08a63cb3cfd0eb13db6a224281b5ae2e3af))

## [0.0.14](https://github.com/furan917/MageComm/compare/v0.0.13...v0.0.14) (2023-04-25)


### Bug Fixes

* Correct tests to work with new functionality ([5b82107](https://github.com/furan917/MageComm/commit/5b82107be816ccf9534d20b11b90cbf8f2b012ad))

## [0.0.13](https://github.com/furan917/MageComm/compare/v0.0.12...v0.0.13) (2023-04-25)


### Bug Fixes

* Include Restricted & Required argument configuration ([ebd175b](https://github.com/furan917/MageComm/commit/ebd175b51bb53367aabb0712409d0ea7bb9ed110))

## [0.0.12](https://github.com/furan917/MageComm/compare/v0.0.11...v0.0.12) (2023-04-20)


### Bug Fixes

* ref_name for autodeployments ([40bbfd2](https://github.com/furan917/MageComm/commit/40bbfd22e0faa01e5384140107b3ce61cf8da6e7))

## [0.0.11](https://github.com/furan917/MageComm/compare/v0.0.10...v0.0.11) (2023-04-20)


### Bug Fixes

* Correct release.yml behaviour and allow manual running of QOL actions ([644dfaa](https://github.com/furan917/MageComm/commit/644dfaa666385c967709cd61c8ba75ae5d2bfe13))

## [0.0.10](https://github.com/furan917/MageComm/compare/v0.0.9...v0.0.10) (2023-04-20)


### Bug Fixes

* Allow ReleasePlease to create build assets for releases ([d7f44ba](https://github.com/furan917/MageComm/commit/d7f44bac257e32dbd280750261119c277e961ff8))

## [0.0.9](https://github.com/furan917/MageComm/compare/v0.0.8...v0.0.9) (2023-04-20)


### Bug Fixes

* SQS queues did not work correctly with CorrelationIDs ([cdf48cb](https://github.com/furan917/MageComm/commit/cdf48cbe93157ad97da9e0cce8377005a80fc591))
