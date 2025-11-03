# Changelog

## [1.9.0](https://github.com/d0ugal/promexporter/compare/v1.8.0...v1.9.0) (2025-11-03)


### Features

* add optional Pyroscope continuous profiling support ([7f7565a](https://github.com/d0ugal/promexporter/commit/7f7565a5ce98b49f5cd87898d69f12ea8e50def4))


### Bug Fixes

* update module github.com/klauspost/compress to v1.18.1 ([4d7c917](https://github.com/d0ugal/promexporter/commit/4d7c917df06768972823cb38c81e3b4574868297))

## [1.8.0](https://github.com/d0ugal/promexporter/compare/v1.7.1...v1.8.0) (2025-11-01)


### Features

* add OpenTelemetry Gin and runtime instrumentation ([27d41de](https://github.com/d0ugal/promexporter/commit/27d41deb567ab1e994a83e75d13b585cc0819f72))
* trigger CI after auto-format workflow completes ([7895323](https://github.com/d0ugal/promexporter/commit/7895323a519e4190fe2f092540cfdee109a8c594))

## [1.7.1](https://github.com/d0ugal/promexporter/compare/v1.7.0...v1.7.1) (2025-10-30)


### Bug Fixes

* update module github.com/prometheus/procfs to v0.19.2 ([80d6eb1](https://github.com/d0ugal/promexporter/commit/80d6eb1510d383606c55d645feba55d1fcbae960))

## [1.7.0](https://github.com/d0ugal/promexporter/compare/v1.6.1...v1.7.0) (2025-10-30)


### Features

* add duplication linter (dupl) to golangci configuration ([35be628](https://github.com/d0ugal/promexporter/commit/35be6281037f937e673e973e56f14c6e3279e19a))
* **ci:** add auto-format workflow ([30212a5](https://github.com/d0ugal/promexporter/commit/30212a516009c92195f8a68e79e04f08d8e05a45))


### Bug Fixes

* update google.golang.org/genproto/googleapis/api digest to ab9386a ([a5c4aa4](https://github.com/d0ugal/promexporter/commit/a5c4aa43749876266f62b58a190e21cd0339f87f))
* update google.golang.org/genproto/googleapis/rpc digest to ab9386a ([9083b5a](https://github.com/d0ugal/promexporter/commit/9083b5aa8c4bca1a5568755ff30c0ad870b682f0))
* update module github.com/gabriel-vasile/mimetype to v1.4.11 ([1e4ce93](https://github.com/d0ugal/promexporter/commit/1e4ce93858f97144375389688d084a2875dfd094))

## [1.6.1](https://github.com/d0ugal/promexporter/compare/v1.6.0...v1.6.1) (2025-10-28)


### Bug Fixes

* update module go.opentelemetry.io/auto/sdk to v1.2.1 ([cab19aa](https://github.com/d0ugal/promexporter/commit/cab19aa72409e6d4e93a9006afe2881c96136229))
* update module go.opentelemetry.io/proto/otlp to v1.8.0 ([2044b49](https://github.com/d0ugal/promexporter/commit/2044b4980726e8cd2ce2106387bdabe9c4cb9223))
* update module google.golang.org/grpc to v1.76.0 ([43b159e](https://github.com/d0ugal/promexporter/commit/43b159e023df56746acf445952394dabc4f5d8c7))
* update opentelemetry-go monorepo to v1.38.0 ([8ef3044](https://github.com/d0ugal/promexporter/commit/8ef3044cd9b484d7fbede726f86eef04d5b23bbb))

## [1.6.0](https://github.com/d0ugal/promexporter/compare/v1.5.1...v1.6.0) (2025-10-28)


### Features

* add dev-tag Makefile target ([b0071ba](https://github.com/d0ugal/promexporter/commit/b0071bad111ebd90c36c9dd322a4f788e29e6011))
* add OpenTelemetry tracing support ([5fd046c](https://github.com/d0ugal/promexporter/commit/5fd046c95182c23f207cd78fa2e20e3671a9e2a1))


### Bug Fixes

* **lint:** update golangci-lint config to use wsl_v5 instead of deprecated wsl ([edbfe99](https://github.com/d0ugal/promexporter/commit/edbfe999e4e88961e1ac00d01c21c54fb0a5a667))
* resolve linting issues ([76a15d7](https://github.com/d0ugal/promexporter/commit/76a15d721844c4b91aaf5628914678b7da37f135))
* **tracing:** correct OTLP endpoint URL handling ([ac677d0](https://github.com/d0ugal/promexporter/commit/ac677d02fdeb4752192b3351eb5045ab9d8f5698))
* update google.golang.org/genproto/googleapis/api digest to 3a174f9 ([92ecf70](https://github.com/d0ugal/promexporter/commit/92ecf7042d2e67b27b16243927ed85b1c0907d06))
* update google.golang.org/genproto/googleapis/rpc digest to 3a174f9 ([3160235](https://github.com/d0ugal/promexporter/commit/316023543bb02acbdb40d86c1cd368f0cd945738))
* update module github.com/bytedance/sonic to v1.14.2 ([9701c42](https://github.com/d0ugal/promexporter/commit/9701c42234661b985a4ff966f64c7c4535292af4))
* update module github.com/cenkalti/backoff/v4 to v5 ([318efac](https://github.com/d0ugal/promexporter/commit/318efac2d02b54662fad6e9c569073cbbac85033))
* update module github.com/grpc-ecosystem/grpc-gateway/v2 to v2.27.3 ([8e4e3d8](https://github.com/d0ugal/promexporter/commit/8e4e3d889bb2971e03b3f32fa7ad042e40f42d98))
* update module github.com/prometheus/common to v0.67.2 ([4345424](https://github.com/d0ugal/promexporter/commit/4345424c7d1fd39773969c9253aaf1b0b22d6ba8))
* update module github.com/ugorji/go/codec to v1.3.1 ([f8c4ac9](https://github.com/d0ugal/promexporter/commit/f8c4ac953b23a56d659f6d94ceba9a57011e786e))

## [1.5.1](https://github.com/d0ugal/promexporter/compare/v1.5.0...v1.5.1) (2025-10-26)


### Bug Fixes

* pass custom version info to web UI ([4ab394f](https://github.com/d0ugal/promexporter/commit/4ab394f8c54c1b0129dd9dd391cda148bfada4cd))

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
