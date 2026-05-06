{
  description = "FRED MCP server — MCP tools for Federal Reserve Economic Data";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "aarch64-linux"
        "x86_64-linux"
        "aarch64-darwin"
        "x86_64-darwin"
      ];
      forEachSystem = f: nixpkgs.lib.genAttrs systems
        (system: f nixpkgs.legacyPackages.${system});
    in
    {
      devShells = forEachSystem (pkgs: {
        default = pkgs.mkShell {
          packages = with pkgs; [ go_1_25 golangci-lint ];
        };
      });

      overlays.default = final: _prev: {
        fred-mcp = final.callPackage ./nix/package.nix { };
      };

      packages = forEachSystem (pkgs: {
        default = pkgs.callPackage ./nix/package.nix { };
      });
    };
}
