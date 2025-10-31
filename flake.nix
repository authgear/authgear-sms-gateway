{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils = {
      url = "github:numtide/flake-utils";
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
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            (final: prev: {
              go = (
                prev.go.overrideAttrs {
                  version = "1.24.9";
                  src = prev.fetchurl {
                    url = "https://go.dev/dl/go1.24.9.src.tar.gz";
                    hash = "sha256-xy+BulT+AO/n8+dJnUAJeSRogbE7d16am7hVQcEb5pU=";
                  };
                }
              );
            })
          ];
        };
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go

            (pkgs.golangci-lint.overrideAttrs (
              prev:
              let
                version = "2.5.0";
              in
              {
                inherit version;
                src = pkgs.fetchFromGitHub {
                  owner = "golangci";
                  repo = "golangci-lint";
                  rev = "v${version}";
                  hash = "sha256-7dHr7cd+yYofIb+yR2kKfj0k0onLH2W/YuxNor7zPeo=";
                };
                vendorHash = "sha256-QEYbFz7SJxLMblkNqaRLDn/PO+mtSPvNYiEUmZh0sLQ=";
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
          ];
        };
      }
    );
}
