{
  description = "A basic flake with a shell";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    systems.url = "github:nix-systems/default";
    flake-utils = {
      url = "github:numtide/flake-utils";
      inputs.systems.follows = "systems";
    };
  };

  outputs =
    {
      nixpkgs,
      flake-utils,
      ...
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            # This is go 1.23.4 unless you update flake.lock
            pkgs.go

            (pkgs.golangci-lint.overrideAttrs (
              prev:
              let
                version = "1.63.4";
              in
              {
                inherit version;
                src = pkgs.fetchFromGitHub {
                  owner = "golangci";
                  repo = "golangci-lint";
                  rev = "v${version}";
                  hash = "sha256-7nIo6Nuz8KLuQlT7btjnTRFpOl+KVd30v973HRKzh08=";
                };
                vendorHash = "sha256-atr4HMxoPEfGeaNlHqwTEAcvgbSyzgCe262VUg3J86c=";
                # We do not actually override anything here,
                # but if we do not repeat this, ldflags refers to the original version.
                ldflags = [
                  "-s"
                  "-X main.version=${version}"
                  "-X main.commit=v${version}"
                  "-X main.date=19700101-00:00:00"
                ];
              }
            ))

            (pkgs.buildGoModule {
              name = "govulncheck";
              src = pkgs.fetchgit {
                url = "https://go.googlesource.com/vuln";
                rev = "refs/tags/v1.1.3";
                hash = "sha256-ydJ8AeoCnLls6dXxjI05+THEqPPdJqtAsKTriTIK9Uc=";
              };
              vendorHash = "sha256-jESQV4Na4Hooxxd0RL96GHkA7Exddco5izjnhfH6xTg=";
              subPackages = [ "cmd/govulncheck" ];
              # checkPhase by default run tests. Running tests will result in build error.
              # So we skip it.
              doCheck = false;
            })

            (pkgs.buildGoModule {
              name = "goimports";
              src = pkgs.fetchgit {
                url = "https://go.googlesource.com/tools";
                rev = "refs/tags/v0.28.0";
                hash = "sha256-BCxsVz4f2h75sj1LzDoKvQ9c8P8SYjcaQE9CdzFdt3w=";
              };
              vendorHash = "sha256-MSir25OEmQ7hg0OAOjZF9J5a5SjlJXdOc523uEBSOSs=";
              subPackages = [ "cmd/goimports" ];
            })
          ];
        };
      }
    );
}
