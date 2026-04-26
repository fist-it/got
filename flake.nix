{
  description = "got - a git-like version control system in Go";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
    in
    {
      packages.${system}.default = pkgs.buildGoModule {
        pname = "got";
        version = "0.1.0";
        src = ./.;
        vendorHash = null;
      };

      apps.${system}.default = {
        type = "app";
        program = "${self.packages.${system}.default}/bin/got";
      };

      devShells.${system}.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          go
          gopls
          delve
        ];
      };
    };
}
