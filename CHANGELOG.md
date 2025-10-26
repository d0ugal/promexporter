# Changelog

## [1.5.0](https://github.com/d0ugal/promexporter/compare/v1.4.1...v1.5.0) (2025-10-26)


### Features

* add SkipVersionInfo option to disable automatic version info setting ([8bcc33f](https://github.com/d0ugal/promexporter/commit/8bcc33f96a76de138c149993165f83b0207de0a6))
* add WithVersionInfo method to accept custom version information ([bab712c](https://github.com/d0ugal/promexporter/commit/bab712c0c0a80078aabc76e9939c3bed4ba462b7))


### Bug Fixes

* Add a warning log when falling back ([0eefae5](https://github.com/d0ugal/promexporter/commit/0eefae5d9c7fedd644645046fa2790c3b1593d60))
* linting ([f06e445](https://github.com/d0ugal/promexporter/commit/f06e44502f22f37fbb829bc83e11eac4b7605eea))
* resolve CI linting issues ([6a58262](https://github.com/d0ugal/promexporter/commit/6a5826289f85f8749341b1e70bc12c2b591dc9d9))
* update module github.com/prometheus/procfs to v0.19.1 ([c7ccb84](https://github.com/d0ugal/promexporter/commit/c7ccb84b69ba5a08c0b837c5c77c6820107518c3))

## [1.4.1](https://github.com/d0ugal/promexporter/compare/v1.4.0...v1.4.1) (2025-10-25)


### Bug Fixes

* linting ([c143a20](https://github.com/d0ugal/promexporter/commit/c143a20b95aa83cc87191e80922db69bf9b48e37))
* update module github.com/prometheus/procfs to v0.19.0 ([0801fe0](https://github.com/d0ugal/promexporter/commit/0801fe08394bc4915654cf85a9c4dafde2f2a16a))

## [1.4.0](https://github.com/d0ugal/promexporter/compare/v1.3.1...v1.4.0) (2025-10-24)


### Features

* add custom template rendering system ([6086ba2](https://github.com/d0ugal/promexporter/commit/6086ba2c5e66b062ee3f3816f932e326dc8c9fa6))
* add optional web UI and health endpoint configuration ([6f2da4f](https://github.com/d0ugal/promexporter/commit/6f2da4f0dc7a0e0ec982b5201c2c3ddabb070836))
* add optional web UI and health endpoint configuration ([4fdfca4](https://github.com/d0ugal/promexporter/commit/4fdfca450361e21b091598d48f4076b4ee128af2))
* improve configuration display formatting ([c6d2716](https://github.com/d0ugal/promexporter/commit/c6d27166de7e329fecc0a5c7bcb8d3021247430d))
* improve nested map formatting ([2ee04e0](https://github.com/d0ugal/promexporter/commit/2ee04e02759e6ce30f6c27c05dc994a96b40b3c0))
* support custom HTML fragments for config rendering ([2b3a5af](https://github.com/d0ugal/promexporter/commit/2b3a5afd01ff912d7b7ed70b5b20c0ff5a67e5f3))


### Bug Fixes

* handle array of maps in template formatting ([b0abbf0](https://github.com/d0ugal/promexporter/commit/b0abbf02cd29d4d43acb9481ebbae11bedb7eb61))
* pass custom HTML as template.HTML to prevent escaping ([3845e5b](https://github.com/d0ugal/promexporter/commit/3845e5b3c25bf1067e29ff4748394e7cd4e865a3))
* register safeHTML template function ([c98ba6e](https://github.com/d0ugal/promexporter/commit/c98ba6eb939b6dc4908dfcd4cfd8c1e90c5bd300))
* use safeHTML instead of html filter for custom HTML fragments ([80b7744](https://github.com/d0ugal/promexporter/commit/80b7744c64cf390d91ac79c35a4d33f442a50b0c))

## [1.3.1](https://github.com/d0ugal/promexporter/compare/v1.3.0...v1.3.1) (2025-10-24)


### Bug Fixes

* use concrete type check instead of interface assertion ([271b007](https://github.com/d0ugal/promexporter/commit/271b007fb44b9911ea81c1ef3350ae4980d1188f))
* use interface method check instead of type assertion ([789c856](https://github.com/d0ugal/promexporter/commit/789c8565a5af4e3dc04f42adcae29d73d06773d7))

## [1.3.0](https://github.com/d0ugal/promexporter/compare/v1.2.0...v1.3.0) (2025-10-24)


### Features

* improve configuration display with type-based sensitivity detection ([a51340c](https://github.com/d0ugal/promexporter/commit/a51340ce7d5c07068e97072f911fbeb39c3d3210))


### Bug Fixes

* use concrete type check instead of interface assertion ([271b007](https://github.com/d0ugal/promexporter/commit/271b007fb44b9911ea81c1ef3350ae4980d1188f))

## [1.2.0](https://github.com/d0ugal/promexporter/compare/v1.1.0...v1.2.0) (2025-10-24)


### Features

* support custom configuration types in app and server ([deb4530](https://github.com/d0ugal/promexporter/commit/deb4530423cb46b2d0ce2e7310134f2507b6219e))

## [1.1.0](https://github.com/d0ugal/promexporter/compare/v1.0.2...v1.1.0) (2025-10-23)


### Features

* add SensitiveString type for explicit sensitive configuration handling ([2984151](https://github.com/d0ugal/promexporter/commit/298415176ba2b8fdb08b8bb5c249aa91f079c409))


### Bug Fixes

* actually use ConfigDisplay in getConfigData method ([79809a4](https://github.com/d0ugal/promexporter/commit/79809a40c05b1e41e9a0190540f05da805f3b82e))
* format code with go fmt ([611f9bd](https://github.com/d0ugal/promexporter/commit/611f9bd3136116d9b0beed3ad2ef547610821a9a))
* resolve remaining wsl linting issues ([fb25803](https://github.com/d0ugal/promexporter/commit/fb25803e3b4de5cd590d493824602c2f1aee352d))
* resolve wsl linting issues in sensitive.go ([d9cb539](https://github.com/d0ugal/promexporter/commit/d9cb53912f827424d2c5d63e9a060e26c4e7ea90))
* update golangci-lint config for Go version compatibility ([905faa3](https://github.com/d0ugal/promexporter/commit/905faa304747fd4ff112caf3f887459cc91cbfd0))

## [1.0.2](https://github.com/d0ugal/promexporter/compare/v1.0.1...v1.0.2) (2025-10-23)


### Bug Fixes

* update module github.com/prometheus/procfs to v0.18.0 ([efc02be](https://github.com/d0ugal/promexporter/commit/efc02be5989fbd5297925b5477952c069b48ad6e))

## [1.0.1](https://github.com/d0ugal/promexporter/compare/v1.0.0...v1.0.1) (2025-10-19)


### Bug Fixes

* resolve funcorder linting issues ([c7ed478](https://github.com/d0ugal/promexporter/commit/c7ed478f9cdad858799438158d37523d9aeb38f3))
* update linting configuration and format code ([fb2b3c3](https://github.com/d0ugal/promexporter/commit/fb2b3c3a7be2196c5c77aaf1899fb517482aec1d))

## 1.0.0 (2025-10-18)


### Features

* add comprehensive renovate configuration in JSON5 format ([d430d39](https://github.com/d0ugal/promexporter/commit/d430d397f3758d6b61ba731bf6e5b14b5c630a0d))
* add Go runtime metrics to promexporter library ([d261a5b](https://github.com/d0ugal/promexporter/commit/d261a5b71ad1e4f3c6fbb9d9ac20483f4a25971b))
* configure release-please to start with version 0.1.0 ([8ef1fa8](https://github.com/d0ugal/promexporter/commit/8ef1fa822c0c504ce3b8887b67dd6a7c6dc2ac90))
* initial promexporter library implementation ([c6c2ec1](https://github.com/d0ugal/promexporter/commit/c6c2ec110fc73d1cd5223161dc81d6e8109e5145))


### Bug Fixes

* restore proper golangci-lint config from working exporters ([6788c70](https://github.com/d0ugal/promexporter/commit/6788c701090db248ea8b26f007526bbc051a7d6e))
* simplify golangci-lint config to exclude typecheck yaml errors ([68b480e](https://github.com/d0ugal/promexporter/commit/68b480eb01fe9dc4b810cf63580855293e1a398a))
* simplify golangci-lint config to work in CI environment ([0a66a88](https://github.com/d0ugal/promexporter/commit/0a66a8873e0ef16a170d3e5492710755996dfc29))
* update golangci-lint config to match other exporters with typecheck exclusion ([27538e4](https://github.com/d0ugal/promexporter/commit/27538e4568dac835ad121b9734bf527b4d8f3b2a))
