# Changelog

## [0.7.0](https://github.com/Marcel2603/tfcoach/compare/v0.6.0..0.7.0) - 2025-11-12

### Features

- *(app)* Add config from home directory in loader hierarchy, add flag for custom config location ([#39](https://github.com/marcel2603/tfcoach/issues/39)) - ([5c1509f](https://github.com/Marcel2603/tfcoach/commit/5c1509f5f2b569775ce5760f35bb000044c3ba24))  by @Erkenbend
- *(rule)* Ignore rules based on comments ([#40](https://github.com/marcel2603/tfcoach/issues/40)) - ([3c2bc9c](https://github.com/Marcel2603/tfcoach/commit/3c2bc9c99141b73d243a4a9398fcc11a4de9becb))  by @Marcel2603

## [0.6.0](https://github.com/Marcel2603/tfcoach/compare/v0.5.0..v0.6.0) - 2025-10-26

### Features

- *(app)* Add educational output format (grouped by rule and ordered by severity) ([#33](https://github.com/marcel2603/tfcoach/issues/33)) - ([16c9d03](https://github.com/Marcel2603/tfcoach/commit/16c9d03a09d1584235b70232047a4cf552272f82))  by @Erkenbend
- *(rule)* Adding 'avoid_null_provider' rule ([#37](https://github.com/marcel2603/tfcoach/issues/37)) - ([f6cb4c9](https://github.com/Marcel2603/tfcoach/commit/f6cb4c983fb0eaa9791df336912b4c303265e9f1))  by @Marcel2603

### Performance

- *(formatting)* Format the whole repo ([#35](https://github.com/marcel2603/tfcoach/issues/35)) - ([e97c89c](https://github.com/Marcel2603/tfcoach/commit/e97c89ca1f1135cc7b32dc96201c738df0f39b21))  by @Marcel2603

### Miscellaneous Tasks

- *(precommit)* Enable precommit ([#34](https://github.com/marcel2603/tfcoach/issues/34)) - ([9cd8c2c](https://github.com/Marcel2603/tfcoach/commit/9cd8c2cd3d7b05dafc92596eca95829da5773da3))  by @Marcel2603

## [0.5.0](https://github.com/Marcel2603/tfcoach/compare/v0.4.0..v0.5.0) - 2025-10-21

### Features

- *(app)* Add flag to lint command for pretty output format ([#26](https://github.com/marcel2603/tfcoach/issues/26)) - ([1f0c05f](https://github.com/Marcel2603/tfcoach/commit/1f0c05fb3b514fbe7d39b0868c42d307a6699b3f))  by @Erkenbend
- *(rules)* Add new rule "core.enforce_variable_description" ([#31](https://github.com/marcel2603/tfcoach/issues/31)) - ([a7df136](https://github.com/Marcel2603/tfcoach/commit/a7df13670fb28f57a1e0d49fc0f836e992c22462))  by @Erkenbend

### Bug Fixes

- *(release)* Version and commit are getting set correctly - ([65f80b7](https://github.com/Marcel2603/tfcoach/commit/65f80b798041e167187eb5e7725879e6437c123a))  by @Marcel2603

### Documentation

- *(README)* Add badges to readme, reoder the first 3 sections - ([b6601b7](https://github.com/Marcel2603/tfcoach/commit/b6601b7c679de3cb4367654915b0b5367de97538))  by @Marcel2603
- *(changelog)* Add contribute on every commit, add dependency-section - ([df01639](https://github.com/Marcel2603/tfcoach/commit/df01639e61c3f970745ba2c79f54b008e0a91362))  by @Marcel2603
- *(generation)* Generate usage and rules-overview from code ([#29](https://github.com/marcel2603/tfcoach/issues/29)) - ([2528ce0](https://github.com/Marcel2603/tfcoach/commit/2528ce03b10a862bc5b54909479fa2df180c140f))  by @Marcel2603
- *(release)* Document release strategy - ([98063aa](https://github.com/Marcel2603/tfcoach/commit/98063aa05d49314021f73c48a5856fe957dd6858))  by @Marcel2603
- *(usage)* Fix available outputs ([#27](https://github.com/marcel2603/tfcoach/issues/27)) - ([5b57cd0](https://github.com/Marcel2603/tfcoach/commit/5b57cd02872b2cd736b628f50121ad59183d7feb))  by @Marcel2603

### Testing

- *(revive)* Enforce linting of go code, let ci fail if warnings/err… ([#30](https://github.com/marcel2603/tfcoach/issues/30)) - ([4ccbf2e](https://github.com/Marcel2603/tfcoach/commit/4ccbf2ee117fd112cefbbc450a46eada67dd71c9))  by @Marcel2603

## [0.4.0](https://github.com/Marcel2603/tfcoach/compare/v0.3.0..v0.4.0) - 2025-10-15

### Features

- *(app)* Add flag to lint command for JSON output format ([#12](https://github.com/marcel2603/tfcoach/issues/12)) - ([9c42309](https://github.com/Marcel2603/tfcoach/commit/9c42309df0e846b483ae0858c8109d682918b1db))  by @Erkenbend
- *(config)* Make tfcoach configurable ([#19](https://github.com/marcel2603/tfcoach/issues/19)) - ([01ef141](https://github.com/Marcel2603/tfcoach/commit/01ef141eb460e45957a6b527545e8523593b7455))  by @Marcel2603

### Documentation

- *(codeowners)* Secure all files - ([4979c78](https://github.com/Marcel2603/tfcoach/commit/4979c78cfa12cd0cd0849e2653f6f33edae5d3a0))  by @Marcel2603

### Miscellaneous Tasks

- *(preview)* Disable workflow - ([3dfe6d7](https://github.com/Marcel2603/tfcoach/commit/3dfe6d72e3599748e37f175bc092e011682f13ec))  by @Marcel2603
- *(preview)* Disable workflow - ([baccb89](https://github.com/Marcel2603/tfcoach/commit/baccb89d6056ecb332837c2cf4d89c7f6a733b3d))  by @Marcel2603

### Dependencies

- *(deps)* Update actions/setup-python action to v6 ([#24](https://github.com/marcel2603/tfcoach/issues/24)) - ([7ca767d](https://github.com/Marcel2603/tfcoach/commit/7ca767d0cea0bf0942b710f0492e931f8c9c387a))  by @renovate[bot]
- *(deps)* Update actions/setup-go action to v6 ([#23](https://github.com/marcel2603/tfcoach/issues/23)) - ([c937f38](https://github.com/Marcel2603/tfcoach/commit/c937f3818e457b32248b241554d19b96f2561ed0))  by @renovate[bot]
- *(deps)* Update davidanson/markdownlint-cli2-action action to v20 ([#25](https://github.com/marcel2603/tfcoach/issues/25)) - ([53ebb60](https://github.com/Marcel2603/tfcoach/commit/53ebb60cdb509619d3dbaa3e6c54e1b3e0df7f45))  by @renovate[bot]
- *(deps)* Update python docker tag to v3.14 ([#21](https://github.com/marcel2603/tfcoach/issues/21)) - ([c50e9d2](https://github.com/Marcel2603/tfcoach/commit/c50e9d25204709a2dedddcedfe0c8516d6e33f1d))  by @renovate[bot]
- *(deps)* Update actions/checkout action to v5 ([#22](https://github.com/marcel2603/tfcoach/issues/22)) - ([5162c4c](https://github.com/Marcel2603/tfcoach/commit/5162c4cce76fdcc6a0fb84536460edd4a76cd16a))  by @renovate[bot]

## New Contributors ❤️

* @renovate[bot] made their first contribution in [#24](https://github.com/Marcel2603/tfcoach/pull/24)
## [0.3.0](https://github.com/Marcel2603/tfcoach/compare/v0.2.0..v0.3.0) - 2025-10-09

### Features

- *(app)* Add new rule 'required_provider_must_be_declared' ([#9](https://github.com/marcel2603/tfcoach/issues/9)) - ([22d4b37](https://github.com/Marcel2603/tfcoach/commit/22d4b377579c0306ca599358d48258e642f3bd26))  by @Erkenbend

### Bug Fixes

- *(rules)* Add terraform block handling in core.file_naming ([#11](https://github.com/marcel2603/tfcoach/issues/11)) - ([56c0cff](https://github.com/Marcel2603/tfcoach/commit/56c0cff00fce2b180b357c470a29d54d4a0431f0))  by @Marcel2603

### Documentation

- *(CHANGELOG)* Adjust changelog after git-cliff change - ([a5a2f94](https://github.com/Marcel2603/tfcoach/commit/a5a2f946838bc83b0906f353079a071f8ecb2a36))  by @Marcel2603
- *(issues)* Fixing typos - ([e4bc08a](https://github.com/Marcel2603/tfcoach/commit/e4bc08a21dc8f11caf1de1656f872cc1bfb69f46))  by @Marcel2603
- *(repo)* Adding first simple CODEOWNERS file - ([10d594e](https://github.com/Marcel2603/tfcoach/commit/10d594e79325e255546f70ceb9991010917a5564))  by @Marcel2603
- Fix typos - ([1989ccf](https://github.com/Marcel2603/tfcoach/commit/1989ccf835a155cd03b9792f05b6b9360f8c01b8))  by @Marcel2603
### Miscellaneous Tasks

- *(release)* Switch to git cliff ([#5](https://github.com/marcel2603/tfcoach/issues/5)) - ([48254d2](https://github.com/Marcel2603/tfcoach/commit/48254d20ad4214e998a13c0a9825e089f7394d95))  by @Marcel2603

## New Contributors ❤️

* @Erkenbend made their first contribution in [#9](https://github.com/Marcel2603/tfcoach/pull/9)
## [0.2.0](https://github.com/Marcel2603/tfcoach/compare/v0.1.0..v0.2.0) - 2025-10-05

### Features

- *(engine)* Process all files (and all rules within those files) asynchronously ([#2](https://github.com/marcel2603/tfcoach/issues/2)) - ([0a84251](https://github.com/Marcel2603/tfcoach/commit/0a8425155b4fd92dc2606f881d4da53c469ebd8e))  by @Marcel2603
- Process all files (and all rules within those files) asynchronously - ([6023a64](https://github.com/Marcel2603/tfcoach/commit/6023a642544b02358e8ae5359cced35bbb9166c9))  by @BusyAnt
### Documentation

- *(github)* Addjust rule issue - ([0f93fe1](https://github.com/Marcel2603/tfcoach/commit/0f93fe1753c82393650030458175d5a246919eb7))  by @Marcel2603
- *(github)* Addjust rule issue - ([7a38eae](https://github.com/Marcel2603/tfcoach/commit/7a38eaeba959001acc74fbb69477de77ce3fe4e8))  by @Marcel2603
- *(github)* Add rule issue - ([457273e](https://github.com/Marcel2603/tfcoach/commit/457273eec7dca15914ea284b6db1706554d63d48))  by @Marcel2603
- *(tooling)* Add hcl test tool, adjust mkdocs, add tool docs - ([b03d99e](https://github.com/Marcel2603/tfcoach/commit/b03d99e368318d0caa514eb675821c0c6dd0eaa6))  by @Marcel2603
- Fix typos in docs, add docker installation - ([336a7de](https://github.com/Marcel2603/tfcoach/commit/336a7de39dc3b11cc5b8021ccf0eeeb3cc17c649))  by @Marcel2603
### Testing

- Add more test cases for engine - ([00b8587](https://github.com/Marcel2603/tfcoach/commit/00b8587e15bdea02db3c62aeefa8462076f2efa3))  by @BusyAnt
### Miscellaneous Tasks

- Fix run command in README - ([fc62c8b](https://github.com/Marcel2603/tfcoach/commit/fc62c8b8ea656d613c1a2d58beec7f9aa206bc46))  by @BusyAnt- Use new sync function WaitGroup.Go instead of WaitGroup.Add/Done - ([7333f49](https://github.com/Marcel2603/tfcoach/commit/7333f492f45ac543385a2c16388b3e9f9ab35be9))  by @BusyAnt
## New Contributors ❤️

* @BusyAnt made their first contribution
## [0.1.0](https://github.com/Marcel2603/tfcoach/compare/v0.0.0..v0.1.0) - 2025-10-04

### Features

- *(app)* Create first release - ([dfd92a1](https://github.com/Marcel2603/tfcoach/commit/dfd92a12b0449e6eb528efd06a56ad155ce78503))  by @Marcel2603
- *(app)* Create first release - ([96f5efe](https://github.com/Marcel2603/tfcoach/commit/96f5efe44e2fee90ababca23643cc39a35a6ae80))  by @Marcel2603
- *(app)* Create first release - ([494b59a](https://github.com/Marcel2603/tfcoach/commit/494b59a7bb77e1f34c57ffd9bb2f73c2a3440af0))  by @Marcel2603

### Documentation

- *(git)* Ignore vscode - ([209896e](https://github.com/Marcel2603/tfcoach/commit/209896e4f4ade14137a5ef87b37df4836ab437c5))  by @Marcel2603
- *(rules)* Add link to rule from overview - ([9f75ba5](https://github.com/Marcel2603/tfcoach/commit/9f75ba59a1af0dcc435094c738df7d754d8043a0))  by @Marcel2603

### Testing

- *(internal+rules)* Adding tests - ([bcad616](https://github.com/Marcel2603/tfcoach/commit/bcad616a7bb83d2c7b04ad9b4bd92d2191af2d3a))  by @Marcel2603

### Miscellaneous Tasks

- *(release)* Add release pipeline - ([d43341f](https://github.com/Marcel2603/tfcoach/commit/d43341fefdc49ab492bf974f2df2a3ed30fe666e))  by @Marcel2603
- *(test)* Fix tests - ([58004ae](https://github.com/Marcel2603/tfcoach/commit/58004ae089073e712924859f7cc3181244103d97))  by @Marcel2603
- *(test)* Rename test workflow - ([f6b5ecb](https://github.com/Marcel2603/tfcoach/commit/f6b5ecb8a02ecfaf99f15f055596881df032336c))  by @Marcel2603

## New Contributors ❤️

* @semantic-release-bot made their first contribution
## [0.0.0]

## New Contributors ❤️

* @Marcel2603 made their first contribution

