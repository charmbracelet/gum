{ pkgs }:

pkgs.buildGoModule rec {
  pname = "gum";
  version = "0.14.0";

  src = ./.;

  vendorHash = "sha256-UNBDVIz2VEizkhelCjadkzd2S2yTYXecTFUpCf+XtxY=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
