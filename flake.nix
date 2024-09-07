{
  description = "A codegen tool for NanoPack";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?tag=24.05";
  };

  outputs = { self, nixpkgs, ... }:
    let
      supportedSystems = [ "x86_64-linux" "x86_64-darwin" "aarch64-linux" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        rec {
          nanoc = pkgs.buildGoModule {
            pname = "nanoc";
            version = "0.1.0";
            src = ./.;
            vendorHash = "sha256-wyg35Xnw2TJirFCHX6DQY9OaeOBJf+xKnYvXk3AKzDU=";
            buildInputs = [
              # nanoc requires clang-format in clang-tools
              pkgs.clang-tools
              # nanoc uses biome to format typescript code
              pkgs.biome
              # nanoc uses swift-format to format swift code
              pkgs.swift-format
            ];
          };

          default = nanoc;
        });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.gotools
              # nanoc requires clang-format in clang-tools
              pkgs.clang-tools
              # nanoc uses biome to format typescript code
              pkgs.biome
              # nanoc uses swift-format to format swift code
              pkgs.swift-format
              # used to build c++ examples
              pkgs.cmake
            ];
          };
        });
    };
}
