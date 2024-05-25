{ pkgs }:

pkgs.buildGoModule rec {
  pname = "gum";
  version = "0.14.0";

  src = ./.;

  vendorHash = "sha256-gDDaKrwlrJyyDzgyGf9iP/XPnOAwpkvIyzCXobXrlF4=";

  ldflags = [ "-s" "-w" "-X=main.Version=${version}" ];
}
