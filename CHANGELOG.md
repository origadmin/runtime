# Changelog

## [0.2.9](https://github.com/origadmin/runtime/compare/v0.2.8...v0.2.9) (2025-11-03)


### Features

* **runtime:** Remove storage system design document as it's outdated ([80d8721](https://github.com/origadmin/runtime/commit/80d8721fe58da4b67f0bccd6e0be3418ab1beeb9))

## [0.2.8](https://github.com/origadmin/runtime/compare/v0.2.7...v0.2.8) (2025-11-03)


### Features

* **api:** Added protobuf definition and validation code for middleware and transport layer configurations ([a96619f](https://github.com/origadmin/runtime/commit/a96619f3ad28aa312d1404f703421b25f8af5b81))
* **broker:** Add detailed broker configurations for multiple message queue systems ([6b4c05f](https://github.com/origadmin/runtime/commit/6b4c05fea8bbf2a522ae0dde6d0060c875e77bbd))
* **broker:** Extend broker configuration with multiple message queue implementations and validation ([5184786](https://github.com/origadmin/runtime/commit/5184786c8ce1158854c3beead16a2fed908dfa61))
* **runtime:** Add broker list and enhance storage proto definitions with OSS lifecycle rules and Redis config options ([46e3545](https://github.com/origadmin/runtime/commit/46e35450d0fe4f5d586c586b6d904e57d0528378))
* **runtime:** Add quick start app example with bootstrap configuration ([954e0c0](https://github.com/origadmin/runtime/commit/954e0c0fa67001faa24b46c45d808f9f216a5dcf))
* **selector:** Added SelectorConfig protocol buffer definition ([aefd179](https://github.com/origadmin/runtime/commit/aefd17990eca613dbce8d208ba98df6f4674e059))

## [0.2.7](https://github.com/origadmin/runtime/compare/v0.2.6...v0.2.7) (2025-10-28)


### Features

* **discovery:** Add Discoveries message type to support multiple discovery configurations ([073409d](https://github.com/origadmin/runtime/commit/073409d4983e7a1add5c0b85f4e4020a9003c68d))

## [0.2.6](https://github.com/origadmin/runtime/compare/v0.2.5...v0.2.6) (2025-10-13)


### Features

* Add go.work file for monorepo management ([6bb255d](https://github.com/origadmin/runtime/commit/6bb255d1d44c1bb264d267971682411a435a45c5))
* **api:** add Apollo and Consul source configurations ([a1dbc62](https://github.com/origadmin/runtime/commit/a1dbc62fee189a89ec3f5aa1d5a3e8388ac90d9d))
* **api:** add configuration protobufs for client and registry ([6f75ef2](https://github.com/origadmin/runtime/commit/6f75ef29a8f634a32996664d791470c4b999e194))
* **api:** Add protobuf definitions for app, optimize middleware and transport selector ([dc6fa23](https://github.com/origadmin/runtime/commit/dc6fa23ffd5e9a9309e70fa97f8a506a36d56acc))
* **apierrors:** add new error categories and messages ([ab3f83a](https://github.com/origadmin/runtime/commit/ab3f83aa60c05b514bcaffe6510cfa28ee3cf0cc))
* **bootstrap:** Add app config support and refactor provider to container ([b05f17e](https://github.com/origadmin/runtime/commit/b05f17e9d3b86be8862f025d6011a6be9ea39fc8))
* **bootstrap:** add Bootstrapper interface for component provider and config access ([48fa0e7](https://github.com/origadmin/runtime/commit/48fa0e7bd5f5d6b839e4666a266bbef0077799e7))
* **bootstrap:** add dynamic component creation and registration mechanism ([b038114](https://github.com/origadmin/runtime/commit/b038114eb4456f2b3fcebd66fb08d1383d1e7a2f))
* **bootstrap:** enhance config loading with custom config and transformer support ([ac1dca7](https://github.com/origadmin/runtime/commit/ac1dca7ddb77ef7c27f4a36de0711b26129a5ef9))
* **bootstrap:** implement new component provider and decoder interfaces ([bccb07f](https://github.com/origadmin/runtime/commit/bccb07f3fe911c50b317d0bbf1114348a986a727))
* **bootstrap:** Implements application startup configuration and component initialization ([7c52593](https://github.com/origadmin/runtime/commit/7c52593a9833276fd8338e010271794d99198b81))
* **cache:** add cleanup interval for memory cache ([e6e2953](https://github.com/origadmin/runtime/commit/e6e2953ce10f76f137126fc75463e03df719e03b))
* **config:** add custom config parser example with bootstrap and config files ([31c1336](https://github.com/origadmin/runtime/commit/31c1336e54ff22f13a6a5e46752fafe810137632))
* **config:** Add helper package for config file operations and update test configs ([cc675da](https://github.com/origadmin/runtime/commit/cc675da30e0c2c782c88119ce098a0438f58f7ae))
* **config:** add multi-format test configs and enhance config loading tests ([f497994](https://github.com/origadmin/runtime/commit/f4979949fd320f86c767236066afdf1b60d99eea))
* **config:** add multi-format test configs and enhance config loading tests ([2d82e5d](https://github.com/origadmin/runtime/commit/2d82e5d855e9c3dd5134077e574d717a206ee827))
* **config:** add priority field to SourceConfig and update environment variable prefix handling ([52977b7](https://github.com/origadmin/runtime/commit/52977b76b577469dc8f2b47ee787ada9f93e906c))
* **config:** add RegisterSourceFactory and RegisterSourceFunc for flexible source registration ([43aa108](https://github.com/origadmin/runtime/commit/43aa108d0e3979fdc7a968cbf79b0bdfdf7ac6bb))
* **config:** assign default priorities for common config sources ([85de18e](https://github.com/origadmin/runtime/commit/85de18ee15d51140e575204ed7ad7439be9a5e61))
* **config:** enhance TLS configuration and remove redundant use_tls fields ([48d3493](https://github.com/origadmin/runtime/commit/48d3493f607e44e37b9259896e0146a8829654f1))
* **config:** implement Decode method in Resolved interface ([4c9be3d](https://github.com/origadmin/runtime/commit/4c9be3d1011d2adc9c5e2be83763f8faacf027a6))
* **config:** refactor config decoder to use BaseDecoder embedding ([c266418](https://github.com/origadmin/runtime/commit/c266418c1d54bbb90ff8d5f0e5a19d1b92788698))
* **config:** refactor source configuration and add new source types ([e088aad](https://github.com/origadmin/runtime/commit/e088aad8c567940f9ac38effc006e66d33ad3147))
* **config:** Support for decoding configurations using protojson ([50a5ab7](https://github.com/origadmin/runtime/commit/50a5ab7501996b45f56bf3c37c75615d992365b5))
* **config:** Update the configuration and implementation of the local registration example ([ee3ff94](https://github.com/origadmin/runtime/commit/ee3ff94259bc3bf11ad35b267d4c80b4b51f9805))
* **cors:** Refactor CORS configuration with enhanced security options and better documentation ([b9bd552](https://github.com/origadmin/runtime/commit/b9bd5520cd50a95613f549b8d33acd4d5edd4b85))
* **dev:** implement subtree push workflow for monorepo management ([fcebfb7](https://github.com/origadmin/runtime/commit/fcebfb797655e17ae32abe19cdb66851936c0fb6))
* **devops:** Add unified build, lint, pre-commit, and CI/CD configurations ([4248399](https://github.com/origadmin/runtime/commit/4248399e476ef74c61e63555a7153b4000fcc755))
* **discovery:** add Endpoint and Selector protobuf definitions with validation ([e6b4763](https://github.com/origadmin/runtime/commit/e6b47631f5237d533ec87ec5f80cbb80ba4bbe85))
* **errors:** add error handling and conversion package ([8f2c956](https://github.com/origadmin/runtime/commit/8f2c9567dc1ef4aece86604584daf091e77c342d))
* **errors:** add metadata classification helpers for error handling ([680c150](https://github.com/origadmin/runtime/commit/680c15067a2113ab2bbbce77c1fc84cda4e848a9))
* **errors:** Add structured error handling with Kratos integration and metadata support ([d629e31](https://github.com/origadmin/runtime/commit/d629e31641fd07af540b6b41a41bdde25ce4cd63))
* **errors:** enhance error handling with metadata and chaining support ([54d38c1](https://github.com/origadmin/runtime/commit/54d38c1d346b8c1b131e35cf859b9da6b274e1bb))
* **http:** Add CORS support and reorder server initialization options ([3f33601](https://github.com/origadmin/runtime/commit/3f33601556f87569a00bc6addc5d618471c983c2))
* **http:** Add pprof support and refactor CORS configuration in HTTP server ([4a8cbdc](https://github.com/origadmin/runtime/commit/4a8cbdc58855180756e292f2efaf3ebcfc0bb3ee))
* **loader:** update config types and transport configurations ([11d3971](https://github.com/origadmin/runtime/commit/11d397115cd093ad7284380ccbfcf13908130306))
* **logger:** Add context support and utility functions for Logger ([b98d1b9](https://github.com/origadmin/runtime/commit/b98d1b974f32bc99687f57c4e687994329f2a0d5))
* **meta:** implement directory indexing and improve file handling ([499d6ab](https://github.com/origadmin/runtime/commit/499d6abdb99a1d8e9ed5ec196a64e4394032e37a))
* **middleware:** add logging middleware support ([43647ba](https://github.com/origadmin/runtime/commit/43647ba55eb00bc6270b3547fd77176af0845d91))
* **middleware:** Add Metadata configuration and refactor MiddlewareConfig to use optional fields ([84448f0](https://github.com/origadmin/runtime/commit/84448f04dd4d1c04374889f3eb0cd274fd81d56a))
* **middleware:** Add middleware configuration support to bootstrap system ([f4c3afa](https://github.com/origadmin/runtime/commit/f4c3afacfd9614b61469928d85b905864f176d68))
* **middleware:** Add middleware provider interface and refactor config decoding ([7b0af5c](https://github.com/origadmin/runtime/commit/7b0af5cc54692b8f4f9271b87971a5dbc04ccf2e))
* **middleware:** Add name field to MiddlewareConfig and enhance CORS middleware ([b6d0da4](https://github.com/origadmin/runtime/commit/b6d0da4bc5819dd9a6f8210f8960072a5183aa7f))
* **middleware:** Add name field to MiddlewareConfig and enhance CORS middleware ([e7c4f5f](https://github.com/origadmin/runtime/commit/e7c4f5f92bd048665be207a7471b695a19a8f40b))
* **middleware:** Add type alias for Option and enhance test coverage ([ee61a0d](https://github.com/origadmin/runtime/commit/ee61a0dcee9300d0ffe0f3baccc63cdca6ab1ea9))
* **middleware:** Refactor options package usage and add integration tests ([047ab44](https://github.com/origadmin/runtime/commit/047ab4420f70653730c1541a549cdc1cb0656bca))
* **optionutil:** Add context support and new ApplyNew function for option handling ([3f60b03](https://github.com/origadmin/runtime/commit/3f60b03c159ec14737ab21a77f274f1d1f28415e))
* **optionutil:** Add ValueOr, ApplyContext, If and Group utility functions and rename opt to ctx for clarity ([7f1bb76](https://github.com/origadmin/runtime/commit/7f1bb764aa204da83dffa147ad617b0d221af244))
* Recreate root go.mod file ([7788964](https://github.com/origadmin/runtime/commit/7788964e26a5f9b8627ca5ece1c52f5316ddce19))
* **registry:** implement registry error handling ([d7f5b7b](https://github.com/origadmin/runtime/commit/d7f5b7bab0dd135f228d49b0c3a35feffea4ecf0))
* **runtime:** add error handling and improve registry implementation ([6a45d51](https://github.com/origadmin/runtime/commit/6a45d514e0ae87c2130691e5701e2335d919d0d3))
* **runtime:** add optionutil package and move option-related functions ([78c3ff0](https://github.com/origadmin/runtime/commit/78c3ff065a5f981027bdb57d7e3a09e9692ff381))
* **runtime:** add pagination support and remove customize config ([9ae94ac](https://github.com/origadmin/runtime/commit/9ae94ac1a4694d4fa24e8728d45f842716252bb2))
* **runtime:** Add proto definition for App and refactor component provider with improved logging ([222baf4](https://github.com/origadmin/runtime/commit/222baf44e9214ab41c2f3c6814f06eb5045534fa))
* **runtime:** Add server-side and client-side middleware acquisition interfaces ([08413c2](https://github.com/origadmin/runtime/commit/08413c221f220d25af7382aca8af1b291b75ab3e))
* **runtime:** Add StructuredConfig interface for type-safe configuration decoding ([c299341](https://github.com/origadmin/runtime/commit/c2993415fcc36a996a1a36fda88d38b5ca6adb0d))
* **runtime:** add support for decoding string to durationpb.Duration ([2258141](https://github.com/origadmin/runtime/commit/225814123d17b23d8bab76aa92ade208dd8faf08))
* **runtime:** add WithLogger method and update logging ([f85f0d6](https://github.com/origadmin/runtime/commit/f85f0d616e519479718e1cb24341de42a0d47853))
* **runtime:** implement bootstrap configuration loading and sorting ([34a4e64](https://github.com/origadmin/runtime/commit/34a4e648ccbb3042adb26bd1f3bae05fb3ee8260))
* **runtime:** implement config-driven initialization and add app metadata support ([4f98ff2](https://github.com/origadmin/runtime/commit/4f98ff236f0f2a805cd670582e632d60bacbf713))
* **runtime:** implement new runtime interface and add layout package ([6f4f57c](https://github.com/origadmin/runtime/commit/6f4f57c0ecd9c7625124cd650edd9858d1b3b029))
* **runtime:** introduce WithAppInfo option and enhance configuration flexibility ([72e01ab](https://github.com/origadmin/runtime/commit/72e01ab8d6c3c9cff65da50b6b77e80c603dc7d7))
* **runtime:** Introducing a runtime option aggregation mechanism ([b7c38cf](https://github.com/origadmin/runtime/commit/b7c38cf4c9cd988af10e14f38dcb9ef8435358e9))
* **runtime:** Refactoring runtime packages ([86eb1af](https://github.com/origadmin/runtime/commit/86eb1af510721ca78c1eff9520cebdc0ee840af1))
* **runtime:** Support custom names for middlewares with fallback to type ([d42d473](https://github.com/origadmin/runtime/commit/d42d473be5ecba9dcfc0cf822276915478f0b9fe))
* **selector:** Enhance selector proto with include/exclude fields and update middleware logic ([f2232f9](https://github.com/origadmin/runtime/commit/f2232f9ef375451ecceaca34909681dfb33dc0da))
* **server:** Add service registration mechanism and example services ([5edf3c5](https://github.com/origadmin/runtime/commit/5edf3c580baaade5de37cc502f7c989943d60733))
* **storage:** add delete and rename functionality ([132b972](https://github.com/origadmin/runtime/commit/132b9729171c09987d0dc1bcdd7ee91cbfd501c0))
* **storage:** add protobuf definitions for file system operations ([6d0c1e6](https://github.com/origadmin/runtime/commit/6d0c1e6dd384321a275884cc956f6e12fa363a4d))
* **storage:** implement advanced breadcrumb truncation and normalize paths ([c5eaada](https://github.com/origadmin/runtime/commit/c5eaadab47f9e5036e3a7e04b6dc66734e0b6408))
* **storage:** implement chunked file storage and directory indexing ([5dae991](https://github.com/origadmin/runtime/commit/5dae9915003ea4053c5c99f6162ae7cb63c66edf))
* **storage:** implement content streaming writer and optimize file metadata handling ([2536dd4](https://github.com/origadmin/runtime/commit/2536dd483d41856ab097f5f1cc20e4102a2355ad))
* **storage:** implement directory and file creation with metadata handling ([bdc7042](https://github.com/origadmin/runtime/commit/bdc70428e6e603adbda29ac4f4c93fc4c68e8451))
* **storage:** implement local file system storage and enhance UI ([9e27db7](https://github.com/origadmin/runtime/commit/9e27db7510dc84a86890e32c537f4e0fe69d8c6a))
* **storage:** implement new file metadata and blob storage interfaces ([530c9b3](https://github.com/origadmin/runtime/commit/530c9b368276252e1dea59c7a56fb5a1028214f3))
* **storage:** implement WriteContent method in content assembler ([0306b41](https://github.com/origadmin/runtime/commit/0306b414c9cfb621de67f111a93bf5e34a37ec17))
* **storage:** refactor Meta service for improved file handling ([7b77dc8](https://github.com/origadmin/runtime/commit/7b77dc88eb495c10eeaa4b4c732d7b11b2f3b541))
* **storage:** refactor storage configuration and implement cache ([894430f](https://github.com/origadmin/runtime/commit/894430f924979ee2b5b5db9b4cce0702b26f3854))
* **transport:** Add CORS support to HTTP server configuration ([400737e](https://github.com/origadmin/runtime/commit/400737e9cd4e93164a58d625050c1b2cdc10578b))
* **transport:** Add network field to GrpcServerConfig and update related protobuf definitions ([364a618](https://github.com/origadmin/runtime/commit/364a6189376beecff789547ea32079d2fa7cfeeb))
* **transport:** Enhanced HTTP transport configuration support ([aa44bc7](https://github.com/origadmin/runtime/commit/aa44bc739a1594bd1be9ea976b3d9ea6d27cec63))


### Bug Fixes

* Remove go.work from .gitignore ([d405949](https://github.com/origadmin/runtime/commit/d4059490a17e46a8f7acb6c68176b00026296bba))
* **runtime:** enforce AppInfo validation and update configuration flow ([160168f](https://github.com/origadmin/runtime/commit/160168f5c7a62113103cb4e85369073fbc352eed))
* **runtime:** ensure config sources are properly loaded and validated ([3e63a5e](https://github.com/origadmin/runtime/commit/3e63a5e20327a7a1cac74113760ad86d536fed42))
* **runtime:** Handle empty Discoveries array and update config decoder tests ([e69fb00](https://github.com/origadmin/runtime/commit/e69fb003850e5563169fd0ad1aaf812e2b7e7a57))
