{
  description = "trigo - a tree command replacement written in Go";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
  };

  outputs =
    { self, nixpkgs }:
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.buildGoModule {
            pname = "trigo";
            version = "0.1.0";
            ldflags = [
              "-s"
              "-w"
            ];
            postInstall = ''
              installShellCompletion --cmd trigo \
                --bash <($out/bin/trigo completion bash) \
                --fish <($out/bin/trigo completion fish) \
                --zsh <($out/bin/trigo completion zsh)
            '';
            nativeBuildInputs = [ pkgs.installShellFiles ];
            src = ./.;
            vendorHash = "sha256-zseKrrna1jGda6lES3MCRBDGlv42fwKQmT5KiJDhmd0=";
          };
        }
      );
    };
}
