# Changelog

## [v0.3.3](https://github.com/masawada/mdp/compare/v0.3.2...v0.3.3) - 2026-02-27
- Fix deprecated archives.format in goreleaser config by @masawada in https://github.com/masawada/mdp/pull/33
- Exclude gosec taint analysis rules (G703, G705) for CLI tool by @masawada in https://github.com/masawada/mdp/pull/36
- Bump Songmu/tagpr from 1.15.0 to 1.17.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/35
- Update installation section in README by @masawada in https://github.com/masawada/mdp/pull/37

## [v0.3.2](https://github.com/masawada/mdp/compare/v0.3.1...v0.3.2) - 2026-02-07
- Bump actions/setup-go from 6.1.0 to 6.2.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/23
- Bump Songmu/tagpr from 1.11.0 to 1.12.1 by @dependabot[bot] in https://github.com/masawada/mdp/pull/24
- Bump actions/checkout from 6.0.1 to 6.0.2 by @dependabot[bot] in https://github.com/masawada/mdp/pull/25
- Bump Songmu/tagpr from 1.12.1 to 1.14.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/27
- Bump Songmu/tagpr from 1.14.0 to 1.15.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/28
- Replace panic with error returns in config package by @masawada in https://github.com/masawada/mdp/pull/29
- Add errorf helper method to reduce error output boilerplate in cli by @masawada in https://github.com/masawada/mdp/pull/30
- Optimize rendering process by reusing goldmark instance and deduplicating convert logic by @masawada in https://github.com/masawada/mdp/pull/31
- Deduplicate config loading in CLI by @masawada in https://github.com/masawada/mdp/pull/32

## [v0.3.1](https://github.com/masawada/mdp/compare/v0.3.0...v0.3.1) - 2026-01-11
- Bump github.com/yuin/goldmark from 1.7.13 to 1.7.16 by @dependabot[bot] in https://github.com/masawada/mdp/pull/20
- Bump Songmu/tagpr from 1.10.0 to 1.11.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/19

## [v0.3.0](https://github.com/masawada/mdp/compare/v0.2.0...v0.3.0) - 2025-12-20
- Bump Songmu/tagpr from 1.9.1 to 1.10.0 by @dependabot[bot] in https://github.com/masawada/mdp/pull/16
- Add title extraction support for theme templates by @masawada in https://github.com/masawada/mdp/pull/18

## [v0.2.0](https://github.com/masawada/mdp/compare/v0.1.1...v0.2.0) - 2025-12-14
- Add config path fallback to $HOME/.config and support .yml extension by @masawada in https://github.com/masawada/mdp/pull/9
- Update README with installation and config documentation by @masawada in https://github.com/masawada/mdp/pull/11
- Fix tilde expansion for output_dir config by @masawada in https://github.com/masawada/mdp/pull/12
- Run tests in Docker to isolate from local config by @masawada in https://github.com/masawada/mdp/pull/13
- Add --watch flag for automatic regeneration on file changes by @masawada in https://github.com/masawada/mdp/pull/14
- Specify MIT License for this project by @masawada in https://github.com/masawada/mdp/pull/15

## [v0.1.1](https://github.com/masawada/mdp/compare/v0.1.0...v0.1.1) - 2025-12-13
- Fix goreleaser config option typo by @masawada in https://github.com/masawada/mdp/pull/7

## [v0.1.0](https://github.com/masawada/mdp/commits/v0.1.0) - 2025-12-13
- Add --list option to display generated files by @masawada in https://github.com/masawada/mdp/pull/1
- Add `make lint` target for golangci-lint by @masawada in https://github.com/masawada/mdp/pull/2
- Add GitHub Actions workflow for automated releases with tagpr by @masawada in https://github.com/masawada/mdp/pull/3
- Fix release workflow with proper permissions by @masawada in https://github.com/masawada/mdp/pull/5
- Add goreleaser integration for automated binary releases by @masawada in https://github.com/masawada/mdp/pull/6
